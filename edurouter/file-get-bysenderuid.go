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

type FileGetBySenderUIDRouter struct {
	edunet.BaseRouter
}

type FileGetBySenderUIDData struct {
	Sender string `json:"sender"`
	Skip   int64  `json:"skip"`
	Limit  int64  `json:"limit"`
}

type FileGetBySenderUIDReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

var filegetbysenderuidReplyStatus string
var filegetbysenderuidReplyData FileGetBySenderUIDReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileGetBySenderUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, filegetbysenderuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	filegetbysenderuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filegetbysenderuidReplyStatus = "data_format_error"
		return
	}

	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		questiongetbyclassnameReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		questiongetbyclassnameReplyStatus = "session_error"
		return
	}

	senderUIDData := gjson.GetBytes(reqMsgInJSON.Data, "sender")

	if !senderUIDData.Exists() && senderUIDData.String() != "" {
		questiongetbyclassnameReplyStatus = "sender_uid_cannot_be_empty"
		return
	}

	user := edumodel.GetUserByUID(senderUIDData.String())
	if user == nil {
		filegetbysenderuidReplyStatus = "user_not_found"
		return
	}

	if user.Class == "" {
		filegetbysenderuidReplyStatus = "not_in_class"
		return
	}

	if placeString != "manager" {
		ok = edumodel.CheckUserInClass(user.Class, reqMsgInJSON.UID, placeString)
		if !ok {
			questiongetbyclassnameReplyStatus = "permission_error_not_in_same_class"
			return
		}
	}

	senderUID := senderUIDData.String()
	fileList := edumodel.GetFileBySenderUID(int(Skip), int(Limit), senderUID)
	if fileList != nil {
		filegetbysenderuidReplyStatus = "success"
		filegetbysenderuidReplyData.FileList = fileList
	} else {
		filegetbysenderuidReplyStatus = "model_fail"
	}

}

// Handle 用于将请求的处理结果发回客户端
func (router *FileGetBySenderUIDRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileGetBySenderUIDRouter: ", filegetbysenderuidReplyStatus)

	var jsonMsg []byte
	var err error

	if filegetbysenderuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbysenderuidReplyStatus, filegetbysenderuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filegetbysenderuidReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("FileGetBySenderUIDRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
