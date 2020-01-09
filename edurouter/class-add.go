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

// ClassAddRouter 处理添加课程的请求
type ClassAddRouter struct {
	edunet.BaseRouter
}

// ClassAddData 定义添加课程请求的参数
type ClassAddData struct {
	ClassName  string `json:"class"`
	AlterName  string `json:"alter"`
	TeacherUID string `json:"teacher"`
}

// 添加课程请求路由的返回状态
var classaddReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *ClassAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, classaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	classaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classaddReplyStatus = "data_format_error"
		return
	}

	// 解码请求数据Data段
	newClassData := gjson.ParseBytes(reqMsgInJSON.Data)
	// 试图从请求数据中获取班级数据
	classNameData := newClassData.Get("class")
	// 如果班级数据不存在则返回
	if !classNameData.Exists() || classNameData.String() == "" {
		classaddReplyStatus = "classname_cannot_be_empty"
		return
	}

	className := classNameData.String()

	// 试图从请求数据中获取teacher数据
	teacherUIDData := newClassData.Get("teacher")
	// 如果不存在则返回
	if !teacherUIDData.Exists() || teacherUIDData.String() == "" {
		classaddReplyStatus = "init_teacher_cannot_be_empty"
		return
	}

	teacherUID := teacherUIDData.String()

	//权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	sessionPlace, err := c.GetSession("place")
	// 若不存在则报错返回
	if err != nil {
		classaddReplyStatus = "seesion_place_not_found"
		return
	}

	// 试图将其转换为字符串类型
	sessionPlaceString, ok := sessionPlace.(string)
	// 若转换失败则报错退出
	if !ok {
		classaddReplyStatus = "session_place_data_error"
		return
	}

	// 如果身份不是管理员则报错返回
	if sessionPlaceString != "manager" {
		classaddReplyStatus = "permission_error"
	}

	// 检查试图存在同名班级
	class := edumodel.GetClassByName(className)
	// 如果存在则返回
	if class != nil {
		classaddReplyStatus = "same_class_exist"
		return
	}

	// 检查教师uid对应教师是否存在
	teacher := edumodel.GetUserByUID(teacherUID)
	// 若不存在则返回
	if teacher == nil {
		classaddReplyStatus = "teacher_not_found"
		return
	}

	// 构建新的班级数据
	var newClass edumodel.Class
	newClass.ClassName = className
	newClass.AlterName = newClassData.Get("alter").String()
	newClass.TeacherList = []string{teacherUID}
	newClass.StudentList = []string{}
	newClass.CreateDate = time.Now()

	// 更新数据库
	ok = edumodel.AddClass(&newClass) && edumodel.AddUserToClassByUID([]string{teacherUID}, className)
	if ok == true {
		classaddReplyStatus = "success"
	} else {
		classaddReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *ClassAddRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", ClassAddRouter: ", classaddReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(classaddReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("ClassAddRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
