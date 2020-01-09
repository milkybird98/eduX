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

// ClassSetAlterNameRouter 处理添加课程的请求
type ClassSetAlterNameRouter struct {
	edunet.BaseRouter
}

// ClassSetAlterNameData 定义添加课程请求的参数
type ClassSetAlterNameData struct {
	ClassName string `json:"class"`
	AlterName string `json:"alter"`
}

// 添加课程请求路由的返回状态
var classsetalternameReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *ClassSetAlterNameRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, classsetalternameReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	classsetalternameReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classsetalternameReplyStatus = "data_format_error"
		return
	}

	// 解码请求数据Data段
	newClassData := gjson.ParseBytes(reqMsgInJSON.Data)
	// 试图从请求数据中获取班级数据
	classNameData := newClassData.Get("class")
	// 如果班级数据不存在则返回
	if !classNameData.Exists() || classNameData.String() == "" {
		classsetalternameReplyStatus = "classname_cannot_be_empty"
		return
	}

	className := classNameData.String()

	// 试图从请求数据中获取teacher数据
	alterData := newClassData.Get("alter")
	// 如果不存在则返回
	if !alterData.Exists() || alterData.String() == "" {
		classsetalternameReplyStatus = "altername_cannot_be_empty"
		return
	}

	alterName := alterData.String()

	//权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	sessionPlace, err := c.GetSession("place")
	// 若不存在则报错返回
	if err != nil {
		classsetalternameReplyStatus = "seesion_place_not_found"
		return
	}

	// 试图将其转换为字符串类型
	sessionPlaceString, ok := sessionPlace.(string)
	// 若转换失败则报错退出
	if !ok {
		classsetalternameReplyStatus = "session_place_data_error"
		return
	}

	// 如果身份不是管理员则报错返回
	if sessionPlaceString != "manager" && edumodel.CheckUserInClass(className, reqMsgInJSON.UID, "teacher") {
		classsetalternameReplyStatus = "permission_error"
	}

	// 检查试图存在同名班级
	class := edumodel.GetClassByName(className)
	// 如果存在则返回
	if class == nil {
		classsetalternameReplyStatus = "class_not_found"
		return
	}

	// 更新数据库
	ok = edumodel.UpdateClassAlterName(className, alterName)
	if ok == true {
		classsetalternameReplyStatus = "success"
	} else {
		classsetalternameReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *ClassSetAlterNameRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", ClassSetAlterNameRouter: ", classsetalternameReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(classsetalternameReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("ClassSetAlterNameRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
