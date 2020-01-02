package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"time"
)

// PwdGetQuestionRouter 处理获取密保问题的请求
type PwdGetQuestionRouter struct {
	edunet.BaseRouter
}

// PwdGetQuestionReplyData 定义获取密保问题请求的参数
type PwdGetQuestionReplyData struct {
	UID       string `json:"uid"`
	QuestionA string `json:"qa"`
	QuestionB string `json:"qb"`
	QuestionC string `json:"qc"`
}

// 返回状态码
var pwdgetquestionReplyStatus string

// 返回数据
var pwdgetquestionReplyData PwdGetQuestionReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PwdGetQuestionRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, pwdgetquestionReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 获取验证数据
	auth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	// 如果未找到验证数据则报错
	if auth == nil {
		pwdgetquestionReplyStatus = "user_not_found"
		return
	}

	// 拼接数据
	pwdgetquestionReplyData.UID = reqMsgInJSON.UID
	pwdgetquestionReplyData.QuestionA = auth.QuestionA
	pwdgetquestionReplyData.QuestionB = auth.QuestionB
	pwdgetquestionReplyData.QuestionC = auth.QuestionC

	// 设定状态码
	pwdgetquestionReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *PwdGetQuestionRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdGetQuestionRouter: ", pwdgetquestionReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if pwdgetquestionReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(pwdgetquestionReplyStatus, pwdgetquestionReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(pwdgetquestionReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PwdGetQuestionRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
