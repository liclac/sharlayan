package afhack

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/afero"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ afero.File = &TraceFile{}

type TraceFile struct {
	fs *TraceFs
	FD uint64
	F  afero.File
	l  *zap.Logger
}

func NewTraceFile(fs *TraceFs, fd uint64, f afero.File) *TraceFile {
	return &TraceFile{fs: fs, FD: fd, F: f,
		l: fs.L.Named("fd" + strconv.FormatUint(fd, 10))}
}

func (f TraceFile) String() string {
	return fmt.Sprintf("&TraceFile{ %d, %#v }", f.FD, f.F)
}

func (f *TraceFile) Close() error {
	err := f.F.Close()
	if le := f.l.Check(debugOrWarn(err), "Close()"); le != nil {
		le.Write(zap.Error(err))
	}
	return err
}

func (f *TraceFile) Read(p []byte) (int, error) {
	n, err := f.F.Read(p)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("Read([%d]byte{}) %d", len(p), n)
		le.Write(zap.Error(err))
	}
	return n, err
}

func (f *TraceFile) ReadAt(p []byte, off int64) (int, error) {
	n, err := f.F.ReadAt(p, off)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("ReadAt([%d]byte{}, %d) %d", len(p), off, n)
		le.Write(zap.Error(err))
	}
	return n, err
}

func (f *TraceFile) Seek(offset int64, whence int) (int64, error) {
	pos, err := f.F.Seek(offset, whence)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("Seek(%d, %d) %d", offset, whence, pos)
		le.Write(zap.Error(err))
	}
	return pos, err
}

func (f *TraceFile) Write(p []byte) (int, error) {
	n, err := f.F.Write(p)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("Write(%q) %d", p, n)
		le.Write(zap.Error(err))
	}
	return n, err
}

func (f *TraceFile) WriteAt(p []byte, off int64) (int, error) {
	n, err := f.F.WriteAt(p, off)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("WriteAt(%q, %d) %d", p, off, n)
		le.Write(zap.Error(err))
	}
	return n, err
}

func (f *TraceFile) Name() string {
	name := f.F.Name()
	if le := f.l.Check(zapcore.DebugLevel, ""); le != nil {
		le.Message = fmt.Sprintf("Name() %q", name)
		le.Write()
	}
	return name
}

func (f *TraceFile) Readdir(count int) ([]os.FileInfo, error) {
	infos, err := f.F.Readdir(count)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("Readdir(%d) %#v", count, wrapFileInfos(infos...))
		le.Write(zap.Error(err))
	}
	return infos, err
}

func (f *TraceFile) Readdirnames(n int) ([]string, error) {
	names, err := f.F.Readdirnames(n)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("Readdirnames(%d) %#v", n, names)
		le.Write(zap.Error(err))
	}
	return names, err
}

func (f *TraceFile) Stat() (os.FileInfo, error) {
	info, err := f.F.Stat()
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("Stat() %#v", wrapFileInfo(info))
		le.Write(zap.Error(err))
	}
	return info, err
}

func (f *TraceFile) Sync() error {
	err := f.F.Sync()
	if le := f.l.Check(debugOrWarn(err), "Sync()"); le != nil {
		le.Write(zap.Error(err))
	}
	return err
}

func (f *TraceFile) Truncate(size int64) error {
	err := f.F.Truncate(size)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("Truncate(%d)", size)
		le.Write(zap.Error(err))
	}
	return err
}

func (f *TraceFile) WriteString(s string) (int, error) {
	ret, err := f.F.WriteString(s)
	if le := f.l.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf("WriteString(%q) %d", s, ret)
		le.Write(zap.Error(err))
	}
	return ret, err
}
