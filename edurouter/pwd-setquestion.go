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

// PwdSetQuestionRouter 处理重设密保问题的请求
type PwdSetQuestionRouter struct {
	edunet.BaseRouter
}

// PwdSetQuestionData 定义重设密保问题请求的参数
type PwdSetQuestionData struct {
	Pwd       string `json:"pwd"`
	QuestionA string `json:"qa"`
	AnswerA   string `json:"aa"`
	QuestionB string `json:"qb"`
	AnswerB   string `json:"ab"`
	QuestionC string `json:"qc"`
	AnswerC   string `json:"ac"`
}

// 返回状态码
var pwdsetquestionReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PwdSetQuestionRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, pwdsetquestionReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	pwdsetquestionReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		pwdsetquestionReplyStatus = "data_format_error"
		return
	}

	newQuestionData := gjson.ParseBytes(reqMsgInJSON.Data)
	// 从Data段获取密码
	pwdData := newQuestionData.Get("pwd")
	// 如果密码不存在则报错返回
	if !pwdData.Exists() || pwdData.String() == "" {
		pwdsetquestionReplyStatus = "password_cannot_be_empty"
		return
	}

	// 密码去盐
	pwdInByte := []byte(pwdData.String())
	// 如果去盐失败则报错返回
	pwdInDecode, err := PwdRemoveSalr(pwdInByte)
	if err != nil {
		pwdsetquestionReplyStatus = "pwd_format_error"
		return
	}

	// 权限检查

	// 获取授权数据
	auth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	// 检查用户授权数据是否存在
	if auth == nil {
		pwdsetquestionReplyStatus = "user_not_found"
		return
	}
	// 判断密码是否正确
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
	if ok { // 如果更新成功则返回success,否则返回错误码
		pwdsetquestionReplyStatus = "success"
		return
	} else {
		pwdsetquestionReplyStatus = "model_fail"
		return
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *PwdSetQuestionRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdSetQuestionRouter: ", pwdsetquestionReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(pwdsetquestionReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PwdSetQuestionRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
