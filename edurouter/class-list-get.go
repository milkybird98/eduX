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

// ClassListGetRouter 用于给管理员获取全部班级列表
type ClassListGetRouter struct {
	edunet.BaseRouter
}

// ClassListGetData 定义了获取班级列表请求的参数
type ClassListGetData struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

// ClassListGetReplyData 定义了返回班级列表请求时的参数
type ClassListGetReplyData struct {
	ClassList *[]edumodel.Class `json:"classlist"`
}

// 返回状态码
var classlistgetReplyStatus string

// 返回数据
var classlistgetReplyData ClassListGetReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *ClassListGetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, classlistgetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	classlistgetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classjoiningetReplyStatus = "data_format_error"
		return
	}

	// 从请求数据Data段中获取skip的limit值
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classlistgetReplyStatus = err.Error()
		return
	}

	// 如果身份不为管理员则报错返回
	if placeString != "manager" {
		classlistgetReplyStatus = "permission_error"
		return
	}

	//获取班级信息
	classList := edumodel.GetClassByOrder(int(Skip), int(Limit))
	// 如果成功获取则返回数据
	if classList != nil {
		classlistgetReplyStatus = "success"
		classlistgetReplyData.ClassList = classList
	} else {
		// 未找到合法数据
		classlistgetReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *ClassListGetRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", ClassListGetRouter: ", classlistgetReplyStatus)

	var jsonMsg []byte
	var err error

	// 生成返回数据
	if classlistgetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(classlistgetReplyStatus, classlistgetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(classlistgetReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("ClassListGetRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
