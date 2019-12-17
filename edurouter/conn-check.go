package edurouter

import (
	"fmt"
	"eduX/utils"
	"crypto/md5"
	"encoding/json"
	"eduX/eduiface"
	"eduX/edunet"
)

type PingRouter struct {
	edunet.BaseRouter
}

type PingData struct{
	ping				string
}

var conncheck_replyStatus string

func (this *PingRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJson ReqMsg
	var reqDataInJson PingData

	conncheck_replyStatus=""
	reqMsgOrigin := request.GetData()

	checksumFlag = false

	err := json.Unmarshal(reqMsgOrigin,&reqMsgInJson)
	if err!=nil{
		fmt.Println(err)
		conncheck_replyStatus="json_format_error"
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJson.uid))
  md5Ctx.Write(reqMsgInJson.data)

	if utils.SliceEqual(reqMsgInJson.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		conncheck_replyStatus="check_sum_error"
		return
	}

	err = json.Unmarshal(reqMsgInJson.data,&reqDataInJson)
	if err!=nil{
		fmt.Println(err)
		conncheck_replyStatus="json_format_error"
		return
	}

	c := request.GetConnection()
	value,err := c.GetSession("isLogined")
	if err!= nil {
		conncheck_replyStatus="session_error"
		return
	}

	if value == false{
		conncheck_replyStatus="not_login"
		return
	}

	if reqDataInJson.ping == "ping"{
		conncheck_replyStatus="pong"
	}
}


func (this *PingRouter) Handle(request eduiface.IRequest) {
	var replyMsg ResMsg

	replyMsg.status = conncheck_replyStatus

	md5Ctx = md5.New()
	md5Ctx.Write([]byte(replyMsg.status))
	md5Ctx.Write(replyMsg.data)
	replyMsg.checksum = md5Ctx.Sum(nil)

	jsonMsg,err :=json.Marshal(replyMsg)
	if err!= nil{
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(),jsonMsg)
}
