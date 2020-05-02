package afhack

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ afero.Fs = &TraceFs{}

// A TraceFs wraps an afero.Fs filesystem and logs all IO operations.
// As with all low level tracing, it adds overhead, and the output is quite noisy.
type TraceFs struct {
	L      *zap.Logger
	FS     afero.Fs
	NextFD uint64
}

func NewTraceFs(L *zap.Logger, fs afero.Fs) *TraceFs {
	return &TraceFs{L: L, FS: fs}
}

func (fs *TraceFs) nextFD() uint64 {
	fs.NextFD++
	return fs.NextFD
}

func (fs *TraceFs) Create(name string) (afero.File, error) {
	f, err := fs.FS.Create(name)
	f = NewTraceFile(fs, fs.nextFD(), f)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Create("%s") %s`, name, f)
		le.Write(zap.Error(err))
	}
	return f, err
}

func (fs *TraceFs) Mkdir(name string, perm os.FileMode) error {
	err := fs.FS.Mkdir(name, perm)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Mkdir("%s", %#o)`, name, perm)
		le.Write(zap.Error(err))
	}
	return err
}

func (fs *TraceFs) MkdirAll(path string, perm os.FileMode) error {
	err := fs.FS.MkdirAll(path, perm)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`MkdirAll("%s", %#o)`, path, perm)
		le.Write(zap.Error(err))
	}
	return err
}

func (fs *TraceFs) Open(name string) (afero.File, error) {
	f, err := fs.FS.Open(name)
	f = NewTraceFile(fs, fs.nextFD(), f)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Open("%s") %s`, name, f)
		le.Write(zap.Error(err))
	}
	return f, err
}

func (fs *TraceFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.FS.OpenFile(name, flag, perm)
	f = NewTraceFile(fs, fs.nextFD(), f)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`OpenFile("%s", %d, %#o) %s`, name, flag, perm, f)
		le.Write(zap.Error(err))
	}
	return f, err
}

func (fs *TraceFs) Remove(name string) error {
	err := fs.FS.Remove(name)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Remove("%s")`, name)
		le.Write(zap.Error(err))
	}
	return err
}

func (fs *TraceFs) RemoveAll(path string) error {
	err := fs.FS.RemoveAll(path)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`RemoveAll("%s")`, path)
		le.Write(zap.Error(err))
	}
	return err
}

func (fs *TraceFs) Rename(oldname, newname string) error {
	err := fs.FS.Rename(oldname, newname)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Rename("%s", "%s")`, oldname, newname)
		le.Write(zap.Error(err))
	}
	return err
}

func (fs *TraceFs) Stat(name string) (os.FileInfo, error) {
	info, err := fs.FS.Stat(name)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Stat("%s") %#v`, name, wrapFileInfo(info))
		le.Write(zap.Error(err))
	}
	return info, err
}

func (fs *TraceFs) Name() string {
	name := fs.FS.Name()
	if le := fs.L.Check(zapcore.DebugLevel, ""); le != nil {
		le.Message = fmt.Sprintf(`Name() "%s"`, name)
		le.Write()
	}
	return name
}

func (fs *TraceFs) Chmod(name string, mode os.FileMode) error {
	err := fs.FS.Chmod(name, mode)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Chmod("%s", %#o)`, name, mode)
		le.Write(zap.Error(err))
	}
	return err
}

func (fs *TraceFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	err := fs.FS.Chtimes(name, atime, mtime)
	if le := fs.L.Check(debugOrWarn(err), ""); le != nil {
		le.Message = fmt.Sprintf(`Chtimes("%s", "%s", "%s")`, name, atime, mtime)
		le.Write(zap.Error(err))
	}
	return err
}

func debugOrWarn(err error) zapcore.Level {
	if err != nil {
		return zapcore.WarnLevel
	}
	return zapcore.DebugLevel
}

type TraceFileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	Sys     interface{}
}

func wrapFileInfo(info os.FileInfo) TraceFileInfo {
	return TraceFileInfo{
		Name: info.Name(), Size: info.Size(),
		Mode: info.Mode(), ModTime: info.ModTime(),
		IsDir: info.IsDir(), Sys: info.Sys(),
	}
}

func wrapFileInfos(infos ...os.FileInfo) []TraceFileInfo {
	traceInfos := make([]TraceFileInfo, len(infos))
	for i, info := range infos {
		traceInfos[i] = wrapFileInfo(info)
	}
	return traceInfos
}
