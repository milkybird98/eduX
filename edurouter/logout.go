package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"
	"fmt"
)

type LogoutRouter struct {
	edunet.BaseRouter
}

var logoutReplyStatus string

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

	sessionUID, err := request.GetConnection().GetSession("UID")
	if err != nil {
		logoutReplyStatus = "session_error"
		return
	}

	if sessionUID == reqMsgInJSON.UID {
		logoutReplyStatus = "success"
	} else {
		logoutReplyStatus = "logout_fail"
	}

}

func (router *LogoutRouter) Handle(request eduiface.IRequest) {
	fmt.Println("LogoutRouter: ", logoutReplyStatus)
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
