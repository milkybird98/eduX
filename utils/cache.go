package utils

import (
	"eduX/eduiface"
	"net"
	"time"

	"github.com/bluele/gcache"
)

// FileTransmitTag 是cache中存储文件传输数据的格式
type FileTransmitTag struct {
	FileName      string
	Size          uint64
	ClientAddress net.Addr
	ServerToC     bool
	ClientToS     bool
}

// FileTransmitCache 是文件传输队列缓存
var FileTransmitCache gcache.Cache

func initFileTranCache() {
	FileTransmitCache = gcache.New(int(GlobalObject.FileTransCacheSize)).
		LRU().
		Build()
}

// SetFileTranCacheExpire 用于设定
func SetFileTranCacheExpire(key string, value FileTransmitTag, expireTime int) {
	FileTransmitCache.SetWithExpire(key, value, time.Second*time.Duration(expireTime))
}

func GetFileTranCache(key string) (*FileTransmitTag, error) {
	value, err := FileTransmitCache.Get(key)
	if err != nil {
		return nil, err
	}

	file, ok := value.(FileTransmitTag)
	if !ok {
		return nil, err
	}

	return &file, nil
}

type UserOnlineTag struct {
	LastSeenTime  time.Time
	ClientAddress net.Addr
	UID           string
	Connection    eduiface.IConnection
}

var UserOnlineCache gcache.Cache

func initUserOnlineCache() {
	UserOnlineCache = gcache.New(int(GlobalObject.UserOnlineCacheSize)).
		LRU().
		Build()
}
