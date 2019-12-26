package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"
)

type ClassJoinInGetRouter struct {
	edunet.BaseRouter
}

type ClassJoinInGetReplyData struct {
	ClassName   string
	TeacherList []string
	StudentList []string
}

var classjoiningetReplyStatus string
var classjoiningetReplyData ClassJoinInGetReplyData

func (router *ClassJoinInGetRouter) PreHandle(request eduiface.IRequest) {
	classjoiningetReplyData = ClassJoinInGetReplyData{}

	reqMsgInJSON, classjoiningetReplyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("ClassJoinInGetRouter: ", classjoiningetReplyStatus)
		return
	}

	classjoiningetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		fmt.Println("ClassJoinInGetRouter: ", classjoiningetReplyStatus)
		return
	}

	c := request.GetConnection()
	place, err := c.GetSession("place")
	if err != nil {
		classjoiningetReplyStatus = "session_error"
		fmt.Println("ClassJoinInGetRouter: ", classjoiningetReplyStatus)
		return
	}

	placeString, ok := place.(string)
	if ok != true {
		classjoiningetReplyStatus = "session_error"
		fmt.Println("ClassJoinInGetRouter: ", classjoiningetReplyStatus)
		return
	}

	class := edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)

	if class == nil {
		classjoiningetReplyStatus = "not_join_class"
	} else {
		classjoiningetReplyStatus = "success"
		classjoiningetReplyData.ClassName = class.ClassName
		classjoiningetReplyData.StudentList = class.StudentList
		classjoiningetReplyData.TeacherList = class.TeacherList
	}
}

func (router *ClassJoinInGetRouter) Handle(request eduiface.IRequest) {
	var jsonMsg []byte
	var err error
	if classjoiningetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(classjoiningetReplyStatus, classjoiningetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(classjoiningetReplyStatus, nil)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
