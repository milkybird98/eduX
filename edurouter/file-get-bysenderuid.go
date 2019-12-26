package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"

	"github.com/tidwall/gjson"
)

type FileGetBySenderUIDRouter struct {
	edunet.BaseRouter
}

type FileGetBySenderUIDData struct {
	Sender string `json:"sender,omitempty"`
	Skip   int64  `json:"skip"`
	Limit  int64  `json:"limit"`
}

type FileGetBySenderUIDReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

var filegetbysenderuidReplyStatus string
var filegetbysenderuidReplyData FileGetBySenderUIDReplyData

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

	var Skip int64
	skipData := gjson.GetBytes(reqMsgInJSON.Data, "skip")
	if skipData.Exists() {
		Skip = skipData.Int()
	} else {
		Skip = 0
	}

	var Limit int64
	limitData := gjson.GetBytes(reqMsgInJSON.Data, "limit")
	if limitData.Exists() {
		Limit = limitData.Int()
	} else {
		Limit = 10
	}

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
	if !senderUIDData.Exists() {
		fileList := edumodel.GetFileBySenderUID(int(Skip), int(Limit), reqMsgInJSON.UID)
		if fileList != nil {
			filegetbysenderuidReplyStatus = "success"
			filegetbysenderuidReplyData.FileList = fileList
		} else {
			filegetbysenderuidReplyStatus = "model_fail"
		}
	} else {
		if placeString != "manager" {
			user := edumodel.GetUserByUID(senderUIDData.String())
			if user == nil {
				filegetbysenderuidReplyStatus = "user_not_found"
				return
			}
			
			ok = edumodel.CheckUserInClass(user.Class, reqMsgInJSON.UID, placeString)
			if !ok {
				filegetbysenderuidReplyStatus = "not_in_class"
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
}

func (router *FileGetBySenderUIDRouter) Handle(request eduiface.IRequest) {
	fmt.Println("FileGetBySenderUIDRouter: ", filegetbysenderuidReplyStatus)

	var jsonMsg []byte
	var err error

	if filegetbysenderuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbysenderuidReplyStatus, classlistgetReplyData)
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
