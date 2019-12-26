package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

type PersonInfoGetRouter struct {
	edunet.BaseRouter
}

type PersonInfoGetData struct {
	UID string `json:"uid"`
}

type PersonInfoGetReplyData struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	ClassName string `json:"class"`
	Gender    int    `json:"gender"`
}

var persongetReplyStatus string
var persongetReplyData PersonInfoGetReplyData

func (router *PersonInfoGetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, persongetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	persongetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	personInfoGetData := gjson.GetBytes(reqMsgInJSON.Data, "uid")
	if !personInfoGetData.Exists() {
		persongetReplyStatus = "data_format_error"
		return
	}

	personUID := personInfoGetData.String()

	c := request.GetConnection()
	var userData *edumodel.User

	sessionUID, err := c.GetSession("UID")
	if err != nil {
		persongetReplyStatus = "57session_error"
		return
	}

	userData = edumodel.GetUserByUID(personUID)
	if userData == nil {
		persongetReplyStatus = "user_not_found"
		return
	}

	if sessionUID != personUID {
		sessionPlace, err := c.GetSession("place")
		if err != nil {
			persongetReplyStatus = "64session_error"
			return
		}
		if sessionPlace == "student" {
			persongetReplyStatus = "permission_error"
			return
		} else if sessionPlace == "teacher" {
			sessionClass, err := c.GetSession("class")
			if err != nil {
				persongetReplyStatus = "73session_error"
				return
			}
			if userData.Class != sessionClass {
				persongetReplyStatus = "permission_error"
				return
			}
		}
	}

	fmt.Println(userData.Class)

	persongetReplyData.ClassName = userData.Class
	persongetReplyData.Gender = userData.Gender
	persongetReplyData.Name = userData.Name
	persongetReplyData.UID = userData.UID

	persongetReplyStatus = "success"
}

func (router *PersonInfoGetRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonInfoGetRouter: ", persongetReplyStatus)
	var jsonMsg []byte
	var err error
	if persongetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(persongetReplyStatus, persongetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(persongetReplyStatus, nil)

	}
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
