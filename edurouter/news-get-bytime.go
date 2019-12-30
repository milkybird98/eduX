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
		questiongetbyclassnameReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		questiongetbyclassnameReplyStatus = "session_error"
		return
	}

	if placeString != "manager" {
		newgetbytimeorderReplyStatus = "permission_error"
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
	fmt.Println("[ROUTER] Time: ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsGetByTimeOrderRouter: ", newgetbytimeorderReplyStatus)

	var jsonMsg []byte
	var err error

	if newgetbytimeorderReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(newgetbytimeorderReplyStatus, newgetbytimeorderReplyData)
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
