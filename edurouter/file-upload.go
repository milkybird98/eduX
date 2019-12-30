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

type FileAddRouter struct {
	edunet.BaseRouter
}

type FileAddData struct {
	FileName  string   `json:"filename"`
	ClassName string   `json:"classname"`
	FileTag   []string `json:"filetag"`
	Size      int64    `json:"size"`
}

type FileAddReplyData struct {
	SerectID string `json:"serect"`
}

var fileaddReplyStatus string
var fileaddSerectID string

func (router *FileAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, fileaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	fileaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		fileaddReplyStatus = "data_format_error"
		return
	}

	newFileData := gjson.ParseBytes(reqMsgInJSON.Data)

	fileNameData := newFileData.Get("filename")
	if !fileNameData.Exists() {
		fileaddReplyStatus = "fileName_cannot_be_empty"
		return
	}

	classNameData := newFileData.Get("classname")
	if !classNameData.Exists() {
		fileaddReplyStatus = "className_cannot_be_empty"
		return
	}

	sizeData := newFileData.Get("size")
	if !sizeData.Exists() {
		fileaddReplyStatus = "size_cannot_be_empty"
		return
	}

	// 权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		fileaddReplyStatus = "session_error"
		return
	}

	if sessionPlace != "teacher" && sessionPlace != "student" {
		fileaddReplyStatus = "permission_error"
		return
	}

	sessionPlaceString, ok := sessionPlace.(string)
	if !ok {
		fileaddReplyStatus = "session_error"
		return
	}

	ok = edumodel.CheckUserInClass(classNameData.String(), reqMsgInJSON.UID, sessionPlaceString)
	if !ok {
		fileaddReplyStatus = "not_in_class"
		return
	}

	fileaddSerectID = primitive.NewObjectID().Hex()

	// 拼接数据
	var newFileTag utils.FileTransmitTag
	newFileTag.FileName = fileNameData.String()
	newFileTag.ID = fileaddSerectID
	newFileTag.Size = sizeData.Int()
	newFileTag.ClientAddress = c.GetTCPConnection().RemoteAddr()
	newFileTag.ClassName = classNameData.String()
	newFileTag.UpdaterUID = reqMsgInJSON.UID
	newFileTag.UpdateTime = time.Now().In(utils.GlobalObject.TimeLocal)
	newFileTag.ServerToC = false
	newFileTag.ClientToS = true

	// 更新cache
	utils.SetFileTranCacheExpire(fileaddSerectID, newFileTag)
	fileaddReplyStatus = "success"

}

// Handle 返回处理结果
func (router *FileAddRouter) Handle(request eduiface.IRequest) {
	var jsonMsg []byte
	var err error
	var fileAddReplyData FileAddReplyData
	fileAddReplyData.SerectID = fileaddSerectID

	fmt.Println("FileAddRouter: ", fileaddReplyStatus)

	if fileaddReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(fileaddReplyStatus, &fileAddReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(fileaddReplyStatus, nil)
	}

	if err != nil {
		fmt.Println("FileAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
