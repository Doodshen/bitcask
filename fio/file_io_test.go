package fio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func destroyFile(name string) {
	if err := os.Remove(name); err != nil {
		panic(err)
	}
}

func TestNewFileIOManager(t *testing.T) {
	file, err := NewFileIOManager("./test-file")
	require.NoError(t, err)
	defer destroyFile("./test-file")
	_, err = file.Write([]byte("test"))
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)
}

func TestFileIO_Write(t *testing.T) {
	path := filepath.Join("../tmp", "a.data")
	fio, err := NewFileIOManager(path)
	defer os.Remove(path)

	assert.Nil(t, err)
	assert.NotNil(t, fio)

	n, err := fio.Write([]byte(""))
	assert.Equal(t, 0, n)
	assert.Nil(t, err)

	n, err = fio.Write([]byte("bitcask kv"))
	assert.Equal(t, 10, n)
	assert.Nil(t, err)

	n, err = fio.Write([]byte("storage"))
	assert.Equal(t, 7, n)
	assert.Nil(t, err)
}

func TestFileIO_Sync(t *testing.T) {
	path := filepath.Join("../tmp", "a.data")
	fio, err := NewFileIOManager(path)
	//defer destroyFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, fio)

	err = fio.Sync()
	assert.Nil(t, err)
}

func TestFileIO_Close(t *testing.T) {
	path := filepath.Join("../tmp", "a.data")
	fio, err := NewFileIOManager(path)
	defer destroyFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, fio)

	err = fio.Close()
	assert.Nil(t, err)
}
