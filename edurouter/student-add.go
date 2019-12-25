package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"encoding/base64"
	"fmt"

	"github.com/tidwall/gjson"
)

type StudentAddRouter struct {
	edunet.BaseRouter
}

type StudentAddData struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	ClassName string `json:"class"`
}

var studentaddReplyStatus string

func (router *StudentAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	var reqDataInJSON StudentAddData
	reqMsgInJSON, studentaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	studentaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		studentaddReplyStatus = "data_format_error"
		return
	}

	newStudentData := gjson.ParseBytes(reqMsgInJSON.Data)
	reqDataInJSON.UID = newStudentData.Get("uid").String()
	reqDataInJSON.Name = newStudentData.Get("name").String()
	reqDataInJSON.ClassName = newStudentData.Get("class").String()

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		studentaddReplyStatus = "seesion_error"
		return
	}

	if sessionPlace != "manager" {
		studentaddReplyStatus = "permission_error"
	}

	//数据库操作
	var newUser edumodel.User

	newUser.UID = reqDataInJSON.UID
	newUser.Name = reqDataInJSON.Name
	newUser.Pwd = base64.StdEncoding.EncodeToString([]byte(newUser.UID))
	newUser.Plcae = "student"
	newUser.Class = reqDataInJSON.ClassName
	newUser.Gender = 0

	res := edumodel.AddUser(&newUser)
	if res {
		studentaddReplyStatus = "success"
	} else {
		studentaddReplyStatus = "model_fail"
	}

}

func (router *StudentAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("StudentAddRouter: ", studentaddReplyStatus)
	jsonMsg, err := CombineReplyMsg(studentaddReplyStatus, nil)
	if err != nil {
		fmt.Println("StudentAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
