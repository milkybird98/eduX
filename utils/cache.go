package utils

import (
	"net"
	"time"

	"github.com/bluele/gcache"
)

/*
	文件传输队列元素结构体
*/
type FileTransmitTag struct {
	FileName      string
	Size          uint64
	ClientAddress net.Addr
	ServerToC     bool
	ClientToS     bool
}

//文件传输队列
var FileTransmitCache gcache.Cache

func initCache() {
	FileTransmitCache = gcache.New(GlobalObject.FileTransCacheSize).
		LRU().
		Build()
}

func SetCacheExpire(key string, value interface{}, expireTime int) {
	FileTransmitCache.SetWithExpire(key, value, time.Second * time.Duration(expireTime))
}

func GetCache(key string) (interface{}, error) {
	value, err := FileTransmitCache.Get(key)
	return value, err
}
