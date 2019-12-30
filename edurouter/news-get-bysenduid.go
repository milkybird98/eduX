package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
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

	var Skip int64
	skipData := gjson.GetBytes(reqMsgInJSON.Data, "skip")
	if skipData.Exists() && skipData.Int() >= 0 {
		Skip = skipData.Int()
	} else {
		Skip = 0
	}

	var Limit int64
	limitData := gjson.GetBytes(reqMsgInJSON.Data, "limit")
	if limitData.Exists() && limitData.Int() > 0 {
		Limit = limitData.Int()
	} else {
		Limit = 10
	}

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

func (router *NewsGetBySenderUIDRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsGetBySenderUIDRouter: ", newsgetbysenderuidReplyStatus)

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
