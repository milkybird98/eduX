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

// PersonCountRouter 用于处理管理员统计文件数量的请求
type PersonCountRouter struct {
	edunet.BaseRouter
}

// PersonCountReplyData 定义了管理员统计文件数量请求的返回数据Data段参数
type PersonCountReplyData struct {
	Number int `json:"num"`
}

// 返回状态码
var personcountReplyStatus string

// 返回数据
var personcountReplyData PersonCountReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonCountRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, personcountReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	personcountReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personcountReplyStatus = "data_format_error"
		return
	}

	// 权限检查
	/*
		c := request.GetConnection()

		// 试图从session中获取身份数据
		userPlace, err := GetSessionPlace(c)
		if err != nil {
			personcountReplyStatus = err.Error()
			return
		}

		// 如果身份不为管理员则权限错误
		if userPlace != "manager" {
			personcountReplyStatus = "permission_error"
			return
		}
	*/

	personcountReplyData.Number = edumodel.GetUserNumber()

	// 如果查询成功在返回success,否则返回错误码
	if personcountReplyData.Number != -1 {
		personcountReplyStatus = "success"
	} else {
		personcountReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonCountRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonCountRouter: ", personcountReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if personcountReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(personcountReplyStatus, personcountReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(personcountReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PersonCountRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
