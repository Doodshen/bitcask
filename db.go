package bitcask

import (
	"bitcask/data"
	"bitcask/index"
	"sync"
)

//DB bitcask存储引实例

type DB struct {
	options    Options //配置项
	mu         *sync.RWMutex
	activeFile *data.DataFile            // 当前活跃数据文件，可以用于写入
	olderFiles map[uint32]*data.DataFile //旧的数据文件，只能用于读
	index      index.Indexer             //内存索引
}

// Put 写入数据到bitcask活跃文件中，并更新索引
func (db *DB) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return ErrKeyEmpty
	}
	//构建LogRecord结构体
	LogRecord := &data.LogRecord{
		Key:   key,
		Value: value,
		Type:  data.LogRecordNormal,
	}
	//1,拿到内存索引信息
	pos, err := db.appendLogRecord(LogRecord)
	if err != nil {
		return nil
	}
	//2 更新内存索引信息
	if ok := db.index.Put(key, pos); !ok {
		return ErrIndexUpdateFailed
	}
	return nil
}

// Get 根据key查找数据
func (db *DB) Get(key []byte) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	//1 判断key
	if len(key) == 0 {
		return nil, ErrKeyEmpty
	}
	//2 获取内存索引信息
	logRecordPos := db.index.Get(key)
	//如果key不存在内存索引种，说明key不存在
	if logRecordPos == nil {
		return nil, ErrKeyNotFound
	}

	//3 根据文件id找到对应的数据文件
	var dataFile *data.DataFile
	if db.activeFile.Field == logRecordPos.Fid {
		dataFile = db.activeFile
	} else {
		dataFile = db.olderFiles[logRecordPos.Fid]
	}
	//如果数据文件为空，说明数据文件不存在
	if dataFile == nil {
		return nil, ErrDataFileNotFound
	}

	//4 从数据文件中读取数据
	logRecord, err := dataFile.ReadLogRecord(logRecordPos.Offset)
	if err != nil {
		return nil, err
	}
	//5 判断数据类型
	if logRecord.Type == data.LogRecordDelate {
		return nil, ErrKeyNotFound
	}
	return logRecord.Value, nil

}

// appendLogRecord 写入数据到活跃文件，并返回当前写入的offset
func (db *DB) appendLogRecord(logRecoed *data.LogRecord) (*data.LogRecordPos, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	//判断活跃文件是否存在 因为数据库初始化时可能没有活跃文件
	//如果没有活跃文件生成活跃文件
	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}
	//写入数据到活跃文件
	//1 编码
	encRecord, size := data.EncodeLogRecord(logRecoed)
	//如果写入数据长度已经达到了活跃文件阈值，关闭活跃文件，并生成新的活跃文件
	if db.activeFile.WriteOff+size > db.options.DataFileSize {
		//先持久化数据文件，保证数据持久化到磁盘中
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}

		//将当前活跃文件转化为旧的文件
		db.olderFiles[db.activeFile.Field] = db.activeFile

		//打开新的数据文件
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}
	//写入数据到活跃文件
	//1 先将当前文件写入的offset记录下来
	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(encRecord); err != nil {
		return nil, err
	}
	//2 提供配置项决定是否持久化
	if db.options.SyncWrites {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}
	//3构建内存索引信息，返回当前写入的offset
	pos := &data.LogRecordPos{
		Fid:    db.activeFile.Field,
		Offset: writeOff,
	}
	return pos, nil

}

// setActiveDataFile() 初始化当前活跃文件  在该执行该方法之前必须持有锁
func (db *DB) setActiveDataFile() error {
	var initialField uint32 = 0
	//如果当前活跃文件不为空，新创建的活跃文件id +1
	if db.activeFile != nil {
		initialField = db.activeFile.Field + 1
	}
	//打开新的数据文件  定义文件目录配置项，打开数据让用户传递过来
	dataFile, err := data.OpenDataFile(db.options.DirPath, initialField)
	if err != nil {
		return err
	}
	db.activeFile = dataFile
	return nil
}
