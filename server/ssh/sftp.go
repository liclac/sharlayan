package ssh

import (
	"io"
	"os"

	"github.com/pkg/sftp"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/liclac/sharlayan/config"
)

type sftpSubsystem struct{}

func SFTP(cfg *config.Config) Subsystem {
	if !cfg.SSH.SFTP.Enable {
		zap.L().Named("ssh.sftp").Debug("Not enabled, skipping...")
		return nil
	}
	return &sftpSubsystem{}
}

func (sftpSubsystem) ID() string {
	return "sftp"
}

func (sftpSubsystem) Serve(c *SSHConn, ch ssh.Channel) {
	L := c.L.Named("sftp")

	srv := sftp.NewRequestServer(ch, sftpHandlers(c.Server.cfg, L, c.Fs))
	defer func() {
		if err := srv.Close(); err != nil && err != io.EOF {
			L.Warn("Error closing SFTP server", zap.Error(err))
		} else {
			L.Debug("SFTP session terminated")
		}
	}()
	if err := srv.Serve(); err != nil && err != io.EOF {
		L.Error("SFTP server error", zap.Error(err))
	}
}

type sftpFS struct {
	afero.Fs
	cfg *config.Config
	L   *zap.Logger
}

func sftpHandlers(cfg *config.Config, L *zap.Logger, fs afero.Fs) sftp.Handlers {
	f := sftpFS{fs, cfg, L}
	return sftp.Handlers{
		FileGet:  f,
		FilePut:  f,
		FileCmd:  f,
		FileList: f,
	}
}

func (f sftpFS) Fileread(req *sftp.Request) (io.ReaderAt, error) {
	switch req.Method {
	case "Get":
		if f.cfg.SSH.Trace {
			f.L.Debug("[Trace] Get", zap.String("path", req.Filepath))
		}
		return f.Fs.Open(req.Filepath)
	default:
		f.L.Warn("Unsupported read command", zap.String("method", req.Method),
			zap.String("path", req.Filepath), zap.String("target", req.Target))
		return nil, sftp.ErrSSHFxOpUnsupported
	}
}

func (f sftpFS) Filelist(req *sftp.Request) (sftp.ListerAt, error) {
	switch req.Method {
	case "List":
		if f.cfg.SSH.Trace {
			f.L.Debug("[Trace] List", zap.String("path", req.Filepath))
		}
		infos, err := afero.ReadDir(f, req.Filepath)
		if err != nil {
			return nil, err
		}
		return sftpLister(infos), nil
	case "Stat":
		if f.cfg.SSH.Trace {
			f.L.Debug("[Trace] Stat", zap.String("path", req.Filepath))
		}
		info, err := f.Stat(req.Filepath)
		if err != nil {
			return nil, err
		}
		return sftpLister{info}, nil
	default:
		f.L.Warn("Unsupported list command", zap.String("method", req.Method),
			zap.String("path", req.Filepath), zap.String("target", req.Target))
		return nil, sftp.ErrSSHFxOpUnsupported
	}
}

func (f sftpFS) Filewrite(req *sftp.Request) (io.WriterAt, error) {
	f.L.Debug("Write denied", zap.String("method", req.Method),
		zap.String("path", req.Filepath), zap.String("target", req.Target))
	return nil, sftp.ErrSSHFxPermissionDenied
}

func (f sftpFS) Filecmd(req *sftp.Request) error {
	switch req.Method {
	case "Setstat", "Rename", "Rmdir", "Mkdir", "Link", "Symlink", "Remove":
		f.L.Debug("Write denied", zap.String("method", req.Method),
			zap.String("path", req.Filepath), zap.String("target", req.Target))
		return sftp.ErrSSHFxPermissionDenied
	default:
		f.L.Warn("Unsupported file command", zap.String("method", req.Method),
			zap.String("path", req.Filepath), zap.String("target", req.Target))
		return sftp.ErrSSHFxOpUnsupported
	}
}

// Implements sftp.ListerAt for a slice of FileInfos.
type sftpLister []os.FileInfo

// Copy file descriptors at l[offset:] into ls, return how many were copied.
// If there are no more file descriptors, return io.EOF.
func (l sftpLister) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l)) {
		return 0, io.EOF
	}
	n := copy(ls, l[offset:])
	if n < len(ls) || offset+int64(n) == int64(len(l)) {
		return n, io.EOF
	}
	return n, nil
}
