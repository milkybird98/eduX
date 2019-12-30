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

type FileCountRouter struct {
	edunet.BaseRouter
}

type FileCountData struct {
	ClassName string    `json:"classname"`
	Date      time.Time `json:"time"`
}

type FileCountReplyData struct {
	Number int `json:"num"`
}

var filecountReplyStatus string
var filecountReplyData FileCountReplyData

func (router *FileCountRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, filecountReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	filecountReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filecountReplyStatus = "data_format_error"
		return
	}

	// 权限检查
	c := request.GetConnection()

	userPlace, err := GetSessionPlace(c)
	if err != nil {
		filecountReplyStatus = err.Error()
		return
	}

	if userPlace != "manager" {
		filecountReplyStatus = "permission_error"
		return
	}

	timeData := gjson.GetBytes(reqMsgInJSON.Data, "time")
	targetTime, err := time.Parse(time.RFC3339, timeData.String())
	var isTimeRequired bool
	if err != nil || targetTime.IsZero() {
		isTimeRequired = false
	} else {
		isTimeRequired = true
	}

	className := gjson.GetBytes(reqMsgInJSON.Data, "classname").String()

	if isTimeRequired {
		filecountReplyData.Number = edumodel.GetFileNumberByDate(className, targetTime)
	} else {
		filecountReplyData.Number = edumodel.GetFileNumberAll(className)
	}

	if filecountReplyData.Number != -1 {
		filecountReplyStatus = "success"
	} else {
		filecountReplyStatus = "model_fail"
	}
}

func (router *FileCountRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileCountRouter: ", filecountReplyStatus)

	var jsonMsg []byte
	var err error

	if filecountReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filecountReplyStatus, filecountReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filecountReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("FileCountRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
