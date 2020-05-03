package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/liclac/sharlayan/server"
	"github.com/liclac/sharlayan/server/ssh"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a server",
	Long:  `Run a server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		L := zap.L().Named("serve")

		// Render the whole tree into an in-memory, read-only filesystem.
		fs := traceFS(cfg, afero.NewMemMapFs())
		if err := buildToFs(cfg, fs); err != nil {
			return err
		}
		fs = afero.NewReadOnlyFs(fs)

		// Make a context, and cancel it if we receive a signal.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // Prevent context goroutine leak.

		go func() {
			defer cancel()

			sigC := make(chan os.Signal, 1)
			defer close(sigC)
			signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
			defer signal.Stop(sigC)

			select {
			case sig := <-sigC:
				L.Info("Signal received, terminating.", zap.Stringer("signal", sig))
			case <-ctx.Done():
				L.Debug("Ceasing signal capture, context expired", zap.Error(ctx.Err()))
			}
		}()

		// Spawn some servers, wait for them to finish, return their error(s).
		return collect(server.Serve(ctx, fs,
			server.HTTP(cfg),
			ssh.Server(cfg, ssh.SFTP(cfg)),
		))
	},
}

func collect(errC <-chan error) (rerr error) {
	for err := range errC {
		rerr = multierror.Append(rerr, err)
	}
	return
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Bool("http.enable", true, "enable the HTTP server")
	serveCmd.Flags().StringP("http.addr", "a", "127.0.0.1:3300", "address for the HTTP server")

	serveCmd.Flags().Bool("ssh.enable", false, "enable the SSH server")
	serveCmd.Flags().String("ssh.addr", "127.0.0.1:3322", "address for the SSH server")
	serveCmd.Flags().Bool("ssh.trace", false, "enable trace logging")
	serveCmd.Flags().String("ssh.host-key", "", "path to host private key (default \"${config.dir}/host_key.pem\")")
	serveCmd.Flags().Bool("ssh.sftp.enable", true, "enable the SFTP subsystem")

	viper.BindPFlags(serveCmd.Flags())
}
