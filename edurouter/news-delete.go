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

// NewsDeleteRouter 负责问题删除业务的路由
type NewsDeleteRouter struct {
	edunet.BaseRouter
}

// NewsDeleteData 用于客户端指定需要删除的问题
type NewsDeleteData struct {
	NewsID string `json:"id"`
}

// 返回状态码
var newdeleteReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *NewsDeleteRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验码是否正确
	reqMsgInJSON, newdeleteReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	newdeleteReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newdeleteReplyStatus = "data_format_error"
		return
	}

	// 试图从Data段获取消息id数据
	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	// 若不存在则返回错误码
	if !innerIDData.Exists() || innerIDData.String() == "" {
		newdeleteReplyStatus = "newsid_cannot_be_empty"
		return
	}
	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 检查id对应的消息数据是否存在
	newsData := edumodel.GetNewsByInnerID(innerIDString)
	// 如果消息不存在则报错
	if newsData == nil {
		newdeleteReplyStatus = "news_not_found"
		return
	}

	// 如果当前用户不是管理员或者当前用户不是消息的发送者则报错
	if placeString != "manager" && reqMsgInJSON.UID != newsData.SenderUID {
		newdeleteReplyStatus = "permission_error"
		return
	}

	// 从数据库删除消息
	ok = edumodel.DeleteNewsByInnerID(innerIDString)
	if ok {
		newdeleteReplyStatus = "success"
	} else {
		newdeleteReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *NewsDeleteRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsDeleteRouter: ", newdeleteReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(newdeleteReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("NewsDeleteRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
