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

var studentadd_replyStatus string

func (this *StudentAddRouter) PreHandle(request eduiface.IRequest) {
	var reqDataInJSON StudentAddData
	reqMsgInJSON, studentadd_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("StudentAddRouter: ", studentadd_replyStatus)
		return
	}

	studentadd_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		studentadd_replyStatus = "data_format_error"
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
		studentadd_replyStatus = "seesion_error"
		return
	}

	if sessionPlace != "manager" {
		studentadd_replyStatus = "permission_error"
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
		studentadd_replyStatus = "add_success"
		return
	}else{
		studentadd_replyStatus = "add_fail"
		return
	}

}

func (this *StudentAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("StudentAddRouter: ", studentadd_replyStatus)
	jsonMsg, err := CombineReplyMsg(studentadd_replyStatus, nil)
	if err != nil {
		fmt.Println("StudentAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
