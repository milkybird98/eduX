package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

type PersonInfoPutRouter struct {
	edunet.BaseRouter
}

type PersonInfoPutData struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	ClassName string `json:"class"`
	Gender    int    `json:"gender"`
}

var personputReplyStatus string

func (router *PersonInfoPutRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, personputReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	classjoiningetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classjoiningetReplyStatus = "data_format_error"
		return
	}

	var reqDataInJSON PersonInfoPutData
	newPersonInfoData := gjson.ParseBytes(reqMsgInJSON.Data)
	reqDataInJSON.UID = newPersonInfoData.Get("uid").String()
	reqDataInJSON.Name = newPersonInfoData.Get("name").String()
	reqDataInJSON.ClassName = newPersonInfoData.Get("class").String()
	reqDataInJSON.Gender = int(newPersonInfoData.Get("gender").Int())

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		persongetReplyStatus = "session_error"
		return
	}

	userData := edumodel.GetUserByUID(reqDataInJSON.UID)
	if userData == nil {
		persongetReplyStatus = "user_not_found"
		return
	}

	if sessionUID != reqDataInJSON.UID {
		sessionPlace, err := c.GetSession("plcae")
		if err != nil {
			persongetReplyStatus = "session_error"
			return
		}

		if sessionPlace == "student" {
			persongetReplyStatus = "permission_error"
			return
		} else if sessionPlace == "teacher" {
			sessionClass, err := c.GetSession("Class")
			if err != nil {
				persongetReplyStatus = "session_error"
				return
			}
			if userData.Class != sessionClass {
				persongetReplyStatus = "permission_error"
				return
			}
		}
	}

	//修改个人信息
	res := edumodel.UpdateUserByID(reqMsgInJSON.UID, reqDataInJSON.ClassName, reqDataInJSON.Name, "", reqDataInJSON.Gender)
	if res {
		personputReplyStatus = "update_success"
	} else {
		personputReplyStatus = "update_fail"
	}
}

func (router *PersonInfoPutRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonInfoPutRouter: ", personputReplyStatus)
	jsonMsg, err := CombineReplyMsg(personputReplyStatus, nil)
	if err != nil {
		fmt.Println("PersonInfoPutRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
