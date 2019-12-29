package utils

import (
	"eduX/eduiface"
	"encoding/json"
	"io/ioutil"
	"os"
)

/*
	存储一切有关eduX框架的全局参数，供其他模块使用
	一些参数也可以通过 用户根据 eduX.json来配置
*/
type GlobalObj struct {
	/*
		Server
	*/
	TcpServer eduiface.IServer //当前eduX的全局Server对象
	Host      string           //当前服务器主机IP
	TcpPort   int              //当前服务器主机监听端口号
	Name      string           //当前服务器名称

	/*
		eduX
	*/
	Version          string //当前eduX版本号
	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度

	/*
		config file path
	*/
	ConfFilePath string

	/*
		DataBase
	*/
	DataBaseUrl  string
	DataBaseName string

	/*
		Cache
	*/
	CacheTableSize int
}

/*
	定义一个全局的对象
*/
var GlobalObject *GlobalObj

//判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//读取用户的配置文件
func (g *GlobalObj) Reload() {

	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		//fmt.Println("Config File ", g.ConfFilePath , " is not exist!!")
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "eduXServerApp",
		Version:          "V0.1",
		TcpPort:          23333,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     "conf/eduX.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		DataBaseUrl:      "mongodb://localhost:27017",
		DataBaseName:     "eduPlatform",
		CacheTableSize:   4096,
	}

	InitCache()

	//从配置文件中重新加载一些用户配置的参数
	GlobalObject.Reload()
}
