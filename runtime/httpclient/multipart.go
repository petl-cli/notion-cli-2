package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// MultipartBody describes a multipart/form-data request body. It is the
// alternative to Request.Body — when set, the runtime ignores Body and builds
// a multipart body instead. JSON-mode requests never touch this type.
type MultipartBody struct {
	// Files maps the form field name to one or more file paths. Each path
	// becomes a separate part whose Content-Type is sniffed from the file.
	// Multiple files under the same field name (e.g. "files") are emitted as
	// repeated parts with the same field name, which is the standard
	// multipart pattern for arrays of files.
	Files map[string][]string

	// Fields holds non-file form fields. Each entry becomes a text part with
	// no filename. Use this for scalar form values mixed alongside files.
	Fields map[string]string
}

// IsEmpty reports whether the multipart body has nothing to send.
// The runtime uses this to skip multipart construction entirely when a
// command declares multipart but the user passed no body inputs.
func (m *MultipartBody) IsEmpty() bool {
	return m == nil || (len(m.Files) == 0 && len(m.Fields) == 0)
}

// buildMultipartBody serializes a MultipartBody into a wire-format multipart
// payload. Returns the body reader and the full Content-Type header value
// (which carries the boundary string the server needs to parse parts).
//
// Files are opened and streamed during serialization; the entire body is
// materialized in memory because the http.Request expects a sized body and
// we already buffer JSON bodies the same way. For very large file uploads
// this could be reworked to stream, but the common case (a handful of MBs)
// is fine and matches the JSON path's memory profile.
func buildMultipartBody(m *MultipartBody) (io.Reader, string, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	// Files first so partial failures don't leave half-written field parts
	// that the server has to reconcile with a missing file.
	for fieldName, paths := range m.Files {
		for _, path := range paths {
			if err := writeFilePart(mw, fieldName, path); err != nil {
				return nil, "", err
			}
		}
	}

	for name, value := range m.Fields {
		if err := mw.WriteField(name, value); err != nil {
			return nil, "", fmt.Errorf("writing field %q: %w", name, err)
		}
	}

	if err := mw.Close(); err != nil {
		return nil, "", fmt.Errorf("finalizing multipart body: %w", err)
	}

	return &buf, mw.FormDataContentType(), nil
}

// writeFilePart appends one file as a part to the multipart writer. The
// filename in Content-Disposition is the basename of path — servers
// typically only look at the basename anyway, and exposing the full local
// path leaks user filesystem layout.
func writeFilePart(mw *multipart.Writer, fieldName, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %q for upload: %w", path, err)
	}
	defer f.Close()

	part, err := mw.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return fmt.Errorf("creating form part for %q: %w", path, err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return fmt.Errorf("copying %q into multipart body: %w", path, err)
	}
	return nil
}
