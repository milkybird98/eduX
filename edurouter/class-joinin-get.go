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

var classjoininget_replyStatus string
var classjoininget_replyData ClassJoinInGetReplyData

func (this *ClassJoinInGetRouter) PreHandle(request eduiface.IRequest) {
	classjoininget_replyData = ClassJoinInGetReplyData{}

	reqMsgInJSON, classjoininget_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("ClassJoinInGetRouter: ", classjoininget_replyStatus)
		return
	}

	classjoininget_replyStatus, ok = CheckConnectionLogin(request)
	if ok != true {
		fmt.Println("ClassJoinInGetRouter: ", classjoininget_replyStatus)
		return
	}

	c := request.GetConnection()
	place, err := c.GetSession("plcae")
	if err != nil {
		classjoininget_replyStatus = "session_error"
		fmt.Println("ClassJoinInGetRouter: ", classjoininget_replyStatus)
		return
	}

	placeString, ok := place.(string)
	if ok != true {
		classjoininget_replyStatus = "session_error"
		fmt.Println("ClassJoinInGetRouter: ", classjoininget_replyStatus)
		return
	}

	class := edumodel.GetClassByUID(reqMsgInJSON.uid, placeString)

	classjoininget_replyData.ClassName = class.ClassName
	classjoininget_replyData.StudentList = class.StudentList
	classjoininget_replyData.TeacherList = class.TeacherList
}

func (this *ClassJoinInGetRouter) Handle(request eduiface.IRequest) {
	jsonMsg, err := CombineReplyMsg(classjoininget_replyStatus, classjoininget_replyData)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
