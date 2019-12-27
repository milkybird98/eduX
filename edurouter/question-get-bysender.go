package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

type QuestionGetBySenderUIDRouter struct {
	edunet.BaseRouter
}

type QuestionGetBySenderUIDData struct {
	SenderUID   string `json:"uid"`
	Skip        int64  `json:"skip"`
	Limit       int64  `json:"limit"`
	DeferSolved bool   `json:"defer"`
	IsSolved    bool   `json:"issolved"`
}

type QuestionGetBySenderUIDReplyData struct {
	QuestionList *[]edumodel.Question `json:"questions"`
}

var questiongetbysenderuidReplyStatus string
var questiongetbysenderuidReplyData QuestionGetBySenderUIDReplyData

func (router *QuestionGetBySenderUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, questiongetbysenderuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	questiongetbysenderuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questiongetbysenderuidReplyStatus = "data_format_error"
		return
	}

	senderUIDData := gjson.GetBytes(reqMsgInJSON.Data, "uid")
	if !senderUIDData.Exists() {
		questiongetbysenderuidReplyStatus = "senderuid_cannot_be_empty"
		return
	}

	senderUID := senderUIDData.String()

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

	var DetectSolved bool
	detectSolvedData := gjson.GetBytes(reqMsgInJSON.Data, "defer")
	if detectSolvedData.Exists() {
		DetectSolved = detectSolvedData.Bool()
	} else {
		DetectSolved = false
	}

	var IsSolved bool
	issolvedData := gjson.GetBytes(reqMsgInJSON.Data, "issolved")
	if issolvedData.Exists() {
		IsSolved = issolvedData.Bool()
	}

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		questiongetbysenderuidReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		questiongetbysenderuidReplyStatus = "session_error"
		return
	}

	class := edumodel.GetClassByUID(senderUID, "student")
	if class == nil {
		questiongetbysenderuidReplyStatus = "class_not_found"
		return
	}

	if placeString != "manager" {
		ok := edumodel.CheckUserInClass(class.ClassName, reqMsgInJSON.UID, placeString)
		if !ok {
			questiongetbysenderuidReplyStatus = "permission_error"
			return
		}
	}

	questionList := edumodel.GetQuestionBySenderUID(int(Skip), int(Limit), DetectSolved, IsSolved, senderUID)
	if questionList != nil {
		questiongetbysenderuidReplyStatus = "success"
		questiongetbysenderuidReplyData.QuestionList = questionList
	} else {
		questiongetbysenderuidReplyStatus = "model_fail"
	}
}

func (router *QuestionGetBySenderUIDRouter) Handle(request eduiface.IRequest) {
	fmt.Println("QuestionGetBySenderUIDRouter: ", questiongetbysenderuidReplyStatus)

	var jsonMsg []byte
	var err error

	if questiongetbysenderuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(questiongetbysenderuidReplyStatus, questiongetbysenderuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(questiongetbysenderuidReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("QuestionGetBySenderUIDRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
