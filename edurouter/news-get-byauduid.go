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

// NewsGetByAudientUIDRouter 处理通过听众UID获取消息的请求
type NewsGetByAudientUIDRouter struct {
	edunet.BaseRouter
}

// NewsGetByAudientUIDData 定义根据听众UID获取消息请求的参数
type NewsGetByAudientUIDData struct {
	Audient    string `json:"audient,omitempty"` // 听众UID,供管理员请求时使用
	Skip       int64  `json:"skip"`              // 跳过项目数
	Limit      int64  `json:"limit"`             // 获取项目数
	IsAnnounce bool   `json:"isannounce"`        // 是否是公告
}

// NewGetByAudientUIDReplyData 定义返回消息的参数格式
type NewGetByAudientUIDReplyData struct {
	NewsList *[]edumodel.News `json:"news"`
}

// 返回状态码
var newgetbyaudientuidReplyStatus string

// 返回数据
var newgetbyaudientuidReplyData NewGetByAudientUIDReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *NewsGetByAudientUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, newgetbyaudientuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	newgetbyaudientuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newgetbyaudientuidReplyStatus = "data_format_error"
		return
	}

	// 获取skip和limit数据
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
		newgetbyaudientuidReplyStatus = err.Error()
		return
	}

	// 如果当前连接用户是教师或者学生
	if placeString != "manager" {
		// 查询数据库,查询自己是听众的数据
		newsList := edumodel.GetNewsByAudientUID(int(Skip), int(Limit), IsAnnounce, reqMsgInJSON.UID)
		// 如果数据存在则返回success和数据,否则返回错误码
		if newsList != nil {
			newgetbyaudientuidReplyStatus = "success"
			newgetbyaudientuidReplyData.NewsList = newsList
		} else {
			newgetbyaudientuidReplyStatus = "model_fail"
		}
	} else { // 如果当前连接用户是管理员
		// 获取听众数据
		audientUIDData := gjson.GetBytes(reqMsgInJSON.Data, "audient")
		// 若听众数据不存在则报错返回
		if !audientUIDData.Exists() || audientUIDData.String() == "" {
			newgetbyaudientuidReplyStatus = "audientuid_cannot_be_empty"
			return
		}
		audientUID := audientUIDData.String()

		// 根据听众数据查询消息
		newsList := edumodel.GetNewsByAudientUID(int(Skip), int(Limit), IsAnnounce, audientUID)
		// 如果数据存在则返回success和数据,否则返回错误码
		if newsList != nil {
			newgetbyaudientuidReplyStatus = "success"
			newgetbyaudientuidReplyData.NewsList = newsList
		} else {
			newgetbyaudientuidReplyStatus = "model_fail"
		}
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *NewsGetByAudientUIDRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsGetByAudientUIDRouter: ", newgetbyaudientuidReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if newgetbyaudientuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newgetbyaudientuidReplyStatus, newgetbyaudientuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(newgetbyaudientuidReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("NewsGetByAudientUIDRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
