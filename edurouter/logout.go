package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"time"
)

type LogoutRouter struct {
	edunet.BaseRouter
}

var logoutReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *LogoutRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, logoutReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	logoutReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	logoutReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *LogoutRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", LogoutRouter: ", logoutReplyStatus)
	jsonMsg, err := CombineReplyMsg(logoutReplyStatus, nil)
	if err != nil {
		fmt.Println("LogoutRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}

func (router *LogoutRouter) PostHandle(request eduiface.IRequest) {
	if logoutReplyStatus == "success" {
		c := request.GetConnection()
		c.SetSession("isLogined", false)
		c.Stop()
	}
}
