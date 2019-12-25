package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

// ClassStudentDelRouter 向班级中添加学生消息路由
type ClassStudentDelRouter struct {
	edunet.BaseRouter
}

// ClassStudentDelData 向班级中添加学生消息数据结构,如果学生添加自己,则StudentListInUID为null
type ClassStudentDelData struct {
	StudentListInUID []string `json:"students"`
	ClassName        string   `json:"class"`
}

var classstudentdelReplyStatus string

// PreHandle 数据检查,权限检查,更新数据库
func (router *ClassStudentDelRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, classstudentdelReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	classstudentdelReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classstudentdelReplyStatus = "data_format_error"
		return
	}

	newMsgData := gjson.ParseBytes(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()
	sessionPlcae, err := c.GetSession("place")
	if err != nil {
		classstudentdelReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlcae.(string)
	if ok != true {
		classstudentdelReplyStatus = "session_error"
		return
	}

	//检查班级是否存在
	classNameData := newMsgData.Get("class")
	if !classNameData.Exists() {
		classstudentdelReplyStatus = "class_cannot_be_empty"
		return
	}

	className := classNameData.String()
	class := edumodel.GetClassByName(className)

	if class == nil {
		classstudentdelReplyStatus = "class_not_found"
		return
	}

	//删除学生
	if placeString == "student" {
		ok := edumodel.DeleteClassStudentByUID(className, []string{reqMsgInJSON.UID})
		if ok == true {
			classstudentdelReplyStatus = "success"
		} else {
			classstudentdelReplyStatus = "model_fail"
		}
	} else if placeString == "teacher" {
		ok := edumodel.CheckUserInClass(className, reqMsgInJSON.UID, "teacher")
		if ok != true {
			classstudentdelReplyStatus = "permission_error"
		} else {
			studentListData := newMsgData.Get("students")
			if !studentListData.Exists() {
				classstudentdelReplyStatus = "students_cannot_be_empty"
				return
			}

			studentList := studentListData.Array()
			var studentListString []string
			for _, stu := range studentList {
				studentListString = append(studentListString, stu.String())
			}

			ok := edumodel.DeleteClassStudentByUID(className, studentListString)
			if ok == true {
				classstudentdelReplyStatus = "success"
			} else {
				classstudentdelReplyStatus = "model_fail"
			}
		}
	} else if placeString == "manager" {
		studentListData := newMsgData.Get("students")
		if !studentListData.Exists() {
			classstudentdelReplyStatus = "students_cannot_be_empty"
			return
		}

		studentList := studentListData.Array()
		var studentListString []string
		for _, stu := range studentList {
			studentListString = append(studentListString, stu.String())
		}

		ok := edumodel.DeleteClassStudentByUID(className, studentListString)
		if ok == true {
			classstudentdelReplyStatus = "success"
		} else {
			classstudentdelReplyStatus = "model_fail"
		}
	}
}

// Handle 返回处理结果
func (router *ClassStudentDelRouter) Handle(request eduiface.IRequest) {
	fmt.Println("ClassStudentDelRouter: ", classstudentdelReplyStatus)
	jsonMsg, err := CombineReplyMsg(classstudentdelReplyStatus, nil)
	if err != nil {
		fmt.Println("ClassStudentDelRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}