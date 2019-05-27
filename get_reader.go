package getter

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

// ReaderGetter is a Getter implementation that will download from a io.Reader
type ReaderGetter struct {
	getter
}

func (r *ReaderGetter) ClientMode(_ *url.URL) (ClientMode, error) {
	return ClientModeFile, nil
}

func (r *ReaderGetter) Get(dst string, u *url.URL) error {
	return fmt.Errorf("reader only works for a single file and not for a folder")
}

func (r *ReaderGetter) GetFile(dst string, u *url.URL) error {
	ctx := r.Context()
	// Create all the parent directories if needed
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return err
	}
	defer f.Close()

	var currentFileSize int64

	if fi, err := f.Stat(); err == nil {
		if _, err = f.Seek(0, os.SEEK_END); err == nil {
			currentFileSize = fi.Size()
		}
	}

	if r.client.Rc != nil && r.client.ProgressListener != nil {
		r.client.Rc = r.client.ProgressListener.TrackProgress(r.client.Src,
			currentFileSize, r.client.RcTotalSize, r.client.Rc)
	}

	n, err := Copy(ctx, f, r.client.Rc)
	if err == nil && n < r.client.RcTotalSize {
		err = io.ErrShortWrite
	}
	return err
}
