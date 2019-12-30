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

type FileGetByTagsRouter struct {
	edunet.BaseRouter
}

type FileGetByTagsData struct {
	Tags  []string `json:"tags"`
	Skip  int64    `json:"skip"`
	Limit int64    `json:"limit"`
}

type FileGetByTagsReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

var filegetbytagsReplyStatus string
var filegetbytagsReplyData FileGetByTagsReplyData

func (router *FileGetByTagsRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, filegetbytagsReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	filegetbytagsReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filegetbytagsReplyStatus = "data_format_error"
		return
	}

	tagsData := gjson.GetBytes(reqMsgInJSON.Data, "tags")
	if !tagsData.Exists() || !tagsData.IsArray() || len(tagsData.Array()) == 0 {
		filegetbytagsReplyStatus = "tags_cannot_be_empty"
		return
	}

	var tagInString []string
	for _, tag := range tagsData.Array() {
		if tag.String() != "" {
			tagInString = append(tagInString, tag.String())
		}
	}

	if len(tagInString) == 0 {
		filegetbytagsReplyStatus = "tags_cannot_be_empty"
		return
	}

	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		filegetbytagsReplyStatus = "session_place_not_found"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		filegetbytagsReplyStatus = "session_place_data_error"
		return
	}

	var className string
	if placeString != "manager" {
		class := edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)
		if class == nil {
			filegetbytagsReplyStatus = "model_fail"
			return
		}

		className = class.ClassName
		if className == "" {
			filegetbytagsReplyStatus = "not_in_class"
			return
		}
	} else {
		className = ""
	}

	fileList := edumodel.GetFileByTags(int(Skip), int(Limit), tagInString, className)
	if fileList != nil {
		filegetbytagsReplyStatus = "success"
		filegetbytagsReplyData.FileList = fileList
	} else {
		filegetbytagsReplyStatus = "model_fail"
	}
}

func (router *FileGetByTagsRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileGetByTagsRouter: ", filegetbytagsReplyStatus)

	var jsonMsg []byte
	var err error

	if filegetbytagsReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbytagsReplyStatus, filegetbytagsReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filegetbytagsReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("FileGetByTagsRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
