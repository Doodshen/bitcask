package bitcask

import "errors"

var (
	ErrKeyEmpty          = errors.New("key is empty")
	ErrIndexUpdateFailed = errors.New("index update failed")
	ErrKeyNotFound       = errors.New("key not found")
	ErrDataFileNotFound  = errors.New("data file not found")
)
