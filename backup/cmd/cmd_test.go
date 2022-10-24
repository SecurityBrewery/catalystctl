package cmd

import (
	"github.com/SecurityBrewery/catalystctl/testdata"
	"io/fs"
	"os"
	"testing"
)

func testFile(t *testing.T, s string) (*os.File, error) {
	t.Helper()

	f, err := os.CreateTemp("", s)
	if err != nil {
		return nil, err
	}

	b, err := fs.ReadFile(testdata.FS, s)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(b)
	return f, err
}
