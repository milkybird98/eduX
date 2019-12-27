package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

type ClassAddRouter struct {
	edunet.BaseRouter
}

type ClassAddData struct {
	ClassName  string `json:"class"`
	TeacherUID string `json:"teacher"`
}

var classaddReplyStatus string

func (router *ClassAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, classaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	classaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classaddReplyStatus = "data_format_error"
		return
	}

	newClassData := gjson.ParseBytes(reqMsgInJSON.Data)
	classNameData := newClassData.Get("class")
	if !classNameData.Exists() {
		classaddReplyStatus = "classname_cannot_be_empty"
		return
	}

	className := classNameData.String()

	teacherUIDData := newClassData.Get("teacher")
	if !teacherUIDData.Exists() {
		classaddReplyStatus = "init_teacher_cannot_be_empty"
		return
	}

	teacherUID := teacherUIDData.String()

	//权限检查
	c := request.GetConnection()

	sessionPlace, err := c.GetSession("place")
	if err != nil {
		classaddReplyStatus = "seesion_place_not_found"
		return
	}

	sessionPlaceString, ok := sessionPlace.(string)
	if !ok {
		classaddReplyStatus = "session_place_data_error"
		return
	}

	if sessionPlaceString != "manager" {
		classaddReplyStatus = "permission_error"
	}

	class := edumodel.GetClassByName(className)
	if class != nil {
		classaddReplyStatus = "same_class_exist"
		return
	}

	teacher := edumodel.GetUserByUID(teacherUID)
	if teacher == nil {
		classaddReplyStatus = "teacher_not_found"
		return
	}

	newClass := edumodel.Class{className, []string{}, []string{teacherUID}}
	ok = edumodel.AddClass(&newClass)
	if ok == true {
		classaddReplyStatus = "success"
	} else {
		classaddReplyStatus = "model_fail"
	}
}

func (router *ClassAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("ClassAddRouter: ", classaddReplyStatus)
	jsonMsg, err := CombineReplyMsg(classaddReplyStatus, nil)
	if err != nil {
		fmt.Println("ClassAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
