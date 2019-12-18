package edurouter

import (
	"fmt"
	"eduX/utils"
	"crypto/md5"
	"encoding/json"
	"eduX/eduiface"
	"eduX/edunet"
	"eduX/edumodel"
)

type LoginRouter struct {
	edunet.BaseRouter
}

type LoginData struct{
	pwd				[]byte
}

var passwordData string 
var passwordCorrect bool
var pwdCorrectFlag bool

func (this *LoginRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON ReqMsg
	var reqDataInJson LoginData
	reqMsgOrigin := request.GetData()

	checksumFlag = false
	pwdCorrectFlag = false

	err := json.Unmarshal(reqMsgOrigin,&reqMsgInJSON)
	if err!=nil{
		fmt.Println(err)
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJSON.uid))
  md5Ctx.Write(reqMsgInJSON.data)

	if utils.SliceEqual(reqMsgInJSON.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		return
	}
	
	userData := edumodel.GetUserByUID(reqMsgInJSON.uid)

	err = json.Unmarshal(reqMsgInJSON.data,&reqDataInJson)
	if err!=nil{
		fmt.Println(err)
		return
	}

	if userData!=nil && utils.SliceEqual(reqDataInJson.pwd,[]byte(userData.Pwd)){
		pwdCorrectFlag = true
	}

	c := request.GetConnection()

	c.SetSession("isLogined",true)
	c.SetSession("UID",userData.UID)
	c.SetSession("place",userData.Plcae)
	c.SetSession("class",userData.Class)
}

func (this *LoginRouter) Handle(request eduiface.IRequest) {
	var replyMsg ResMsg
	if checksumFlag == false{
		replyMsg.status="check_sum_error"
	}else if pwdCorrectFlag{
		replyMsg.status="login_success"
	}else{
		replyMsg.status="login_fail"
	}
	md5Ctx := md5.New()
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

	c.SetSession("isLogined",true)
}
