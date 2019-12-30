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

type NewsGetByAudientUIDRouter struct {
	edunet.BaseRouter
}

type NewsGetByAudientUIDData struct {
	Audient    string `json:"audient,omitempty"`
	Skip       int64  `json:"skip"`
	Limit      int64  `json:"limit"`
	IsAnnounce bool   `json:"isannounce"`
}

type NewGetByAudientUIDReplyData struct {
	NewsList *[]edumodel.News `json:"news"`
}

var newgetbyaudientuidReplyStatus string
var newgetbyaudientuidReplyData NewGetByAudientUIDReplyData

func (router *NewsGetByAudientUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, newgetbyaudientuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	newgetbyaudientuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newgetbyaudientuidReplyStatus = "data_format_error"
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
		newgetbyaudientuidReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		newgetbyaudientuidReplyStatus = "session_error"
		return
	}

	if placeString != "manager" {
		newsList := edumodel.GetNewsByAudientUID(int(Skip), int(Limit), IsAnnounce, reqMsgInJSON.UID)
		if newsList != nil {
			newgetbyaudientuidReplyStatus = "success"
			newgetbyaudientuidReplyData.NewsList = newsList
		} else {
			newgetbyaudientuidReplyStatus = "model_fail"
		}
	} else {
		audientUIDData := gjson.GetBytes(reqMsgInJSON.Data, "audient")
		if !audientUIDData.Exists() {
			newgetbyaudientuidReplyStatus = "audientuid_cannot_be_empty"
			return
		}

		audientUID := audientUIDData.String()
		newsList := edumodel.GetNewsByAudientUID(int(Skip), int(Limit), IsAnnounce, audientUID)
		if newsList != nil {
			newgetbyaudientuidReplyStatus = "success"
			newgetbyaudientuidReplyData.NewsList = newsList
		} else {
			newgetbyaudientuidReplyStatus = "model_fail"
		}
	}
}

func (router *NewsGetByAudientUIDRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ",time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsGetByAudientUIDRouter: ", newgetbyaudientuidReplyStatus)

	var jsonMsg []byte
	var err error

	if newgetbyaudientuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newgetbyaudientuidReplyStatus, newgetbyaudientuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(newgetbyaudientuidReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("NewsGetByAudientUIDRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
