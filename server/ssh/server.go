package ssh

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/liclac/sharlayan/config"
	"github.com/liclac/sharlayan/server"
)

type Subsystem interface {
	ID() string
	Serve(conn *SSHConn, ch ssh.Channel)
}

type SSHServer struct {
	L   *zap.Logger
	cfg *config.Config

	sub map[string]Subsystem
}

func Server(cfg *config.Config, subs ...Subsystem) server.Server {
	L := zap.L().Named("ssh")
	if !cfg.SSH.Enable {
		L.Debug("Not enabled, skipping...")
		return nil
	}
	subMap := map[string]Subsystem{}
	for _, sub := range subs {
		if sub != nil {
			id := sub.ID()
			subMap[id] = sub
			L.Debug("Subsystem registered", zap.String("id", id))
		}
	}
	return &SSHServer{L, cfg, subMap}
}

func (s *SSHServer) Run(ctx context.Context, fs afero.Fs) error {
	// Configure an SSH server...
	sshConfig := &ssh.ServerConfig{
		NoClientAuth: true,
		AuthLogCallback: func(conn ssh.ConnMetadata, method string, err error) {
			L := s.L.With(zap.String("user", conn.User()),
				zap.Stringer("addr", conn.RemoteAddr()))
			if err != nil {
				L.Debug("Auth Failure", zap.Error(err))
			} else {
				L.Debug("Auth Success")
			}
		},
	}
	connShared := SSHConnShared{s, fs, sshConfig}

	// Load or generate a host key.
	hostKeyPath := s.cfg.SSH.HostKey
	if hostKeyPath == "" {
		hostKeyPath = filepath.Join(s.cfg.Config.Dir, "host_key")
	}
	hostKey, err := LoadOrGenerateHostKey(hostKeyPath)
	if err != nil {
		return fmt.Errorf("ssh: couldn't load or generate host key: %w", err)
	}
	sshConfig.AddHostKey(hostKey)

	// Listen until the context terminates.
	l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", s.cfg.SSH.Addr)
	if err != nil {
		return fmt.Errorf("ssh: couldn't listen: %w", err)
	}
	go func() {
		<-ctx.Done()
		if err := l.Close(); err != nil {
			s.L.Error("Error closing listener", zap.Error(err))
		}
	}()
	s.L.Info("Listening", zap.Stringer("addr", l.Addr()))

	// Accept connections.
	for {
		rawConn, err := l.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil // Ignore errors caused by a closed listener.
			default:
				return fmt.Errorf("accept: %w", err)
			}
		}
		if err := accept(&connShared, rawConn); err != nil {
			s.L.Warn("Failed to accept connection", zap.Error(err))
		}
	}
}
