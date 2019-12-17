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

type PersonInfoGetRouter struct {
	edunet.BaseRouter
}

type PersonInfoGetData struct{
	UID					string
}

type PersonInfoGetReplyData struct{
	UID					string
	Name				string
	Class				string
	Gender			int
}

type PersonInfoGetByClassRouter struct {
	edunet.BaseRouter
}

type PersonInfoGetByClassData struct{
	Class					string
}

/*
 *	MsgID 100
 *	
 *
 *
 *
 */

 var personget_replyStatus string
 var personget_replyData PersonInfoGetReplyData

func (this *PersonInfoGetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJson ReqMsg
	var reqDataInJson PersonInfoGetData
	reqMsgOrigin := request.GetData()

	checksumFlag = false
	pwdCorrectFlag = false

	err := json.Unmarshal(reqMsgOrigin,&reqMsgInJson)
	if err!=nil{
		fmt.Println(err)
		personget_replyStatus="json_format_error"
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJson.uid))
  md5Ctx.Write(reqMsgInJson.data)

	if utils.SliceEqual(reqMsgInJson.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		personget_replyStatus="check_sum_error"
		return
	}

	err = json.Unmarshal(reqMsgInJson.data,&reqDataInJson)
	if err!=nil{
		fmt.Println(err)
		personget_replyStatus="json_format_error"
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

	personget_replyData.Class = userData.Class
	personget_replyData.Gender = userData.Gender
	personget_replyData.Name = userData.Name
	personget_replyData.UID = userData.UID
}

func (this *PersonInfoGetRouter) Handle(request eduiface.IRequest) {
	var replyMsg ResMsg
	var err error
	
	replyMsg.status = personget_replyStatus
	replyMsg.data,err = json.Marshal(personget_replyData)
	if err!= nil{
		fmt.Println(err)
		return
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
}

/*
 *	MsgID 101
 *	
 *
 *
 *
 */

var persongetbyclass_replyStatus string
var persongetbyclass_replyData []PersonInfoGetReplyData

func (this *PersonInfoGetByClassRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJson ReqMsg
	var reqDataInJson PersonInfoGetByClassData
	reqMsgOrigin := request.GetData()

	checksumFlag = false
	pwdCorrectFlag = false

	err := json.Unmarshal(reqMsgOrigin,&reqMsgInJson)
	if err!=nil{
		fmt.Println(err)
		personget_replyStatus="json_format_error"
		return
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJson.uid))
  md5Ctx.Write(reqMsgInJson.data)

	if utils.SliceEqual(reqMsgInJson.checksum,md5Ctx.Sum(nil)){
		checksumFlag = true
	}else{
		personget_replyStatus="check_sum_error"
		return
	}

	err = json.Unmarshal(reqMsgInJson.data,&reqDataInJson)
	if err!=nil{
		fmt.Println(err)
		personget_replyStatus="json_format_error"
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

	sessionPlcae,err := c.GetSession("Plcae")
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
		if reqDataInJson.Class != sessionClass{
			personget_replyStatus="permission_error"
			return
		}
	}

	userData := edumodel.GetUserByClass(reqDataInJson.Class)
	for _,personData := range *userData{
		persongetbyclass_replyData = append(
			persongetbyclass_replyData,
			PersonInfoGetReplyData{
				personData.UID,
				personData.Name,
				personData.Class,
				personData.Gender})
	}
}

func (this *PersonInfoGetByClassRouter) Handle(request eduiface.IRequest) {
	var replyMsg ResMsg
	var err error
	
	replyMsg.status = personget_replyStatus
	replyMsg.data,err = json.Marshal(persongetbyclass_replyData)
	if err!= nil{
		fmt.Println(err)
		return
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
}