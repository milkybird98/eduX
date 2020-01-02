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

// ClassDelRouter 处理课程删除请求
type ClassDelRouter struct {
	edunet.BaseRouter
}

// ClassDelData 定义课程删除请求参数
type ClassDelData struct {
	ClassName string `json:"class"`
}

// 课程删除路由的返回状态
var classdelReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *ClassDelRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, classdelReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	classdelReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classdelReplyStatus = "data_format_error"
		return
	}

	// 试图从请求数据Data段获取class数据
	delClassData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	// 如果数据不存在则返回
	if !delClassData.Exists() || delClassData.String() == "" {
		classdelReplyStatus = "data_format_error"
		return
	}

	delClassName := delClassData.String()

	//权限检查

	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 如果身份不为管理员则报错返回
	if placeString != "manager" {
		classdelReplyStatus = "permission_error"
		return
	}

	// 数据库操作

	// 检查班级是否存在
	class := edumodel.GetClassByName(delClassName)
	// 如果不存在则返回
	if class == nil {
		classdelReplyStatus = "class_not_found"
		return
	}

	// 删除班级
	ok = edumodel.DeleteClassByName(delClassName)
	if ok == true {
		classdelReplyStatus = "success"
	} else {
		classdelReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *ClassDelRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", ClassDelRouter: ", classdelReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(classdelReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("ClassDelRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
