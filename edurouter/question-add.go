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

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *QuestionAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, questionaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	questionaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questionaddReplyStatus = "data_format_error"
		return
	}

	newQuestionData := gjson.ParseBytes(reqMsgInJSON.Data)
	// 从Data段获取问题标题
	titleData := newQuestionData.Get("title")
	// 检查问题标题是否存在
	if !titleData.Exists() || titleData.String() == "" {
		questionaddReplyStatus = "title_cannot_be_empty"
		return
	}

	// 从Data段获取问题正文
	textData := newQuestionData.Get("text")
	// 检查问题正文是否存在
	if !textData.Exists() || textData.String() == "" {
		questionaddReplyStatus = "title_cannot_be_empty"
		return
	}

	//身份验证

	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 如果当前用户是学生,则权限错误
	if placeString != "student" {
		questionaddReplyStatus = "permission_error"
		return
	}

	// 试图获取班级数据 
	class := edumodel.GetClassByUID(reqMsgInJSON.UID, "student")
	// 未加入班级,报错返回
	if class == nil {
		questionaddReplyStatus = "not_join_class"
		return
	}

	// 拼接新的问题数据
	var question edumodel.Question
	question.Title = titleData.String()
	question.Text = textData.String()
	question.SenderUID = reqMsgInJSON.UID
	question.ClassName = class.ClassName
	// 获取当前时间
	question.SendTime = time.Now().In(utils.GlobalObject.TimeLocal)
	question.IsSolved = false

	// 更新数据库
	ok = edumodel.AddQuestion(&question)
	if ok {
		questionaddReplyStatus = "success"
	} else {
		questionaddReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *QuestionAddRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionAddRouter: ", questionaddReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(questionaddReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("QuestionAddRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
