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

var conncheck_replyStatus string

func (this *PingRouter) PreHandle(request eduiface.IRequest) {
	reqMsgInJSON, conncheck_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("PingRouter: ", conncheck_replyStatus)
		return
	}

	conncheck_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		return
	}

	pingData := gjson.GetBytes(reqMsgInJSON.Data, "ping")
	if !pingData.Exists() {
		conncheck_replyStatus = "data_format_error"
		return
	}

	reqPing := pingData.String()

	if reqPing == "ping" {
		conncheck_replyStatus = "pong"
	} else {
		conncheck_replyStatus = "data_format_error"
	}
}

func (this *PingRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PingRouter: ", conncheck_replyStatus)
	jsonMsg, err := CombineReplyMsg(conncheck_replyStatus, nil)
	if err != nil {
		fmt.Println("PingRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
