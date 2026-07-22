package storage

import (
	"io"
	"os"
	"path/filepath"
)

// Local stores files on the API host's filesystem.
type Local struct {
	baseDir string
}

func NewLocal(baseDir string) *Local {
	return &Local{baseDir: baseDir}
}

func (l *Local) Save(fileName string, content io.Reader, _ string) (string, error) {
	if err := os.MkdirAll(l.baseDir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(l.baseDir, filepath.Base(fileName))
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(f, content); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return path, nil
}

func (l *Local) Open(ref string) (io.ReadCloser, error) {
	// Refs are constrained to the base directory so a tampered database row
	// cannot read arbitrary files.
	return os.Open(filepath.Join(l.baseDir, filepath.Base(ref)))
}

func (l *Local) Delete(ref string) error {
	return os.Remove(filepath.Join(l.baseDir, filepath.Base(ref)))
}
