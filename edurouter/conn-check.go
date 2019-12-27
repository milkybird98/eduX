package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"

	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

type PingRouter struct {
	edunet.BaseRouter
}

var conncheckReplyStatus string

func (router *PingRouter) PreHandle(request eduiface.IRequest) {
	reqMsgInJSON, conncheckReplyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("PingRouter: ", conncheckReplyStatus)
		return
	}

	conncheckReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	pingData := gjson.GetBytes(reqMsgInJSON.Data, "ping")
	if !pingData.Exists() {
		conncheckReplyStatus = "data_format_error"
		return
	}

	reqPing := pingData.String()

	if reqPing == "ping" {
		conncheckReplyStatus = "pong"
	} else {
		conncheckReplyStatus = "data_format_error"
	}
}

func (router *PingRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ",time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PingRouter: ", conncheckReplyStatus)
	jsonMsg, err := CombineReplyMsg(conncheckReplyStatus, nil)
	if err != nil {
		fmt.Println("PingRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
