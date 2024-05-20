package data

import "bitcask/fio"

//数据文件结构体 每个文件对应一个结构体
type DataFile struct {
	Field     uint32        //文件id
	WriteOff  int64         //文件写到位置
	IOManager fio.IOManager // IO读取管理
}

//OpenDataFile 打开数据文件
func OpenDataFile(dirPath string, fieId uint32) (*DataFile, error) {
	return nil, nil
}

//Sync 持久化文件到磁盘
func (df *DataFile) Sync() error {
	return nil
}

//
func (df *DataFile) Write([]byte) error {
	return nil
}

func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, error) {
	return nil, nil
}
