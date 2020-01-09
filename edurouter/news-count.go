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

// NewsCountRouter 用于处理管理员统计文件数量的请求
type NewsCountRouter struct {
	edunet.BaseRouter
}

// NewsCountData 定义了管理员请求文件数量时Data段的参数
type NewsCountData struct {
	AudientUID string `json:"audiuid"`
	SendUID    string `json:"senduid"`
	NewsType   string `json:"type"`
}

// NewsCountReplyData 定义了管理员统计文件数量请求的返回数据Data段参数
type NewsCountReplyData struct {
	Number int `json:"num"`
}

// 返回状态码
var newscountReplyStatus string

// 返回数据
var newscountReplyData NewsCountReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *NewsCountRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, newscountReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	newscountReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newscountReplyStatus = "data_format_error"
		return
	}

	// 权限检查
	/*
		c := request.GetConnection()

		// 试图从session中获取身份数据
		userPlace, err := GetSessionPlace(c)
		if err != nil {
			newscountReplyStatus = err.Error()
			return
		}

		// 如果身份不为管理员则权限错误
		if userPlace != "manager" {
			newscountReplyStatus = "permission_error"
			return
		}
	*/

	// 尝试从Data段中获取班级名称
	audiUID := gjson.GetBytes(reqMsgInJSON.Data, "audiuid").String()

	sendUID := gjson.GetBytes(reqMsgInJSON.Data, "senduid").String()

	newsType := gjson.GetBytes(reqMsgInJSON.Data, "type").Int()

	newscountReplyData.Number = -1

	if newsType == 4 {
		audiUID = "all"
	}

	newscountReplyData.Number = edumodel.GetNewsNumber(audiUID, sendUID, int(newsType))

	// 如果查询成功在返回success,否则返回错误码
	if newscountReplyData.Number != -1 {
		newscountReplyStatus = "success"
	} else {
		newscountReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *NewsCountRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsCountRouter: ", newscountReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if newscountReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newscountReplyStatus, newscountReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(newscountReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("NewsCountRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
