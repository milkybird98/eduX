package eduiface

import (
	"net"
)

//定义连接接口
type IConnection interface {
	//启动连接，让当前连接开始工作
	Start()
	//启动连接，让当前连接开始工作
	StartFileTransmit()
	//停止连接，结束当前连接状态M
	Stop()

	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	//获取当前连接ID
	GetConnID() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr

	//直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	SendMsg(msgId uint32, data []byte) error

	//设置链接属性
	SetSession(key string, value interface{})
	//获取链接属性
	GetSession(key string)(interface{}, error)
	//移除链接属性
	RemoveSession(key string)

}


