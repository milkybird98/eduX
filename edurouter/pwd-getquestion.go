package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"time"
)

type PwdGetQuestionRouter struct {
	edunet.BaseRouter
}

type PwdGetQuestionReplyData struct {
	UID       string `json:"uid"`
	QuestionA string `json:"qa"`
	QuestionB string `json:"qb"`
	QuestionC string `json:"qc"`
}

var pwdgetquestionReplyStatus string
var pwdgetquestionReplyData PwdGetQuestionReplyData

func (router *PwdGetQuestionRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, pwdgetquestionReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	pwdgetquestionReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	auth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	if auth == nil {
		pwdgetquestionReplyStatus = "user_not_found"
		return
	}

	pwdgetquestionReplyData.UID = reqMsgInJSON.UID
	pwdgetquestionReplyData.QuestionA = auth.QuestionA
	pwdgetquestionReplyData.QuestionB = auth.QuestionB
	pwdgetquestionReplyData.QuestionC = auth.QuestionC

	pwdgetquestionReplyStatus = "success"
}

func (router *PwdGetQuestionRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdGetQuestionRouter: ", pwdgetquestionReplyStatus)

	var jsonMsg []byte
	var err error
	if pwdgetquestionReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(pwdgetquestionReplyStatus, pwdgetquestionReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(pwdgetquestionReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("PwdGetQuestionRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
