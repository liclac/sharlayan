package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/liclac/sharlayan/config"
)

type httpServer struct {
	L   *zap.Logger
	cfg *config.Config
}

func HTTP(cfg *config.Config) Server {
	L := zap.L().Named("http")
	if !cfg.HTTP.Enable {
		L.Debug("Not enabled")
		return nil
	}
	return httpServer{L, cfg}
}

func (s httpServer) Run(ctx context.Context, fs afero.Fs) error {
	// HTTP server which gracefully shuts down with the context.
	srv := (&http.Server{
		Handler:     http.FileServer(afero.NewHttpFs(fs)),
		BaseContext: func(net.Listener) context.Context { return ctx },
	})
	go func() {
		<-ctx.Done()

		timeout := 60 * time.Second
		sctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		s.L.Info("Gracefully shutting down...", zap.Duration("timeout", timeout))
		if err := srv.Shutdown(sctx); err != nil {
			s.L.Warn("Graceful shutdown failed", zap.Error(err))
		}
	}()

	// Not using ListenAndServe only so that we can print the real address.
	l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", s.cfg.HTTP.Addr)
	if err != nil {
		return fmt.Errorf("http: couldn't listen: %w", err)
	}
	s.L.Info("Listening", zap.Stringer("addr", l.Addr()))

	// Serve until the goroutine above calls srv.Shutdown().
	if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http: server error: %w", err)
	}
	return nil
}
