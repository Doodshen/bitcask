package bitcask

type Options struct {
	//数据库数据目录
	DirPath string

	//活跃文件阈值
	DataFileSize int64

	//每次写数据是否持久化
	SyncWrites bool
}
