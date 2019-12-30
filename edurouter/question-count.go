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

type QuestionCountRouter struct {
	edunet.BaseRouter
}

type QuestionCountData struct {
	ClassName string    `json:"classname"`
	Date      time.Time `json:"time"`
	IsSolved  bool      `json:"issolved"`
}

type QuestionCountReplyData struct {
	Number int `json:"num"`
}

var questioncountReplyStatus string
var questioncountReplyData QuestionCountReplyData

func (router *QuestionCountRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, questioncountReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	questioncountReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questioncountReplyStatus = "data_format_error"
		return
	}

	// 权限检查
	c := request.GetConnection()

	userPlace, err := GetSessionPlace(c)
	if err != nil {
		questioncountReplyStatus = err.Error()
		return
	}

	if userPlace != "manager" {
		questioncountReplyStatus = "permission_error"
		return
	}

	// 获取参数
	timeData := gjson.GetBytes(reqMsgInJSON.Data, "time")
	targetTime, err := time.Parse(time.RFC3339, timeData.String())
	var isTimeRequired bool
	if err != nil || targetTime.IsZero() {
		isTimeRequired = false
	} else {
		isTimeRequired = true
	}

	className := gjson.GetBytes(reqMsgInJSON.Data, "classname").String()
	IsSolved := gjson.GetBytes(reqMsgInJSON.Data, "issolved").Bool()

	// 查询数据库
	if isTimeRequired {
		if IsSolved {
			questioncountReplyData.Number = edumodel.GetQuestionAnsweredNumberByDate(className, targetTime)
		} else {
			questioncountReplyData.Number = edumodel.GetQuestionNumberByDate(className, targetTime)
		}
	} else {
		if IsSolved {
			questioncountReplyData.Number = edumodel.GetQuestionAnsweredNumberByDate(className, targetTime)
		} else {
			questioncountReplyData.Number = edumodel.GetQuestionNumberByDate(className, targetTime)
		}
	}

	if questioncountReplyData.Number != -1 {
		questioncountReplyStatus = "success"
	} else {
		questioncountReplyStatus = "model_fail"
	}
}

func (router *QuestionCountRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionCountRouter: ", questioncountReplyStatus)

	var jsonMsg []byte
	var err error

	if questioncountReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(questioncountReplyStatus, questioncountReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(questioncountReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("QuestionCountRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
