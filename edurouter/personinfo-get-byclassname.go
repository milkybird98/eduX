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

/*
 *	MsgID 101
 *
 *
 *
 *
 */

var persongetbyclass_replyStatus string
var persongetbyclass_replyData []PersonInfoGetReplyData

func (this *PersonInfoGetByClassRouter) PreHandle(request eduiface.IRequest) {
	reqMsgInJSON, persongetbyclass_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("PersonInfoGetByClassRouter: ", persongetbyclass_replyStatus)
		return
	}

	persongetbyclass_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		return
	}

	reqClassNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	if !reqClassNameData.Exists() {
		persongetbyclass_replyStatus = "data_format_error"
		return
	}

	reqClassName := reqClassNameData.String()

	c := request.GetConnection()
	sessionPlcae, err := c.GetSession("Plcae")
	if err != nil {
		persongetbyclass_replyStatus = "session_error"
		return
	}
	if sessionPlcae == "student" {
		persongetbyclass_replyStatus = "permission_error"
		return
	} else if sessionPlcae == "teacher" {
		sessionClass, err := c.GetSession("Class")
		if err != nil {
			persongetbyclass_replyStatus = "session_error"
			return
		}
		if reqClassName != sessionClass {
			persongetbyclass_replyStatus = "permission_error"
			return
		}
	}

	userManyData := edumodel.GetUserByClass(reqClassName)
	if userManyData == nil || len(*userManyData) <= 0 {
		persongetbyclass_replyStatus = "data_not_found"
		return
	}

	persongetbyclass_replyStatus = "success"

	for _, personData := range *userManyData {
		persongetbyclass_replyData = append(
			persongetbyclass_replyData,
			PersonInfoGetReplyData{
				personData.UID,
				personData.Name,
				personData.Class,
				personData.Gender})
	}
}

func (this *PersonInfoGetByClassRouter) Handle(request eduiface.IRequest) {
	fmt.Println("PersonInfoGetByClassRouter: ", persongetbyclass_replyStatus)
	jsonMsg, err := CombineReplyMsg(persongetbyclass_replyStatus, persongetbyclass_replyData)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
