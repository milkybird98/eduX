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

// ClassStudentAddRouter 向班级中添加学生消息路由
type ClassStudentAddRouter struct {
	edunet.BaseRouter
}

// ClassStudentAddData 向班级中添加学生消息数据结构,如果学生添加自己,则StudentListInUID为null
type ClassStudentAddData struct {
	StudentListInUID []string `json:"students"`
	ClassName        string   `json:"class"`
}

// 返回数据
var classstudentaddReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据验证和数据库操作
func (router *ClassStudentAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, classstudentaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	classstudentaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classstudentaddReplyStatus = "data_format_error"
		return
	}

	// 解码请求数据中的Data段
	newMsgData := gjson.ParseBytes(reqMsgInJSON.Data)

	//权限检查

	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classstudentaddReplyStatus = err.Error()
		return
	}

	// 检查Data段中是否存在班级名称
	classNameData := newMsgData.Get("class")
	// 若不存在则返回
	if !classNameData.Exists() || classNameData.String() == "" {
		classstudentaddReplyStatus = "class_cannot_be_empty"
		return
	}

	// 根据班级名称从数据库中获取班级数据
	className := classNameData.String()
	class := edumodel.GetClassByName(className)

	// 如果班级不存在则返回
	if class == nil {
		classstudentaddReplyStatus = "class_not_found"
		return
	}

	// 如果当前连接用户的身份为学生
	if placeString == "student" {
		// 检查是否已经加入班级
		class := edumodel.GetClassByUID(reqMsgInJSON.UID, "student")
		// 如果已经加入则返回错误码
		if class != nil {
			classstudentaddReplyStatus = "already_in_class"
			return
		}
		// 如果当前连接用户的身份为教师
	} else if placeString == "teacher" {
		// 检查自己是否在该班级中
		ok := edumodel.CheckUserInClass(className, reqMsgInJSON.UID, "teacher")
		// 若不在则出现权限错误
		if ok != true {
			classstudentaddReplyStatus = "permission_error"
			return
		}
	}

	//数据验证并更新数据库

	// 如果当前连接用户的身份为学生
	if placeString == "student" {
		// 则将自己加入该班级,并更新班级数据和用户数据
		ok := edumodel.UpdateClassStudentByUID(className, []string{reqMsgInJSON.UID}) && edumodel.AddUserToClassByUID([]string{reqMsgInJSON.UID}, className)
		// 若成功则返回success否则返回错误码
		if ok == true {
			classstudentaddReplyStatus = "success"
		} else {
			classstudentaddReplyStatus = "model_fail"
		}
		// 如果当前连接用户的身份为教师或管理员
	} else if placeString == "teacher" || placeString == "manager" {
		// 试图从Data段获取学生数据
		studentListData := newMsgData.Get("students")
		// 如果学生列表为空则返回错误码
		if !studentListData.Exists() || !studentListData.IsArray() || len(studentListData.Array()) == 0 {
			classstudentaddReplyStatus = "students_cannot_be_empty"
			return
		}

		// 获取并检查学生列表中的每一项
		studentList := studentListData.Array()
		var studentListString []string
		for _, stu := range studentList {
			// 如果uid对应学生不存在则跳过
			student := edumodel.GetUserByUID(stu.String())
			if student == nil {
				continue
			}
			// 将学生uid添加进列表中
			studentListString = append(studentListString, stu.String())
		}

		// 则将学生列表加入该班级,并更新班级数据和用户数据
		ok := edumodel.UpdateClassStudentByUID(className, studentListString) && edumodel.AddUserToClassByUID(studentListString, className)
		// 如果成功则返回success,否则返回错误码
		if ok == true {
			classstudentaddReplyStatus = "success"
		} else {
			classstudentaddReplyStatus = "model_fail"
		}
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *ClassStudentAddRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", ClassStudentAddRouter: ", classstudentaddReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(classstudentaddReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("ClassStudentAddRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
