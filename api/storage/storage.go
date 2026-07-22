// Package storage abstracts where uploaded certificate files live.
//
// Two backends exist:
//   - local disk (default) — files under ./uploads/certificates, unchanged
//     from the original behaviour, for local development and Docker.
//   - Vercel Blob — selected automatically when BLOB_READ_WRITE_TOKEN is set,
//     for cloud deployments where the API filesystem is ephemeral.
//
// A stored file is identified by the "ref" string persisted in
// Certificate.FilePath: a filesystem path for local storage, or the blob URL
// for Vercel Blob. Open dispatches on that shape, so databases with a mix of
// old local rows and new blob rows keep working.
package storage

import (
	"io"
	"os"
	"strings"
	"sync"
)

type Storage interface {
	// Save stores content under fileName and returns the ref to persist.
	Save(fileName string, content io.Reader, contentType string) (string, error)
	// Open returns the content of a previously saved ref. The caller must
	// close the reader. Returns os.ErrNotExist if the object is gone.
	Open(ref string) (io.ReadCloser, error)
	// Delete removes the stored object. Best effort — used to roll back a
	// Save when the accompanying database write fails.
	Delete(ref string) error
}

var (
	defaultStorage Storage
	defaultOnce    sync.Once
)

// Default returns the process-wide storage backend, chosen from the
// environment on first use: Vercel Blob when BLOB_READ_WRITE_TOKEN is set,
// local disk otherwise.
func Default() Storage {
	defaultOnce.Do(func() {
		if token := os.Getenv("BLOB_READ_WRITE_TOKEN"); token != "" {
			defaultStorage = NewVercelBlob(token)
		} else {
			defaultStorage = NewLocal("./uploads/certificates")
		}
	})
	return defaultStorage
}

// IsBlobRef reports whether a persisted ref points at remote blob storage
// rather than the local filesystem.
func IsBlobRef(ref string) bool {
	return strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://")
}
