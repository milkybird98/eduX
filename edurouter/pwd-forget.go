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

// PwdForgetRouter 处理忘记密码请求
type PwdForgetRouter struct {
	edunet.BaseRouter
}

// PwdForgetData 定义忘记密码请求的参数
type PwdForgetData struct {
	AnswerA string `json:"aa"`
	AnswerB string `json:"ab"`
	AnswerC string `json:"ac"`
}

// PwdForgetReplyData 定义忘记密码请求返回的参数
type PwdForgetReplyData struct {
	UID    string `json:"uid"`
	Serect string `json:"serect"`
}

// 返回状态码
var pwdforgetReplyStatus string

// 返回数据
var pwdforgetReplyData PwdForgetReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PwdForgetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, pwdforgetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		pwdforgetReplyStatus = "data_format_error"
		return
	}
	pwdForgetData := gjson.ParseBytes(reqMsgInJSON.Data)

	// 从Data段获取问题A的答案
	AnswerAData := pwdForgetData.Get("aa")
	// 如果不存在则返回错误码
	if !AnswerAData.Exists() || AnswerAData.String() == "" {
		pwdforgetReplyStatus = "answer_cannot_be_empty"
		return
	}

	// 从Data段获取问题B的答案
	AnswerBData := pwdForgetData.Get("ab")
	// 如果不存在则返回错误码
	if !AnswerBData.Exists() || AnswerBData.String() == "" {
		pwdforgetReplyStatus = "answer_cannot_be_empty"
		return
	}

	// 从Data段获取问题C的答案
	AnswerCData := pwdForgetData.Get("ac")
	// 如果不存在则返回错误码
	if !AnswerCData.Exists() || AnswerCData.String() == "" {
		pwdforgetReplyStatus = "answer_cannot_be_empty"
		return
	}

	// 试图从授权数据库获取当前用户授权数据
	auth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	// 如果不存在则报错
	if auth == nil {
		pwdforgetReplyStatus = "user_not_found"
		return
	}

	// 如果验证问题正确
	if AnswerAData.String() == auth.AnswerA &&
		AnswerBData.String() == auth.AnswerB &&
		AnswerCData.String() == auth.AnswerC {

		// 生成新的serect
		newSerect := primitive.NewObjectID().Hex()

		// 拼接验证cache
		var newCache utils.ResetPasswordTag
		newCache.UID = reqMsgInJSON.UID

		// 将验证cache添加进入cache表中
		utils.SetResetPasswordCacheExpire(newSerect, newCache)

		// 设定返回数据
		pwdforgetReplyData.UID = reqMsgInJSON.UID
		pwdforgetReplyData.Serect = newSerect

		// 设定返回状态
		pwdforgetReplyStatus = "success"
	} else {
		// 密保问题错误,返回错误码
		pwdforgetReplyStatus = "answer_wrong"
	}

}

// Handle 用于将请求的处理结果发回客户端
func (router *PwdForgetRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdForgetRouter: ", pwdforgetReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if pwdforgetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(pwdforgetReplyStatus, pwdforgetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(pwdforgetReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PwdForgetRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
