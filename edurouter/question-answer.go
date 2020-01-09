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

// QuestionAnswerRouter 处理教师回答问题请求
type QuestionAnswerRouter struct {
	edunet.BaseRouter
}

// QuestionAnswerData 定义教师回答问题请求的参数
type QuestionAnswerData struct {
	QuestionID string `json:"id"`
	Answer     string `json:"answer"`
}

// 返回状态码
var questionanswerReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *QuestionAnswerRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, questionanswerReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		fmt.Println("QuestionAnswerRouter: ", questionanswerReplyStatus)
		return
	}

	// 检查当前连接是否已登录
	questionanswerReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questionanswerReplyStatus = "data_format_error"
		return
	}

	// 试图获取答案数据
	answerData := gjson.GetBytes(reqMsgInJSON.Data, "answer")
	// 如果答案不存在则报错
	if !answerData.Exists() || answerData.String() == "" {
		questionanswerReplyStatus = "answer_cannot_be_empty"
		return
	}
	answer := answerData.String()

	// 试图获取问题id
	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	// 问题id不存在则报错,返回错误码
	if !innerIDData.Exists() || innerIDData.String() == "" {
		questionanswerReplyStatus = "questionid_cannot_be_empty"
		return
	}
	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		questionanswerReplyStatus = err.Error()
		return
	}

	// 如果请求用户不是教师则权限错误
	if placeString != "teacher" {
		questionanswerReplyStatus = "permission_error"
		return
	}

	// 获取问题数据
	questionData := edumodel.GetQuestionByInnerID(innerIDString)
	// 如果问题数据未找到则返回错误码
	if questionData == nil {
		questionanswerReplyStatus = "question_not_found"
		return
	}

	// 检查教师是否和提问者在同一班级,若不在则权限错误
	ok = edumodel.CheckUserInClass(questionData.ClassName, reqMsgInJSON.UID, "teacher")
	if !ok {
		questionanswerReplyStatus = "permission_error"
		return
	}

	// 更新答案
	ok = edumodel.AnserQuestionByInnerID(innerIDString, reqMsgInJSON.UID, answer)
	if ok {
		questionanswerReplyStatus = "success"
	} else {
		questionanswerReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *QuestionAnswerRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionAnswerRouter: ", questionanswerReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(questionanswerReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("QuestionAnswerRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
