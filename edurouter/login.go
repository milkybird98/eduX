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

var passwordData string 
var passwordCorrect bool
var userPlaceFlag string
var checksumFlag bool
var pwdCorrectFlag bool

type LoginRouter struct {
	edunet.BaseRouter
}

type LoginData struct{
	pwd				[]byte
}

func (this *LoginRouter) PreHandle(request eduiface.IRequest) {
	var jsonMsg ReqMsg
	var jsonData LoginData
	originMsg := request.GetData()
	checksumFlag = false
	pwdCorrectFlag = false

	err := json.Unmarshal(originMsg,&jsonMsg)
	if err!=nil{
		fmt.Println(err)
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(jsonMsg.uid))
  md5Ctx.Write(jsonMsg.data)
	

	if utils.SliceEqual(jsonMsg.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		return
	}
	
	userData := edumodel.GetUserByUID(jsonMsg.uid)

	err = json.Unmarshal(jsonMsg.data,&jsonData)
	if err!=nil{
		fmt.Println(err)
		return
	}

	if userData!=nil && utils.SliceEqual(jsonData.pwd,[]byte(userData.Pwd)){
		pwdCorrectFlag = true
		userPlaceFlag = userData.Plcae
	}

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

	c := request.GetConnection()
	jsonMsg,err :=json.Marshal(replyMsg)
	if err!= nil{
		fmt.Println(err)
		return
	}

	c.SendMsg(request.GetMsgID(),jsonMsg)
}

func (this *LoginRouter) PostHandle(request eduiface.IRequest){ 
	if pwdCorrectFlag {
		c := request.GetConnection()
		c.SetSession("isLogined",true)
		c.SetSession("place",userPlaceFlag)
	}
}
