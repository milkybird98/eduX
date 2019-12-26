package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

type PersonInfoGetByClassRouter struct {
	edunet.BaseRouter
}

type PersonInfoGetByClassData struct {
	ClassName string `json:"class"`
}

type PersonInfoGetByClassReplyData struct {
	UserList []PersonInfoGetReplyData `json:"userlist"`
}

/*
 *	MsgID 101
 *
 *
 *
 *
 */

var persongetbyclassReplyStatus string
var persongetbyclassReplyData PersonInfoGetByClassReplyData

func (router *PersonInfoGetByClassRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, persongetbyclassReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	persongetbyclassReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	reqClassNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	if !reqClassNameData.Exists() {
		persongetbyclassReplyStatus = "data_format_error"
		return
	}

	reqClassName := reqClassNameData.String()

	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		persongetbyclassReplyStatus = "59session_error"
		return
	}

	if sessionPlace == "student" {
		persongetbyclassReplyStatus = "permission_error"
		return
	} else if sessionPlace == "teacher" {
		sessionClass, err := c.GetSession("class")
		if err != nil {
			persongetbyclassReplyStatus = "69session_error"
			return
		}
		if reqClassName != sessionClass {
			persongetbyclassReplyStatus = "permission_error"
			return
		}
	}

	userManyData := edumodel.GetUserByClass(reqClassName)
	if userManyData == nil || len(*userManyData) <= 0 {
		persongetbyclassReplyStatus = "data_not_found"
		return
	}
	persongetbyclassReplyStatus = "success"

	for _, personData := range *userManyData {
		persongetbyclassReplyData.UserList = append(
			persongetbyclassReplyData.UserList,
			PersonInfoGetReplyData{
				personData.UID,
				personData.Name,
				personData.Class,
				personData.Gender})
	}
}

func (router *PersonInfoGetByClassRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonInfoGetByClassRouter: ", persongetbyclassReplyStatus)
	var jsonMsg []byte
	var err error

	if persongetbyclassReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(persongetbyclassReplyStatus, persongetbyclassReplyData)

	} else {
		jsonMsg, err = CombineReplyMsg(persongetbyclassReplyStatus, nil)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
