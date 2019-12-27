package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

// NewsDeleteRouter 负责问题删除业务的路由
type NewsDeleteRouter struct {
	edunet.BaseRouter
}

// NewsDeleteData 用于客户端指定需要删除的问题
type NewsDeleteData struct {
	NewsID string `json:"id"`
}

var newdeleteReplyStatus string

// PreHandle 负责进行数据验证,权限验证和问题删除数据库操作
func (router *NewsDeleteRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, newdeleteReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	newdeleteReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newdeleteReplyStatus = "data_format_error"
		return
	}

	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	if !innerIDData.Exists() {
		newdeleteReplyStatus = "newsid_cannot_be_empty"
		return
	}

	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		newdeleteReplyStatus = "session_error"
		return
	}

	UID, ok := sessionUID.(string)
	if !ok {
		newdeleteReplyStatus = "session_error"
		return
	}

	sessionPlace, err := c.GetSession("place")
	if err != nil {
		newdeleteReplyStatus = "session_error"
		return
	}

	newsData := edumodel.GetNewsByInnerID(innerIDString)
	if newsData == nil {
		newdeleteReplyStatus = "news_not_found"
		return
	}

	if sessionPlace != "manager" && UID != newsData.SenderUID {
		newdeleteReplyStatus = "permission_error"
		return
	}

	//
	ok = edumodel.DeleteNewsByInnerID(innerIDString)
	if ok {
		newdeleteReplyStatus = "success"
	} else {
		newdeleteReplyStatus = "model_fail"
	}
}

// Handle 负责将处理结果发回客户端
func (router *NewsDeleteRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ",time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsDeleteRouter: ", newdeleteReplyStatus)
	jsonMsg, err := CombineReplyMsg(newdeleteReplyStatus, nil)
	if err != nil {
		fmt.Println("NewsDeleteRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
