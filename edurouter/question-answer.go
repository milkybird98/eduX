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

type QuestionAnswerRouter struct {
	edunet.BaseRouter
}

type QuestionAnswerData struct {
	QuestionID string `json:"id"`
	Answer     string `json:"answer"`
}

var questionanswerReplyStatus string

func (router *QuestionAnswerRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, questionanswerReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		fmt.Println("QuestionAnswerRouter: ", questionanswerReplyStatus)
		return
	}

	questionanswerReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questionanswerReplyStatus = "data_format_error"
		return
	}

	answerData := gjson.GetBytes(reqMsgInJSON.Data, "answer")
	if !answerData.Exists() {
		questionanswerReplyStatus = "answer_cannot_be_empty"
		return
	}

	answer := answerData.String()

	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	if !innerIDData.Exists() {
		questionanswerReplyStatus = "questionid_cannot_be_empty"
		return
	}

	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		questionanswerReplyStatus = "session_error"
		return
	}

	UID, ok := sessionUID.(string)
	if !ok {
		questionanswerReplyStatus = "session_error"
		return
	}

	sessionPlace, err := c.GetSession("place")
	if err != nil {
		questionanswerReplyStatus = "session_error"
		return
	}

	if sessionPlace != "teacher" {
		questionanswerReplyStatus = "permission_error"
		return
	}

	questionData := edumodel.GetQuestionByInnerID(innerIDString)
	if questionData == nil {
		questionanswerReplyStatus = "question_not_found"
		return
	}

	ok = edumodel.CheckUserInClass(questionData.ClassName, UID, "teacher")
	if !ok {
		questionanswerReplyStatus = "permission_error"
		return
	}

	//
	ok = edumodel.AnserQuestionByInnerID(innerIDString, UID, answer)
	if ok {
		questionanswerReplyStatus = "success"
	} else {
		questionanswerReplyStatus = "model_fail"
	}
}

func (router *QuestionAnswerRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionAnswerRouter: ", questionanswerReplyStatus)
	jsonMsg, err := CombineReplyMsg(questionanswerReplyStatus, nil)
	if err != nil {
		fmt.Println("QuestionAnswerRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
