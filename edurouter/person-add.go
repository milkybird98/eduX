package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"encoding/base64"
	"fmt"

	"github.com/tidwall/gjson"
)

type PersonAddRouter struct {
	edunet.BaseRouter
}

type PersonAddData struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Place string `json:"place"`
}

var personAddReplyStatus string

func (router *PersonAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	var reqDataInJSON PersonAddData
	reqMsgInJSON, personAddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	personAddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personAddReplyStatus = "data_format_error"
		return
	}

	newStudentData := gjson.ParseBytes(reqMsgInJSON.Data)
	reqDataInJSON.UID = newStudentData.Get("uid").String()
	reqDataInJSON.Name = newStudentData.Get("name").String()
	reqDataInJSON.Place = newStudentData.Get("place").String()

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		personAddReplyStatus = "seesion_error"
		return
	}

	if sessionPlace != "manager" {
		personAddReplyStatus = "permission_error"
		return
	}

	//数据库操作
	var newUser edumodel.User

	newUser.UID = reqDataInJSON.UID
	newUser.Name = reqDataInJSON.Name
	newUser.Pwd = base64.StdEncoding.EncodeToString([]byte(newUser.UID))
	newUser.Place = reqDataInJSON.Place
	newUser.Class = ""
	newUser.Gender = 0

	res := edumodel.AddUser(&newUser)
	if res {
		personAddReplyStatus = "success"
	} else {
		personAddReplyStatus = "model_fail"
	}

}

func (router *PersonAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonAddRouter: ", personAddReplyStatus)
	jsonMsg, err := CombineReplyMsg(personAddReplyStatus, nil)
	if err != nil {
		fmt.Println("PersonAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}