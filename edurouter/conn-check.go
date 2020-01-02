package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"
	"eduX/utils"

	"fmt"
	"time"
)

// PingRouter 用于处理客户端连接检查请求
type PingRouter struct {
	edunet.BaseRouter
}

var conncheckReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PingRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, conncheckReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	conncheckReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 设定返回状态
	conncheckReplyStatus = "pong"
}

// Handle 用于将请求的处理结果发回客户端
func (router *PingRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PingRouter: ", conncheckReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(conncheckReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PingRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
