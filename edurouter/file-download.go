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

// FileDownloadRouter 用于处理文件下载请求
type FileDownloadRouter struct {
	edunet.BaseRouter
}

// FileDownloadData 定义文件下载请求中Data段参数
type FileDownloadData struct {
	UUID string `json:"uuid"`
}

// FileDownloadReplyData 定义文件下载请求的回复数据Data段参数
type FileDownloadReplyData struct {
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
	Serect   string `json:"serect"`
}

// 返回错误码
var filedownloadReplyStatus string

// 返回数据
var filedownloadReplyData FileDownloadReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileDownloadRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 尝试对序列化数据进行反序列话,并检查校验和是否正确
	reqMsgInJSON, filedownloadReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接对应用户是否已登录
	filedownloadReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 检查数据Data段是否符合json编码规范
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filedownloadReplyStatus = "data_format_error"
		return
	}

	// 对Data段进行反序列化
	newFileData := gjson.ParseBytes(reqMsgInJSON.Data)

	// 尝试获取uuid数据
	UUIDData := newFileData.Get("uuid")
	// 若文件对应uuid数据不存在,则返回错误码
	if !UUIDData.Exists() {
		filedownloadReplyStatus = "UUID_cannot_be_empty"
		return
	}
	UUID := UUIDData.String()

	//权限验证
	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	//数据验证

	// 检查uuid对应文件是否存在
	file := edumodel.GetFileByUUID(UUID)
	// 若不存在则返回错误码
	if file == nil {
		classdelReplyStatus = "file_not_found"
		return
	}

	// 如果当前用户不是管理员且不是文件归属班级的成员,则权限错误
	if placeString != "manager" && !edumodel.CheckUserInClass(file.ClassName, reqMsgInJSON.UID, placeString) {
		filedownloadReplyStatus = "permission_error"
		return
	}

	// 拼接数据
	var newFileTag utils.FileTransmitTag
	newFileTag.FileName = file.FileName
	newFileTag.ID = UUID
	newFileTag.Size = int64(file.Size)
	// 获取当前连接客户端地址
	newFileTag.ClientAddress = c.GetTCPConnection().RemoteAddr()
	newFileTag.ServerToC = true
	newFileTag.ClientToS = false

	// 生成uuid
	serect := primitive.NewObjectID().Hex()
	// 更新cache
	utils.SetFileTranCacheExpire(serect, newFileTag)
	filedownloadReplyData.FileName = file.FileName
	filedownloadReplyData.Size = int64(file.Size)
	filedownloadReplyData.Serect = serect

	// 返回处理状态
	filedownloadReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *FileDownloadRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileDownloadRouter: ", filedownloadReplyStatus)
	var jsonMsg []byte
	var err error
	// 生成返回数据
	if filedownloadReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filedownloadReplyStatus, &filedownloadReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filedownloadReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileDownloadRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
