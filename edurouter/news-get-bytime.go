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

// NewsGetByTimeOrderRouter 处理按照时间顺序获取消息的请求
type NewsGetByTimeOrderRouter struct {
	edunet.BaseRouter
}

// NewsGetByTimeOrderData 定义按照时间顺序请求消息时的参数
type NewsGetByTimeOrderData struct {
	Skip       int64 `json:"skip"`
	Limit      int64 `json:"limit"`
	IsAnnounce bool  `json:"isannounce"`
}

// NewGetByTimeOrderReplyData 定义根据时间顺序请求消息时的返回参数
type NewGetByTimeOrderReplyData struct {
	NewsList *[]edumodel.News `json:"news"`
}

// 返回状态码
var newgetbytimeorderReplyStatus string

// 返回数据
var newgetbytimeorderReplyData NewGetByTimeOrderReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *NewsGetByTimeOrderRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, newgetbytimeorderReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	newgetbytimeorderReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newgetbytimeorderReplyStatus = "data_format_error"
		return
	}

	// 获取skip和limit值
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	var IsAnnounce bool
	// 从Data段获取公告限定标志位
	isAnnounceData := gjson.GetBytes(reqMsgInJSON.Data, "isannounce")
	// 如果标志位不存在则默认认为是非公告
	if isAnnounceData.Exists() {
		IsAnnounce = isAnnounceData.Bool()
	} else {
		IsAnnounce = false
	}

	//权限检查

	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 如果当前连接用户不是管理员则权限错误
	if placeString != "manager" {
		newgetbytimeorderReplyStatus = "permission_error"
		return
	}

	// 查询数据库
	newsList := edumodel.GetNewsByTimeOrder(int(Skip), int(Limit), IsAnnounce)
	// 如果数据存在则返回success和数据,否则返回错误码
	if newsList != nil {
		newgetbytimeorderReplyStatus = "success"
		newgetbytimeorderReplyData.NewsList = newsList
	} else {
		newgetbytimeorderReplyStatus = "model_fail"
	}

}

// Handle 用于将请求的处理结果发回客户端
func (router *NewsGetByTimeOrderRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] Time: ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsGetByTimeOrderRouter: ", newgetbytimeorderReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if newgetbytimeorderReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newgetbytimeorderReplyStatus, newgetbytimeorderReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(newgetbytimeorderReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("NewsGetByTimeOrderRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
