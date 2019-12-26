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
	ID            string
	Size          int64
	ClientAddress net.Addr
	ClassName     string
	UpdaterUID    string
	UpdateTime    time.Time
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
func SetFileTranCacheExpire(key string, value FileTransmitTag) {
	FileTransmitCache.SetWithExpire(key, value, time.Minute*15)
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
