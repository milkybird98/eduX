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

var personput_replyStatus string

func (this *PersonInfoPutRouter) PreHandle(request eduiface.IRequest) {
	var reqDataInJSON PersonInfoPutData
	reqMsgInJSON, personput_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("PersonInfoPutRouter: ", personput_replyStatus)
		return
	}

	classjoininget_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classjoininget_replyStatus = "data_format_error"
		return
	}

	newPersonInfoData := gjson.ParseBytes(reqMsgInJSON.Data)
	reqDataInJSON.UID = newPersonInfoData.Get("uid").String()
	reqDataInJSON.Name = newPersonInfoData.Get("name").String()
	reqDataInJSON.ClassName = newPersonInfoData.Get("class").String()
	reqDataInJSON.Gender = int(newPersonInfoData.Get("gender").Int())

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		personget_replyStatus = "session_error"
		return
	}

	userData := edumodel.GetUserByUID(reqDataInJSON.UID)
	if userData == nil {
		personget_replyStatus = "user_not_found"
		return
	}

	if sessionUID != reqDataInJSON.UID {
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
			if userData.Class != sessionClass {
				personget_replyStatus = "permission_error"
				return
			}
		}
	}

	//修改个人信息
	res := edumodel.UpdateUserByID(reqMsgInJSON.UID, reqDataInJSON.ClassName, reqDataInJSON.Name, "", reqDataInJSON.Gender)
	if res {
		personput_replyStatus = "update_success"
		return
	} else {
		personput_replyStatus = "update_fail"
		return
	}
}

func (this *PersonInfoPutRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonInfoPutRouter: ", personput_replyStatus)
	jsonMsg, err := CombineReplyMsg(personput_replyStatus, nil)
	if err != nil {
		fmt.Println("PersonInfoPutRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
