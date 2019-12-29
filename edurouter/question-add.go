package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

// QuestionAddRouter 添加问题消息路由
type QuestionAddRouter struct {
	edunet.BaseRouter
}

// QuestionAddData 时间,班级等数据后端查询得到
type QuestionAddData struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

var questionaddReplyStatus string

// PreHandle 数据格式验证,检查学生是否加入班级,权限验证并更新数据库
func (router *QuestionAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, questionaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	questionaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questionaddReplyStatus = "data_format_error"
		return
	}

	newQuestionData := gjson.ParseBytes(reqMsgInJSON.Data)
	titleData := newQuestionData.Get("title")
	if !titleData.Exists() {
		questionaddReplyStatus = "title_cannot_be_empty"
		return
	}

	textData := newQuestionData.Get("text")
	if !textData.Exists() {
		questionaddReplyStatus = "title_cannot_be_empty"
		return
	}

	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		questionaddReplyStatus = "seesion_error"
		return
	}

	sessionPlaceString, ok := sessionPlace.(string)
	if !ok {
		questionaddReplyStatus = "session_error"
		return
	}

	if sessionPlaceString != "student" {
		questionaddReplyStatus = "permission_error"
		return
	}

	class := edumodel.GetClassByUID(reqMsgInJSON.UID, "student")
	if class == nil {
		questionaddReplyStatus = "not_join_class"
		return
	}

	var question edumodel.Question
	question.Title = titleData.String()
	question.Text = textData.String()
	question.SenderUID = reqMsgInJSON.UID
	question.ClassName = class.ClassName
	question.SendTime = time.Now()
	question.IsSolved = false

	ok = edumodel.AddQuestion(&question)
	if ok {
		questionaddReplyStatus = "success"
	} else {
		questionaddReplyStatus = "model_fail"
	}
}

// Handle 返回处理结果
func (router *QuestionAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionAddRouter: ", questionaddReplyStatus)
	jsonMsg, err := CombineReplyMsg(questionaddReplyStatus, nil)
	if err != nil {
		fmt.Println("QuestionAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
