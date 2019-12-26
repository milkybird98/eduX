package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

type NewsGetByTimeOrderRouter struct {
	edunet.BaseRouter
}

type NewsGetByTimeOrderData struct {
	Skip       int64 `json:"skip"`
	Limit      int64 `json:"limit"`
	IsAnnounce bool  `json:"isannounce"`
}

type NewGetByTimeOrderReplyData struct {
	NewsList *[]edumodel.News `json:"news"`
}

var newgetbytimeorderReplyStatus string
var newgetbytimeorderReplyData NewGetByTimeOrderReplyData

func (router *NewsGetByTimeOrderRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, newgetbytimeorderReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	newgetbytimeorderReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newgetbytimeorderReplyStatus = "data_format_error"
		return
	}

	var Skip int64
	skipData := gjson.GetBytes(reqMsgInJSON.Data, "skip")
	if skipData.Exists() {
		Skip = skipData.Int()
	} else {
		Skip = 0
	}

	var Limit int64
	limitData := gjson.GetBytes(reqMsgInJSON.Data, "limit")
	if limitData.Exists() {
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
		questiongetbyclassnameReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		questiongetbyclassnameReplyStatus = "session_error"
		return
	}

	if placeString != "manager" {
		newgetbyaudientuidReplyStatus = "permission_error"
		return
	}

	newsList := edumodel.GetNewsByTimeOrder(int(Skip), int(Limit), IsAnnounce)
	if newsList != nil {
		newgetbytimeorderReplyStatus = "success"
		newgetbytimeorderReplyData.NewsList = newsList
	} else {
		newgetbytimeorderReplyStatus = "model_fail"
	}

}

func (router *NewsGetByTimeOrderRouter) Handle(request eduiface.IRequest) {
	fmt.Println("NewsGetByTimeOrderRouter: ", newgetbytimeorderReplyStatus)

	var jsonMsg []byte
	var err error

	if newgetbytimeorderReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newgetbytimeorderReplyStatus, classlistgetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(newgetbytimeorderReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("NewsGetByTimeOrderRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
