package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

// ClassStudentAddRouter 向班级中添加学生消息路由
type ClassStudentAddRouter struct {
	edunet.BaseRouter
}

// ClassStudentAddData 向班级中添加学生消息数据结构,如果学生添加自己,则StudentListInUID为null
type ClassStudentAddData struct {
	StudentListInUID []string `json:"students"`
	ClassName        string   `json:"class"`
}

var classstudentaddReplyStatus string

// PreHandle 数据检查,权限检查,更新数据库
func (router *ClassStudentAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, classstudentaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	classstudentaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classstudentaddReplyStatus = "data_format_error"
		return
	}

	newMsgData := gjson.ParseBytes(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		classstudentaddReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		classstudentaddReplyStatus = "session_error"
		return
	}

	//检查班级是否存在
	classNameData := newMsgData.Get("class")
	if !classNameData.Exists() {
		classstudentaddReplyStatus = "class_cannot_be_empty"
		return
	}

	className := classNameData.String()
	class := edumodel.GetClassByName(className)

	if class == nil {
		classstudentaddReplyStatus = "class_not_found"
		return
	}

	//添加学生
	if placeString == "student" {
		ok := edumodel.UpdateClassStudentByUID(className, []string{reqMsgInJSON.UID}) && edumodel.AddUserToClassByUID([]string{reqMsgInJSON.UID}, className)
		if ok == true {
			classstudentaddReplyStatus = "success"
		} else {
			classstudentaddReplyStatus = "model_fail"
		}

	} else if placeString == "teacher" {
		ok := edumodel.CheckUserInClass(className, reqMsgInJSON.UID, "teacher")
		if ok != true {
			classstudentaddReplyStatus = "permission_error"
		} else {
			studentListData := newMsgData.Get("students")
			if !studentListData.Exists() {
				classstudentaddReplyStatus = "students_cannot_be_empty"
				return
			}

			studentList := studentListData.Array()
			var studentListString []string
			for _, stu := range studentList {
				student := edumodel.GetUserByUID(stu.String())
				if student == nil {
					continue
				}
				studentListString = append(studentListString, stu.String())
			}

			ok := edumodel.UpdateClassStudentByUID(className, studentListString) && edumodel.AddUserToClassByUID(studentListString, className)
			if ok == true {
				classstudentaddReplyStatus = "success"
			} else {
				classstudentaddReplyStatus = "model_fail"
			}
		}
	} else if placeString == "manager" {
		studentListData := newMsgData.Get("students")
		if !studentListData.Exists() {
			classstudentaddReplyStatus = "students_cannot_be_empty"
			return
		}

		studentList := studentListData.Array()
		var studentListString []string
		for _, stu := range studentList {
			student := edumodel.GetUserByUID(stu.String())
			if student == nil {
				continue
			}
			studentListString = append(studentListString, stu.String())
		}

		ok := edumodel.UpdateClassStudentByUID(className, studentListString) && edumodel.AddUserToClassByUID(studentListString, className)
		if ok == true {
			classstudentaddReplyStatus = "success"
		} else {
			classstudentaddReplyStatus = "model_fail"
		}
	}
}

// Handle 返回处理结果
func (router *ClassStudentAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("ClassStudentAddRouter: ", classstudentaddReplyStatus)
	jsonMsg, err := CombineReplyMsg(classstudentaddReplyStatus, nil)
	if err != nil {
		fmt.Println("ClassStudentAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
