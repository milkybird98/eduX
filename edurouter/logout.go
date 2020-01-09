package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"time"
)

// LogoutRouter 处理用户登出请求
type LogoutRouter struct {
	edunet.BaseRouter
}

// 返回状态码
var logoutReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *LogoutRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, logoutReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	logoutReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 返回success
	logoutReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *LogoutRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", LogoutRouter: ", logoutReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(logoutReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("LogoutRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}

// PostHandle 用于在返回数据后关闭连接
func (router *LogoutRouter) PostHandle(request eduiface.IRequest) {
	if logoutReplyStatus == "success" {
		c := request.GetConnection()
		c.SetSession("isLogined", false)
	}
}
