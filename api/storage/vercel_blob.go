package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Wire protocol of the Vercel Blob API as spoken by @vercel/blob v2.6.x.
const (
	defaultBlobAPIURL = "https://vercel.com/api/blob"
	blobAPIVersion    = "12"
	blobPathPrefix    = "certificates"
)

// VercelBlob stores files in a Vercel Blob store. Refs are the blob URLs
// returned on upload. With a private store (recommended for certificates)
// those URLs only resolve with the read-write token, so files stay
// inaccessible to anyone bypassing the API's own auth checks.
type VercelBlob struct {
	token   string
	storeID string
	apiURL  string
	access  string
	client  *http.Client
}

func NewVercelBlob(token string) *VercelBlob {
	apiURL := os.Getenv("VERCEL_BLOB_API_URL")
	if apiURL == "" {
		apiURL = defaultBlobAPIURL
	}

	access := os.Getenv("BLOB_ACCESS")
	if access != "private" && access != "public" {
		access = "private"
	}

	// Read-write tokens look like vercel_blob_rw_<storeId>_<secret>.
	storeID := ""
	if parts := strings.Split(token, "_"); len(parts) >= 4 {
		storeID = parts[3]
	}

	return &VercelBlob{
		token:   token,
		storeID: storeID,
		apiURL:  apiURL,
		access:  access,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

type blobPutResponse struct {
	URL         string `json:"url"`
	DownloadURL string `json:"downloadUrl"`
	Pathname    string `json:"pathname"`
	ContentType string `json:"contentType"`
}

func (v *VercelBlob) Save(fileName string, content io.Reader, contentType string) (string, error) {
	pathname := path.Join(blobPathPrefix, path.Base(fileName))
	endpoint := fmt.Sprintf("%s/?pathname=%s", v.apiURL, url.QueryEscape(pathname))

	req, err := http.NewRequest(http.MethodPut, endpoint, content)
	if err != nil {
		return "", err
	}
	v.setCommonHeaders(req)
	req.Header.Set("x-vercel-blob-access", v.access)
	req.Header.Set("x-add-random-suffix", "1")
	if contentType != "" {
		req.Header.Set("x-content-type", contentType)
	}

	res, err := v.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return "", fmt.Errorf("vercel blob upload failed: status %d: %s", res.StatusCode, string(body))
	}

	var parsed blobPutResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("vercel blob upload: invalid response: %w", err)
	}
	if parsed.URL == "" {
		return "", fmt.Errorf("vercel blob upload: response missing url")
	}
	return parsed.URL, nil
}

func (v *VercelBlob) Open(ref string) (io.ReadCloser, error) {
	if !IsBlobRef(ref) {
		return nil, os.ErrNotExist
	}

	req, err := http.NewRequest(http.MethodGet, ref, nil)
	if err != nil {
		return nil, err
	}
	// Private blob reads are plain GETs on the blob URL with the token.
	req.Header.Set("Authorization", "Bearer "+v.token)

	res, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusNotFound {
		_ = res.Body.Close()
		return nil, os.ErrNotExist
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		_ = res.Body.Close()
		return nil, fmt.Errorf("vercel blob read failed: status %d", res.StatusCode)
	}
	return res.Body, nil
}

func (v *VercelBlob) Delete(ref string) error {
	payload, err := json.Marshal(map[string][]string{"urls": {ref}})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, v.apiURL+"/delete", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	v.setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	res, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return fmt.Errorf("vercel blob delete failed: status %d: %s", res.StatusCode, string(body))
	}
	return nil
}

func (v *VercelBlob) setCommonHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+v.token)
	req.Header.Set("x-api-version", blobAPIVersion)
	if v.storeID != "" {
		req.Header.Set("x-vercel-blob-store-id", v.storeID)
	}
}
