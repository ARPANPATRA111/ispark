package storage

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLocalSaveOpenDelete(t *testing.T) {
	dir := t.TempDir()
	s := NewLocal(dir)

	ref, err := s.Save("a.pdf", strings.NewReader("%PDF-fake"), "application/pdf")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	r, err := s.Open(ref)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	got, _ := io.ReadAll(r)
	_ = r.Close()
	if string(got) != "%PDF-fake" {
		t.Fatalf("content mismatch: %q", got)
	}

	if err := s.Delete(ref); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Open(ref); err == nil {
		t.Fatal("Open after Delete should fail")
	}
}

func TestLocalOpenRejectsPathTraversal(t *testing.T) {
	dir := t.TempDir()
	s := NewLocal(dir)

	if _, err := s.Save("safe.pdf", strings.NewReader("data"), ""); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// A tampered ref must resolve inside baseDir, so this looks for a file
	// named passwd in baseDir (absent), never the actual /etc/passwd.
	if _, err := s.Open("../../etc/passwd"); err == nil {
		t.Fatal("path traversal ref should not open")
	}
}

func TestVercelBlobSaveOpenDelete(t *testing.T) {
	var uploaded []byte
	var uploadHeaders http.Header
	var deletedURLs []string

	mux := http.NewServeMux()
	var srv *httptest.Server

	mux.HandleFunc("/api/blob/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("upload method = %s", r.Method)
		}
		uploadHeaders = r.Header.Clone()
		uploaded, _ = io.ReadAll(r.Body)
		pathname := r.URL.Query().Get("pathname")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"url":      srv.URL + "/blobstore/" + pathname,
			"pathname": pathname,
		})
	})
	mux.HandleFunc("/api/blob/delete", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			URLs []string `json:"urls"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		deletedURLs = body.URLs
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/blobstore/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer vercel_blob_rw_store123_secret" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte("blob-content"))
	})

	srv = httptest.NewServer(mux)
	defer srv.Close()

	t.Setenv("VERCEL_BLOB_API_URL", srv.URL+"/api/blob")
	s := NewVercelBlob("vercel_blob_rw_store123_secret")

	ref, err := s.Save("cert.pdf", strings.NewReader("%PDF-blob"), "application/pdf")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if string(uploaded) != "%PDF-blob" {
		t.Fatalf("uploaded body = %q", uploaded)
	}
	if got := uploadHeaders.Get("x-api-version"); got != blobAPIVersion {
		t.Fatalf("x-api-version = %q", got)
	}
	if got := uploadHeaders.Get("x-vercel-blob-store-id"); got != "store123" {
		t.Fatalf("store id header = %q", got)
	}
	if got := uploadHeaders.Get("x-vercel-blob-access"); got != "private" {
		t.Fatalf("access header = %q (private should be the default)", got)
	}
	if !IsBlobRef(ref) {
		t.Fatalf("ref should be a URL, got %q", ref)
	}

	r, err := s.Open(ref)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	got, _ := io.ReadAll(r)
	_ = r.Close()
	if string(got) != "blob-content" {
		t.Fatalf("content = %q", got)
	}

	if err := s.Delete(ref); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if len(deletedURLs) != 1 || deletedURLs[0] != ref {
		t.Fatalf("deleted urls = %v", deletedURLs)
	}
}

func TestVercelBlobOpenMissing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	t.Setenv("VERCEL_BLOB_API_URL", srv.URL)
	s := NewVercelBlob("vercel_blob_rw_x_y")

	_, err := s.Open(srv.URL + "/blobstore/gone.pdf")
	if !os.IsNotExist(err) {
		t.Fatalf("want ErrNotExist, got %v", err)
	}
}
