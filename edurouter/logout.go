package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"
	"fmt"
)

type LogoutRouter struct {
	edunet.BaseRouter
}

var logout_replyStatus string
var logoutFlag bool

func (this *LogoutRouter) PreHandle(request eduiface.IRequest) {
	logoutFlag = false

	reqMsgInJSON, logout_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("LogoutRouter: ", logout_replyStatus)
		return
	}

	logout_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		return
	}

	sessionUID,err := request.GetConnection().GetSession("UID")
	if err!=nil{
		logout_replyStatus = "session_error"
		return
	}

	if sessionUID == reqMsgInJSON.UID{
		logout_replyStatus = "success"
		logoutFlag = true
	}else{
		logout_replyStatus = "logout_fail"
	}
	
}

func (this *LogoutRouter) Handle(request eduiface.IRequest) {
	fmt.Println("LogoutRouter: ", logout_replyStatus)
	jsonMsg, err := CombineReplyMsg(logout_replyStatus, nil)
	if err != nil {
		fmt.Println("LogoutRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}

func (this *LogoutRouter) PostHandle(request eduiface.IRequest) {
	if logoutFlag {
		c := request.GetConnection()
		c.SetSession("isLogined", false)
		c.Stop()
	}
}
