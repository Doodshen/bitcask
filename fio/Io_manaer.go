package fio

const DataFilePerm = 0644

// IOManager 抽象IO管理接口，可以接入不同IO类型，
type IOManager interface {
	//Read 从文件指定位置读取指定的数据
	Read([]byte, int64) (int, error)

	//Write 写入字节数据到文件中
	Write([]byte) (int, error)

	//Sync 持久化数据
	Sync() error

	//Close 关闭文件
	Close() error
}
