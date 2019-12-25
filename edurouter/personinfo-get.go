package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"encoding/json"
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
		persongetReplyStatus = "session_error"
		return
	}
	if sessionUID != personUID {
		sessionPlcae, err := c.GetSession("plcae")
		if err != nil {
			persongetReplyStatus = "session_error"
			return
		}
		if sessionPlcae == "student" {
			persongetReplyStatus = "permission_error"
			return
		} else if sessionPlcae == "teacher" {
			sessionClass, err := c.GetSession("Class")
			if err != nil {
				persongetReplyStatus = "session_error"
				return
			}
			userData = edumodel.GetUserByUID(personUID)
			if userData == nil {
				persongetReplyStatus = "user_not_found"
				return
			}
			if userData.Class != sessionClass {
				persongetReplyStatus = "permission_error"
				return
			}
		}
	}

	persongetReplyStatus = "success"

	persongetReplyData.ClassName = userData.Class
	persongetReplyData.Gender = userData.Gender
	persongetReplyData.Name = userData.Name
	persongetReplyData.UID = userData.UID
}

func (router *PersonInfoGetRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonInfoGetRouter: ", persongetReplyStatus)
	data, err := json.Marshal(persongetReplyData)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonMsg, err := CombineReplyMsg(persongetReplyStatus, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
