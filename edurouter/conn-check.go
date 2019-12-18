package edurouter

import (
	"crypto/md5"
	"eduX/eduiface"
	"eduX/edunet"
	"eduX/utils"
	"encoding/json"
	"fmt"
)

type PingRouter struct {
	edunet.BaseRouter
}

type PingData struct {
	ping string
}

var conncheck_replyStatus string

func (this *PingRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON ReqMsg
	var reqDataInJson PingData

	conncheck_replyStatus = ""
	reqMsgOrigin := request.GetData()

	checksumFlag = false

	err := json.Unmarshal(reqMsgOrigin, &reqMsgInJSON)
	if err != nil {
		fmt.Println(err)
		conncheck_replyStatus = "json_format_error"
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJSON.uid))
	md5Ctx.Write(reqMsgInJSON.data)

	if utils.SliceEqual(reqMsgInJSON.checksum, md5Ctx.Sum(nil)) {
		checksumFlag = true
	} else {
		conncheck_replyStatus = "check_sum_error"
		return
	}

	err = json.Unmarshal(reqMsgInJSON.data, &reqDataInJson)
	if err != nil {
		fmt.Println(err)
		conncheck_replyStatus = "json_format_error"
		return
	}

	c := request.GetConnection()
	value, err := c.GetSession("isLogined")
	if err != nil {
		conncheck_replyStatus = "session_error"
		return
	}

	if value == false {
		conncheck_replyStatus = "not_login"
		return
	}

	if reqDataInJson.ping == "ping" {
		conncheck_replyStatus = "pong"
	}
}

func (this *PingRouter) Handle(request eduiface.IRequest) {
	var replyMsg ResMsg

	replyMsg.status = conncheck_replyStatus

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(replyMsg.status))
	md5Ctx.Write(replyMsg.data)
	replyMsg.checksum = md5Ctx.Sum(nil)

	jsonMsg, err := json.Marshal(replyMsg)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
