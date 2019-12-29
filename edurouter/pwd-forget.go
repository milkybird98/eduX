package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tidwall/gjson"
)

type PwdForgetRouter struct {
	edunet.BaseRouter
}

type PwdForgetData struct {
	Pwd     string `json:"pwd"`
	AnswerA string `json:"aa"`
	AnswerB string `json:"ab"`
	AnswerC string `json:"ac"`
}

type PwdForgetReplyData struct {
	UID    string `json:"uid"`
	serect string `json:"serect"`
}

var pwdforgetReplyStatus string
var pwdforgetReplyData PwdForgetReplyData

func (router *PwdForgetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, pwdforgetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		pwdforgetReplyStatus = "data_format_error"
		return
	}

	pwdForgetData := gjson.ParseBytes(reqMsgInJSON.Data)
	AnswerAData := pwdForgetData.Get("aa")
	if !AnswerAData.Exists() || AnswerAData.String() == "" {
		pwdforgetReplyStatus = "answer_cannot_be_empty"
		return
	}

	AnswerBData := pwdForgetData.Get("ab")
	if !AnswerBData.Exists() || AnswerBData.String() == "" {
		pwdforgetReplyStatus = "answer_cannot_be_empty"
		return
	}

	AnswerCData := pwdForgetData.Get("ac")
	if !AnswerCData.Exists() || AnswerCData.String() == "" {
		pwdforgetReplyStatus = "answer_cannot_be_empty"
		return
	}

	auth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	if auth == nil {
		pwdforgetReplyStatus = "user_not_found"
		return
	}

	if AnswerAData.String() == auth.AnswerA &&
		AnswerBData.String() == auth.AnswerB &&
		AnswerCData.String() == auth.AnswerC {
		newSerect := primitive.NewObjectID().Hex()
		newCache := utils.ResetPasswordTag{reqMsgInJSON.UID}

		utils.SetResetPasswordCacheExpire(newSerect, newCache)

		pwdforgetReplyData.UID = reqMsgInJSON.UID
		pwdforgetReplyData.serect = newSerect

		pwdforgetReplyStatus = "success"
	} else {
		pwdforgetReplyStatus = "answer_wrong"
	}

}

func (router *PwdForgetRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdForgetRouter: ", pwdforgetReplyStatus)

	var jsonMsg []byte
	var err error
	if pwdforgetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(pwdforgetReplyStatus, pwdforgetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(pwdforgetReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("PwdForgetRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
