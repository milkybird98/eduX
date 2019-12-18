package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"encoding/json"
	"fmt"
)

type PersonInfoPutRouter struct {
	edunet.BaseRouter
}

type PersonInfoPutData struct {
	UID    string
	Name   string
	Class  string
	Gender int
}

var personput_replyStatus string

func (this *PersonInfoPutRouter) PreHandle(request eduiface.IRequest) {
	reqMsgInJSON, personput_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("PersonInfoPutRouter: ", personput_replyStatus)
		return
	}

	var reqDataInJson PersonInfoPutData
	err := json.Unmarshal(reqMsgInJSON.data, &reqDataInJson)
	if err != nil {
		fmt.Println(err)
		personput_replyStatus = "json_format_error"
		fmt.Println("PersonInfoPutRouter: ", personput_replyStatus)
		return
	}

	classjoininget_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		fmt.Println("PersonInfoPutRouter: ", personput_replyStatus)
		return
	}

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		personget_replyStatus = "session_error"
		return
	}

	if sessionUID != reqDataInJson.UID {
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

			userData := edumodel.GetUserByUID(reqDataInJson.UID)
			if userData.Class != sessionClass {
				personget_replyStatus = "permission_error"
				return
			}
		}
	}

	//修改个人信息
	res := edumodel.UpdateUserByID(reqMsgInJSON.uid, reqDataInJson.Class, reqDataInJson.Name, "", reqDataInJson.Gender)
	if res {
		personput_replyStatus = "update_success"
		return
	} else {
		personput_replyStatus = "update_fail"
		return
	}
}

func (this *PersonInfoPutRouter) Handle(request eduiface.IRequest) {
	jsonMsg, err := CombineReplyMsg(personput_replyStatus, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
