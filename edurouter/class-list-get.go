package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

type ClassListGetRouter struct {
	edunet.BaseRouter
}

type ClassListGetData struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type ClassListGetReplyData struct {
	ClassList *[]edumodel.Class `json:"classlist"`
}

var classlistgetReplyStatus string
var classlistgetReplyData ClassListGetReplyData

func (router *ClassListGetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, classlistgetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	classlistgetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classjoiningetReplyStatus = "data_format_error"
		return
	}

	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		classlistgetReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		classlistgetReplyStatus = "session_error"
		return
	}

	if placeString != "manager" {
		classlistgetReplyStatus = "permission_error"
		return
	}

	//获取班级信息
	classList := edumodel.GetClassByOrder(int(Skip), int(Limit))
	if classList != nil {
		classlistgetReplyStatus = "success"
		classlistgetReplyData.ClassList = classList
	} else {
		classlistgetReplyStatus = "model_fail"
	}
}

func (router *ClassListGetRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", ClassListGetRouter: ", classlistgetReplyStatus)

	var jsonMsg []byte
	var err error

	if classlistgetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(classlistgetReplyStatus, classlistgetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(classlistgetReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("ClassListGetRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
