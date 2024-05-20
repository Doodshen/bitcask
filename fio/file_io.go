package fio

import "os"

// 标准文件IO操作 实现IoManaeger接口

type FileIO struct {
	fd *os.File //文件描述符
}

func NewFileIOManager(fileName string) (*FileIO, error) {
	fd, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		DataFilePerm,
	)
	if err != nil {
		return nil, err
	}
	return &FileIO{fd: fd}, nil
}

//Read 从文件指定位置读取指定的数据
func (fio *FileIO) Read(b []byte, offset int64) (int, error) {
	return fio.fd.ReadAt(b, offset)
}

//Write 写入字节数据到文件中
func (fio *FileIO) Write(b []byte) (int, error) {
	return fio.fd.Write(b)
}

//Sync 持久化数据
func (fio *FileIO) Sync() error {
	return fio.fd.Sync()
}

//Close 关闭文件
func (fio *FileIO) Close() error {
	return fio.fd.Close()
}
