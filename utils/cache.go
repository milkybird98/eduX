package utils

import (
	"net"
	"time"

	"github.com/bluele/gcache"
)

// FileTransmitTag 是cache中存储文件传输数据的格式
type FileTransmitTag struct {
	FileName      string
	ID            string
	FileTags      []string
	Des           string
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
	FileTransmitCache = gcache.New(int(GlobalObject.CacheTableSize)).
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

type RegisterTimerTag struct {
	IP net.Addr
}

var RegisterTimerCache gcache.Cache

func initRegisterTimerCache() {
	RegisterTimerCache = gcache.New(int(GlobalObject.CacheTableSize)).
		LRU().
		Build()
}

func SetRegisterTimerCacheExpire(key string, value RegisterTimerTag) {
	RegisterTimerCache.SetWithExpire(key, value, time.Second*30)
}

func GetRegisterTimerCache(key string) (*RegisterTimerTag, error) {
	value, err := RegisterTimerCache.Get(key)
	if err != nil {
		return nil, err
	}

	file, ok := value.(RegisterTimerTag)
	if !ok {
		return nil, err
	}

	return &file, nil
}

type ResetPasswordTag struct {
	UID string
}

var ResetPasswordCache gcache.Cache

func initResetPasswordCache() {
	ResetPasswordCache = gcache.New(int(GlobalObject.CacheTableSize)).
		LRU().
		Build()
}

func SetResetPasswordCacheExpire(key string, value ResetPasswordTag) {
	ResetPasswordCache.SetWithExpire(key, value, time.Minute*5)
}

func GetResetPasswordCache(key string) (*ResetPasswordTag, error) {
	value, err := ResetPasswordCache.Get(key)
	if err != nil {
		return nil, err
	}

	file, ok := value.(ResetPasswordTag)
	if !ok {
		return nil, err
	}

	return &file, nil
}

func InitCache() {
	initFileTranCache()
	initRegisterTimerCache()
	initResetPasswordCache()
}
