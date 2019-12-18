package edurouter

import (
	"fmt"
	"eduX/utils"
	"crypto/md5"
	"encoding/json"
	"eduX/eduiface"
	"eduX/edunet"
)

var logoutFlag bool

type LogoutRouter struct {
	edunet.BaseRouter
}

type LogoutData struct {
	uid				string
}

func (this *LogoutRouter) Handle(request eduiface.IRequest) {
	var reqMsgInJSON ReqMsg
	var reqDataInJson LogoutData
	var replyMsg ResMsg

	reqMsgOrigin := request.GetData()
	logoutFlag = false
	checksumFlag = false

	err := json.Unmarshal(reqMsgOrigin,&reqMsgInJSON)
	if err!=nil{
		fmt.Println(err)
		checksumFlag = false
		replyMsg.status="json_format_error"
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJSON.uid))
  md5Ctx.Write(reqMsgInJSON.data)

	if utils.SliceEqual(reqMsgInJSON.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		checksumFlag = false
		replyMsg.status="check_sum_error"
	}

	if checksumFlag {
		err = json.Unmarshal(reqMsgInJSON.data,&reqDataInJson)
		if err!=nil{
			fmt.Println(err)
			checksumFlag = false
			replyMsg.status="json_format_error"
		}
	}

	c := request.GetConnection()
	value,err := c.GetSession("isLogined")
	if err!= nil {
		checksumFlag = false
		replyMsg.status="session_error"
	}

	if checksumFlag {
		if value == false{
			checksumFlag = false
			replyMsg.status="not_login"
		}
	}

	if checksumFlag{
		if reqDataInJson.uid == reqMsgInJSON.uid{
			replyMsg.status="logout_success"
			logoutFlag = true
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

	c.SendMsg(request.GetMsgID(),jsonMsg)
}

func (this *LogoutRouter) PostHandle(request eduiface.IRequest){ 
	if logoutFlag {
		c := request.GetConnection()
		c.SetSession("isLogined",false)
		c.Stop()
	}
}