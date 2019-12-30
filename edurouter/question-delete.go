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

// PreHandle 负责进行数据验证,权限验证和问题删除数据库操作
func (router *QuestionDeleteRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, questiondeleteReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	questiondeleteReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questiondeleteReplyStatus = "data_format_error"
		return
	}

	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	if !innerIDData.Exists() {
		questiondeleteReplyStatus = "questionid_cannot_be_empty"
		return
	}

	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		questiondeleteReplyStatus = "session_error"
		return
	}

	UID, ok := sessionUID.(string)
	if !ok {
		questiondeleteReplyStatus = "session_error"
		return
	}

	sessionPlace, err := c.GetSession("place")
	if err != nil {
		questiondeleteReplyStatus = "session_error"
		return
	}

	questionData := edumodel.GetQuestionByInnerID(innerIDString)
	if questionData == nil {
		questiondeleteReplyStatus = "question_not_found"
		return
	}

	if sessionPlace != "manager" && UID != questionData.SenderUID {
		questiondeleteReplyStatus = "permission_error"
		return
	}

	//
	ok = edumodel.DeleteQuestionByInnerID(innerIDString)
	if ok {
		questiondeleteReplyStatus = "success"
	} else {
		questiondeleteReplyStatus = "model_fail"
	}
}

// Handle 负责将处理结果发回客户端
func (router *QuestionDeleteRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ",time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionDeleteRouter: ", questiondeleteReplyStatus)
	jsonMsg, err := CombineReplyMsg(questiondeleteReplyStatus, nil)
	if err != nil {
		fmt.Println("QuestionDeleteRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
