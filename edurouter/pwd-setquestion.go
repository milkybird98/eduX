package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

type PwdSetQuestionRouter struct {
	edunet.BaseRouter
}

type PwdSetQuestionData struct {
	Pwd       string `json:"pwd"`
	QuestionA string `json:"qa"`
	AnswerA   string `json:"aa"`
	QuestionB string `json:"qb"`
	AnswerB   string `json:"ab"`
	QuestionC string `json:"qc"`
	AnswerC   string `json:"ac"`
}

var pwdsetquestionReplyStatus string

func (router *PwdSetQuestionRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, pwdsetquestionReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	pwdsetquestionReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		pwdsetquestionReplyStatus = "data_format_error"
		return
	}

	newQuestionData := gjson.ParseBytes(reqMsgInJSON.Data)

	// 权限检查
	pwdData := newQuestionData.Get("pwd")
	if !pwdData.Exists() || pwdData.String() == "" {
		pwdsetquestionReplyStatus = "password_cannot_be_empty"
		return
	}

	pwdInByte := []byte(pwdData.String())
	pwdInDecode, err := PwdRemoveSalr(pwdInByte)
	if err != nil {
		pwdsetquestionReplyStatus = "pwd_format_error"
		return
	}

	auth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	if auth.Pwd != string(pwdInDecode) {
		pwdsetquestionReplyStatus = "password_wrong"
		return
	}

	// 数据更新
	ok = edumodel.UpdateUserAuthByUID(reqMsgInJSON.UID, "",
		newQuestionData.Get("qa").String(),
		newQuestionData.Get("aa").String(),
		newQuestionData.Get("qb").String(),
		newQuestionData.Get("ab").String(),
		newQuestionData.Get("qc").String(),
		newQuestionData.Get("ac").String())
	if ok {
		pwdsetquestionReplyStatus = "success"
		return
	} else {
		pwdsetquestionReplyStatus = "model_fail"
		return
	}
}

func (router *PwdSetQuestionRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdSetQuestionRouter: ", pwdsetquestionReplyStatus)
	jsonMsg, err := CombineReplyMsg(pwdsetquestionReplyStatus, nil)
	if err != nil {
		fmt.Println("PwdSetQuestionRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
