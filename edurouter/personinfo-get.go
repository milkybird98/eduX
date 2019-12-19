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
	reqMsgInJSON, personget_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("PersonInfoGetRouter: ", personget_replyStatus)
		return
	}

	personget_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		return
	}

	personInfoGetData := gjson.GetBytes(reqMsgInJSON.Data, "uid")
	if !personInfoGetData.Exists() {
		personget_replyStatus = "data_format_error"
		return
	}

	personUID := personInfoGetData.String()

	c := request.GetConnection()
	var userData *edumodel.User

	sessionUID, err := c.GetSession("UID")
	if err != nil {
		personget_replyStatus = "session_error"
		return
	}
	if sessionUID != personUID {
		sessionPlcae, err := c.GetSession("plcae")
		if err != nil {
			personget_replyStatus = "session_error"
			return
		}
		if sessionPlcae == "student" {
			personget_replyStatus = "permission_error"
			return
		} else if sessionPlcae == "teacher" {
			sessionClass, err := c.GetSession("Class")
			if err != nil {
				personget_replyStatus = "session_error"
				return
			}
			userData = edumodel.GetUserByUID(personUID)
			if userData == nil {
				personget_replyStatus = "user_not_found"
				return
			}
			if userData.Class != sessionClass {
				personget_replyStatus = "permission_error"
				return
			}
		}
	}

	personget_replyStatus = "success"

	personget_replyData.ClassName = userData.Class
	personget_replyData.Gender = userData.Gender
	personget_replyData.Name = userData.Name
	personget_replyData.UID = userData.UID
}

func (this *PersonInfoGetRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonInfoGetRouter: ", personget_replyStatus)
	data, err := json.Marshal(personget_replyData)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonMsg, err := CombineReplyMsg(personget_replyStatus, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
