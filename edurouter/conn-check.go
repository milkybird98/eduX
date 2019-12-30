package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"
	"eduX/utils"

	"fmt"
	"time"
)

type PingRouter struct {
	edunet.BaseRouter
}

var conncheckReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PingRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, conncheckReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	conncheckReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	conncheckReplyStatus = "pong"
}

// Handle 用于将请求的处理结果发回客户端
func (router *PingRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PingRouter: ", conncheckReplyStatus)
	jsonMsg, err := CombineReplyMsg(conncheckReplyStatus, nil)
	if err != nil {
		fmt.Println("PingRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
