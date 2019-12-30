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

type QuestionGetByClassNameRouter struct {
	edunet.BaseRouter
}

type QuestionGetByClassNameData struct {
	ClassName   string `json:"class"`
	Skip        int64  `json:"skip"`
	Limit       int64  `json:"limit"`
	DeferSolved bool   `json:"defer"`
	IsSolved    bool   `json:"issolved"`
}

type QuestionGetByClassReplyData struct {
	QuestionList *[]edumodel.Question `json:"questions"`
}

var questiongetbyclassnameReplyStatus string
var questiongetbyclassnameReplyData QuestionGetByClassReplyData

func (router *QuestionGetByClassNameRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, questiongetbyclassnameReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	questiongetbyclassnameReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questiongetbyclassnameReplyStatus = "data_format_error"
		return
	}

	classNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	if !classNameData.Exists() {
		questiongetbyclassnameReplyStatus = "classname_cannot_be_empty"
		return
	}

	className := classNameData.String()

	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

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
	} else {
		IsSolved = false
	}

	//权限检查
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

	class := edumodel.GetClassByName(className)
	if class == nil {
		questiongetbyclassnameReplyStatus = "class_not_found"
		return
	}

	if placeString != "manager" {
		ok := edumodel.CheckUserInClass(className, reqMsgInJSON.UID, placeString)
		if !ok {
			questiongetbyclassnameReplyStatus = "permission_error"
			return
		}
	}

	questionList := edumodel.GetQuestionByClassName(int(Skip), int(Limit), DetectSolved, IsSolved, className)
	if questionList != nil {
		questiongetbyclassnameReplyStatus = "success"
		questiongetbyclassnameReplyData.QuestionList = questionList
	} else {
		questiongetbyclassnameReplyStatus = "model_fail"
	}
}

func (router *QuestionGetByClassNameRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ",time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionGetByClassNameRouter: ", questiongetbyclassnameReplyStatus)

	var jsonMsg []byte
	var err error

	if questiongetbyclassnameReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(questiongetbyclassnameReplyStatus, questiongetbyclassnameReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(questiongetbyclassnameReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("QuestionGetByClassNameRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
