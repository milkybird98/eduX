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

// NewsGetBySenderUIDRouter 处理根据发送者UID查询消息的请求
type NewsGetBySenderUIDRouter struct {
	edunet.BaseRouter
}

// NewsGetBySenderUIDData 定义根据发送者UID查询消息的参数
type NewsGetBySenderUIDData struct {
	Sender   string `json:"sender,omitempty"`
	Skip     int64  `json:"skip"`
	Limit    int64  `json:"limit"`
	NewsType int64  `json:"type"`
}

// NewsGetBySenderUIDReplyData 定义根据发送者UID参数请求消息的返回参数
type NewsGetBySenderUIDReplyData struct {
	NewsList *[]edumodel.News `json:"news"`
}

// 返回状态码
var newsgetbysenderuidReplyStatus string

// 返回数据
var newsgetbysenderuidReplyData NewsGetBySenderUIDReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *NewsGetBySenderUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并验证校验和
	reqMsgInJSON, newsgetbysenderuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接用户是否已登录
	newsgetbysenderuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newsgetbysenderuidReplyStatus = "data_format_error"
		return
	}

	// 获取skip和limit数据
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	// 从Data段获取公告标志位,判断是否是公告
	newsTypeData := gjson.GetBytes(reqMsgInJSON.Data, "type")
	// 如果不存在,则认为默认是非公告
	if !newsTypeData.Exists() || newsTypeData.Int() < 1 || newsTypeData.Int() > 5 {
		newsgetbysenderuidReplyStatus = "type_cannot_be_empty"
		return
	}
	newsType := newsTypeData.Int()

	//权限检查

	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		newsgetbysenderuidReplyStatus = err.Error()
		return
	}

	// 如果当前连接用户是教师
	if placeString != "manager" {
		// 查询数据库,查询自己发送的消息
		newsList := edumodel.GetNewsBySenderUID(int(Skip), int(Limit), newsType, reqMsgInJSON.UID)
		// 如果数据存在则返回success和数据,否则返回错误码
		if newsList != nil {
			newsgetbysenderuidReplyStatus = "success"
			newsgetbysenderuidReplyData.NewsList = newsList
		} else {
			newsgetbysenderuidReplyStatus = "model_fail"
		}
	} else { // 如果当前连接用户是管理员
		// 获取发送者数据
		senderUIDData := gjson.GetBytes(reqMsgInJSON.Data, "sender")
		// 若发送者数据不存在则报错返回
		if !senderUIDData.Exists() || senderUIDData.String() == "" {
			newsgetbysenderuidReplyStatus = "senderuid_cannot_be_empty"
			return
		}
		senderUID := senderUIDData.String()

		// 根据发送者查询消息
		newsList := edumodel.GetNewsBySenderUID(int(Skip), int(Limit), newsType, senderUID)
		// 如果数据存在则返回success和数据,否则返回错误码
		if newsList != nil {
			newsgetbysenderuidReplyStatus = "success"
			newsgetbysenderuidReplyData.NewsList = newsList
		} else {
			newsgetbysenderuidReplyStatus = "model_fail"
		}
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *NewsGetBySenderUIDRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsGetBySenderUIDRouter: ", newsgetbysenderuidReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if newsgetbysenderuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newsgetbysenderuidReplyStatus, newsgetbysenderuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(newsgetbysenderuidReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("NewsGetBySenderUIDRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
