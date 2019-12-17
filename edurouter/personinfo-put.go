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

type PersonInfoPutRouter struct {
	edunet.BaseRouter
}

type PersonInfoPutData struct{
	UID					string
	Name				string
	Class				string
	Gender			int
}

var personput_replyStatus string

func (this *PersonInfoPutRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJson ReqMsg
	var reqDataInJson PersonInfoPutData
	reqMsgOrigin := request.GetData()

	checksumFlag = false
	pwdCorrectFlag = false

	err := json.Unmarshal(reqMsgOrigin,&reqMsgInJson)
	if err!=nil{
		fmt.Println(err)
		personput_replyStatus="json_format_error"
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJson.uid))
  md5Ctx.Write(reqMsgInJson.data)

	if utils.SliceEqual(reqMsgInJson.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		personput_replyStatus="check_sum_error"
		return
	}

	err = json.Unmarshal(reqMsgInJson.data,&reqDataInJson)
	if err!=nil{
		fmt.Println(err)
		personput_replyStatus="json_format_error"
		return
	}

	c := request.GetConnection()
	value,err := c.GetSession("isLogined")
	if err!= nil {
		personget_replyStatus="session_error"
		return
	}

	if value == false{
		personget_replyStatus="not_login"
		return
	}

	userData := edumodel.GetUserByUID(reqDataInJson.UID)

	sessionUID,err := c.GetSession("UID")
	if err!= nil {
		personget_replyStatus="session_error"
		return
	}
	if sessionUID != reqDataInJson.UID {
		sessionPlcae,err := c.GetSession("plcae")
		if err!= nil {
			personget_replyStatus="session_error"
			return
		}
		if sessionPlcae == "student" {
			personget_replyStatus="permission_error"
			return
		}else if sessionPlcae == "teacher"{
			sessionClass,err := c.GetSession("Class")
			if err!= nil {
				personget_replyStatus="session_error"
				return
			}
			if userData.Class != sessionClass{
				personget_replyStatus="permission_error"
				return
			}
		}
	}

	res := edumodel.UpdateUserByID(reqMsgInJson.uid,reqDataInJson.Class,reqDataInJson.Name,"",reqDataInJson.Gender)
	if res {
		personput_replyStatus="update_success"
		return
	}else{
		personput_replyStatus="update_fail"
		return
	}
}

func (this *PersonInfoPutRouter) Handle(request eduiface.IRequest) {
	var replyMsg ResMsg

	replyMsg.status = personput_replyStatus

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
}