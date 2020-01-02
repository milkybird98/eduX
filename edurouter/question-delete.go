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

// QuestionDeleteRouter 负责问题删除业务的路由
type QuestionDeleteRouter struct {
	edunet.BaseRouter
}

// QuestionDeleteData 用于客户端指定需要删除的问题
type QuestionDeleteData struct {
	QuestionID string `json:"id"`
}

var questiondeleteReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *QuestionDeleteRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, questiondeleteReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	questiondeleteReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questiondeleteReplyStatus = "data_format_error"
		return
	}

	// 试图获取问题id
	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	// 如果问题id不存在则返回错误码
	if !innerIDData.Exists() || innerIDData.String() == "" {
		questiondeleteReplyStatus = "questionid_cannot_be_empty"
		return
	}
	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		questiondeleteReplyStatus = err.Error()
		return
	}

	// 根据请求的问题id查询问题
	questionData := edumodel.GetQuestionByInnerID(innerIDString)
	// 如果问题数据不存在则返回错误码
	if questionData == nil {
		questiondeleteReplyStatus = "question_not_found"
		return
	}

	// 如果不是管理员或者问题的提出者,则权限错误
	if placeString != "manager" && reqMsgInJSON.UID != questionData.SenderUID {
		questiondeleteReplyStatus = "permission_error"
		return
	}

	// 在数据库中将问题设定为删除,以便于数据库清查时使用
	ok = edumodel.DeleteQuestionByInnerID(innerIDString)
	if ok {
		questiondeleteReplyStatus = "success"
	} else {
		questiondeleteReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *QuestionDeleteRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionDeleteRouter: ", questiondeleteReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(questiondeleteReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("QuestionDeleteRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
