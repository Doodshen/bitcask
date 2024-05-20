package data

//枚举类型
type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota //正常数据
	LogRecordDelate
)

//LogRecordPos 数据内存索引，主要是描述数据在磁盘上的位置
type LogRecordPos struct {
	Fid    uint32 //文件id 表示数据存放在那个文件中
	Offset int64  //偏移，表示将数据存储到了数据文件中的哪个位置
	Size   uint32 // 标识数据在磁盘上的大小
}

// LogRecord 写入到数据文件的记录
// 之所以叫日志，是因为数据文件中的数据是追加写入的，类似日志的格式
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

//EncodeLogRecord 对LogRecord进行编码 返回编码后的日志记录字节数组以及字节长度
func EncodeLogRecord(LogRecord *LogRecord) ([]byte, int64) {
	return nil, 0
}
