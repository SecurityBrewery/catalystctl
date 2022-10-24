package backup

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type Reader struct {
	f *os.File
	r *zip.Reader
}

func NewReader(f *os.File) (*Reader, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return nil, err
	}

	return &Reader{f: f, r: zr}, nil
}

func (f *Reader) Version() string {
	return f.r.Comment
}

func (f *Reader) Read(p []byte) (n int, err error) {
	return f.f.Read(p)
}

func (f *Reader) Filename() string {
	return f.f.Name()
}

func (f *Reader) Files() []*zip.File {
	return f.r.File
}

func (f *Reader) AsPayload() (*bytes.Buffer, string, error) {
	// set form data
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("backup", f.Filename())
	if err != nil {
		return nil, "", err
	}

	if _, err = io.Copy(fw, f); err != nil {
		return nil, "", err
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return &body, w.FormDataContentType(), nil
}

type Writer struct {
	f *os.File
	w *zip.Writer
}

func NewWriter(file string) (*Writer, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	return &Writer{f: f, w: zip.NewWriter(f)}, nil
}

func NewFromReader(file string, r *Reader) (*Writer, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	w := zip.NewWriter(f)

	for _, zf := range r.Files() {
		zir, err := zf.OpenRaw()
		if err != nil {
			return nil, err
		}

		header := zf.FileHeader
		target, err := w.CreateRaw(&header)
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(target, zir); err != nil {
			return nil, err
		}
	}

	return &Writer{f: f, w: w}, nil
}

func (f *Writer) SetVersion(s string) error {
	return f.w.SetComment(s)
}

func (f *Writer) Close() error {
	if err := f.w.Close(); err != nil {
		return err
	}
	return f.f.Close()
}

func (f *Writer) CreateFolder(s string) error {
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	fh := &zip.FileHeader{Name: s, Method: zip.Store}
	fh.SetMode(os.ModeDir)
	fh.Modified = time.Now()

	_, err := f.w.CreateHeader(fh)
	return err
}
