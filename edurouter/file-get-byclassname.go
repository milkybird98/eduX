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

type FileGetByClassNameRouter struct {
	edunet.BaseRouter
}

type FileGetByClassNameData struct {
	ClassName string `json:"class"`
	Skip      int64  `json:"skip"`
	Limit     int64  `json:"limit"`
}

type FileGetByClassReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

var filegetbyclassnameReplyStatus string
var filegetbyclassnameReplyData FileGetByClassReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileGetByClassNameRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, filegetbyclassnameReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	filegetbyclassnameReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filegetbyclassnameReplyStatus = "data_format_error"
		return
	}

	classNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	if !classNameData.Exists() || classNameData.String() == "" {
		filegetbyclassnameReplyStatus = "classname_cannot_be_empty"
		return
	}

	className := classNameData.String()

	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		filegetbyclassnameReplyStatus = "session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		filegetbyclassnameReplyStatus = "session_error"
		return
	}

	class := edumodel.GetClassByName(className)
	if class == nil {
		filegetbyclassnameReplyStatus = "class_not_found"
		return
	}

	if placeString != "manager" {
		ok := edumodel.CheckUserInClass(className, reqMsgInJSON.UID, placeString)
		if !ok {
			filegetbyclassnameReplyStatus = "permission_error_not_in_same_class"
			return
		}
	}

	fileList := edumodel.GetFileByClassName(int(Skip), int(Limit), className)
	if fileList != nil {
		filegetbyclassnameReplyStatus = "success"
		filegetbyclassnameReplyData.FileList = fileList
	} else {
		filegetbyclassnameReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *FileGetByClassNameRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileGetByClassNameRouter: ", filegetbyclassnameReplyStatus)

	var jsonMsg []byte
	var err error

	if filegetbyclassnameReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbyclassnameReplyStatus, filegetbyclassnameReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filegetbyclassnameReplyStatus, nil)
	}
	if err != nil {
		fmt.Println("FileGetByClassNameRouter : ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
