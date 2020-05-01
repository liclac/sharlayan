package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/config"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a server",
	Long:  `Run a server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		// Render the whole tree into an in-memory, read-only filesystem.
		meta, err := loadMeta(cfg)
		if err != nil {
			return err
		}
		bld, err := createBuilder(cfg)
		if err != nil {
			return err
		}
		fs := afero.NewMemMapFs()
		if err := buildToFs(bld, createNodes(bld, meta), fs); err != nil {
			return err
		}
		fs = afero.NewReadOnlyFs(fs)

		// We'll cancel this context if we get a signal.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Spawn a server.
		httpErrC := make(chan error, 1)
		go func() {
			httpErrC <- serveHTTP(ctx, &cfg, fs)
			close(httpErrC)
		}()

		// Wait for either a signal or an error.
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigC:
			log.WithField("signal", sig).Info("Signal received, shutting down...")
		case err := <-httpErrC:
			return err
		}
		cancel()

		// Wait for any trailing errors.
		return <-httpErrC
	},
}

func serveHTTP(ctx context.Context, cfg *config.Config, fs afero.Fs) error {
	if !cfg.HTTP.Enable {
		log.Debug("HTTP: Disabled")
		return nil
	}
	l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", cfg.HTTP.Addr)
	if err != nil {
		return fmt.Errorf("http: couldn't listen: %w", err)
	}
	log.WithField("addr", l.Addr().String()).Info("HTTP: Listening")

	srv := (&http.Server{
		Handler:     http.FileServer(afero.NewHttpFs(fs)),
		BaseContext: func(net.Listener) context.Context { return ctx },
	})
	go func() {
		<-ctx.Done()
		timeout := 60 * time.Second
		sctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		log.WithField("timeout", timeout).Info("HTTP: Gracefully shutting down...")
		if err := srv.Shutdown(sctx); err != nil {
			log.WithError(err).Warn("HTTP: Couldn't gracefully shut down")
		}
	}()
	if err := srv.Serve(l); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Bool("http.enable", true, "enable the HTTP server")
	serveCmd.Flags().StringP("http.addr", "a", "127.0.0.1:3300", "address for the HTTP server")

	viper.BindPFlags(serveCmd.Flags())
}
