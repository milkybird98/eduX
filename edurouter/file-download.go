package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tidwall/gjson"
)

type FileDownloadRouter struct {
	edunet.BaseRouter
}

type FileDownloadData struct {
	UUID string `json:"uuid"`
}

type FileDownloadReplyData struct {
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
	Serect   string `json:"serect"`
}

var filedownloadReplyStatus string
var filedownloadReplyData FileDownloadReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileDownloadRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, filedownloadReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	filedownloadReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filedownloadReplyStatus = "data_format_error"
		return
	}

	newFileData := gjson.ParseBytes(reqMsgInJSON.Data)

	UUIDData := newFileData.Get("serect")
	if !UUIDData.Exists() {
		filedownloadReplyStatus = "UUID_cannot_be_empty"
		return
	}

	// 权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		filedownloadReplyStatus = "session_error"
		return
	}

	UUID := UUIDData.String()

	file := edumodel.GetFileByUUID(UUID)

	sessionPlaceString, ok := sessionPlace.(string)
	if !ok {
		filedownloadReplyStatus = "session_error"
		return
	}

	if sessionPlace != "manager" && !edumodel.CheckUserInClass(file.ClassName, reqMsgInJSON.UID, sessionPlaceString) {
		filedownloadReplyStatus = "permission_error"
		return
	}

	// 拼接数据
	var newFileTag utils.FileTransmitTag
	newFileTag.FileName = file.FileName
	newFileTag.ID = UUID
	newFileTag.Size = int64(file.Size)
	newFileTag.ClientAddress = c.GetTCPConnection().RemoteAddr()
	newFileTag.ServerToC = true
	newFileTag.ClientToS = false

	// 更新cache
	serect := primitive.NewObjectID().Hex()

	utils.SetFileTranCacheExpire(serect, newFileTag)
	filedownloadReplyData.FileName = file.FileName
	filedownloadReplyData.Size = int64(file.Size)
	filedownloadReplyData.Serect = serect

	filedownloadReplyStatus = "success"
}

// Handle 返回处理结果
// Handle 用于将请求的处理结果发回客户端
func (router *FileDownloadRouter) Handle(request eduiface.IRequest) {
	var jsonMsg []byte
	var err error

	fmt.Println("FileDownloadRouter: ", filedownloadReplyStatus)

	if filedownloadReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filedownloadReplyStatus, &filedownloadReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filedownloadReplyStatus, nil)
	}

	if err != nil {
		fmt.Println("FileDownloadRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
