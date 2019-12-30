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

type NewsGetBySenderUIDRouter struct {
	edunet.BaseRouter
}

type NewsGetBySenderUIDData struct {
	Sender     string `json:"sender,omitempty"`
	Skip       int64  `json:"skip"`
	Limit      int64  `json:"limit"`
	IsAnnounce bool   `json:"isannounce"`
}

type NewsGetBySenderUIDReplyData struct {
	NewsList *[]edumodel.News `json:"news"`
}

var newsgetbysenderuidReplyStatus string
var newsgetbysenderuidReplyData NewsGetBySenderUIDReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *NewsGetBySenderUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, newsgetbysenderuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	newsgetbysenderuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newsgetbysenderuidReplyStatus = "data_format_error"
		return
	}

	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	var IsAnnounce bool
	isAnnounceData := gjson.GetBytes(reqMsgInJSON.Data, "isannounce")
	if isAnnounceData.Exists() {
		IsAnnounce = isAnnounceData.Bool()
	} else {
		IsAnnounce = false
	}

	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		newsgetbysenderuidReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		newsgetbysenderuidReplyStatus = "session_error"
		return
	}

	if placeString == "student" {
		newgetbyaudientuidReplyStatus = "permission_error"
		return
	}

	if placeString != "manager" {
		newsList := edumodel.GetNewsBySenderUID(int(Skip), int(Limit), IsAnnounce, reqMsgInJSON.UID)
		if newsList != nil {
			newsgetbysenderuidReplyStatus = "success"
			newsgetbysenderuidReplyData.NewsList = newsList
		} else {
			newsgetbysenderuidReplyStatus = "model_fail"
		}
	} else {
		senderUIDData := gjson.GetBytes(reqMsgInJSON.Data, "sender")
		if !senderUIDData.Exists() {
			newsgetbysenderuidReplyStatus = "senderuid_cannot_be_empty"
			return
		}

		senderUID := senderUIDData.String()
		newsList := edumodel.GetNewsBySenderUID(int(Skip), int(Limit), IsAnnounce, senderUID)
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
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsGetBySenderUIDRouter: ", newsgetbysenderuidReplyStatus)

	var jsonMsg []byte
	var err error

	if newsgetbysenderuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newsgetbysenderuidReplyStatus, newsgetbysenderuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(newsgetbysenderuidReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("NewsGetBySenderUIDRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
