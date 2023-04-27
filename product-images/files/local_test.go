package files

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func setupLocal(t *testing.T) (*Local, string, func()) {
	// create a temporary directory
	dir, err := os.MkdirTemp("", "files")
	if err != nil {
		t.Fatal(err)
	}

	l, err := NewLocal(dir, 5*1024*10000)
	if err != nil {
		t.Fatal(err)
	}

	return l, dir, func() {
		// cleanup function
		os.RemoveAll(dir)
	}
}

func TestSavesContentsOfReader(t *testing.T) {
	savePath := "/1/test.png"
	fileContents := "Hello World"
	l, dir, cleanup := setupLocal(t)
	defer cleanup()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	// check the file has been correctly written
	f, err := os.Open(filepath.Join(dir, savePath))
	assert.NoError(t, err)

	// check the contents of the file
	d, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, string(d))
}

func TestGetsContentsAndWritesToWriter(t *testing.T) {
	savePath := "/1/test.png"
	fileContents := "Hello World"
	l, _, cleanup := setupLocal(t)
	defer cleanup()

	// save a file
	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	// Read teh file back
	r, err := l.Get(savePath)
	assert.NoError(t, err)
	defer r.Close()

	// read the full contents of the reader
	d, err := io.ReadAll(r)
	assert.Equal(t, fileContents, string(d))
}
