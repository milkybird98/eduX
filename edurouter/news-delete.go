package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

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

var newsdeleteReplyData string

// PreHandle 负责进行数据验证,权限验证和问题删除数据库操作
func (router *NewsDeleteRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, newsdeleteReplyData, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	newsdeleteReplyData, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newsdeleteReplyData = "data_format_error"
		return
	}

	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	if !innerIDData.Exists() {
		newsdeleteReplyData = "newsid_cannot_be_empty"
		return
	}

	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		newsdeleteReplyData = "session_error"
		return
	}

	UID, ok := sessionUID.(string)
	if !ok {
		newsdeleteReplyData = "session_error"
		return
	}

	sessionPlace, err := c.GetSession("place")
	if err != nil {
		newsdeleteReplyData = "session_error"
		return
	}

	newsData := edumodel.GetNewsByInnerID(innerIDString)
	if newsData == nil {
		newsdeleteReplyData = "news_not_found"
		return
	}

	if sessionPlace != "manager" && UID != newsData.SenderUID {
		newsdeleteReplyData = "permission_error"
		return
	}

	//
	ok = edumodel.DeleteNewsByInnerID(innerIDString)
	if ok {
		newsdeleteReplyData = "success"
	} else {
		newsdeleteReplyData = "model_fail"
	}
}

// Handle 负责将处理结果发回客户端
func (router *NewsDeleteRouter) Handle(request eduiface.IRequest) {
	fmt.Println("NewsDeleteRouter: ", newsdeleteReplyData)
	jsonMsg, err := CombineReplyMsg(newsdeleteReplyData, nil)
	if err != nil {
		fmt.Println("NewsDeleteRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
