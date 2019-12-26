package edurouter

import (
	"eduX/eduiface"
	"eduX/edunet"

	"fmt"

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
	fmt.Println("PingRouter: ", conncheckReplyStatus)
	jsonMsg, err := CombineReplyMsg(conncheckReplyStatus, nil)
	if err != nil {
		fmt.Println("PingRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
