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

func (this *PingRouter) Handle(request eduiface.IRequest) {
	var reqMsgInJson ReqMsg
	var reqDataInJson PingData
	var replyMsg ResMsg

	reqMsgOrigin := request.GetData()

	checksumFlag = false

	err := json.Unmarshal(reqMsgOrigin,&reqMsgInJson)
	if err!=nil{
		fmt.Println(err)
		checksumFlag = false
		replyMsg.status="json_format_error"
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJson.uid))
  md5Ctx.Write(reqMsgInJson.data)

	if utils.SliceEqual(reqMsgInJson.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		checksumFlag = false
		replyMsg.status="check_sum_error"
	}

	if checksumFlag {
		err = json.Unmarshal(reqMsgInJson.data,&reqDataInJson)
		if err!=nil{
			fmt.Println(err)
			checksumFlag = false
			replyMsg.status="json_format_error"
		}
	}

	if checksumFlag{
		if reqDataInJson.ping == "ping"{
			replyMsg.status="pong"
		}
	}

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
