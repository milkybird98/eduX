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

// FileAddRouter 处理文件上传请求
type FileAddRouter struct {
	edunet.BaseRouter
}

// FileAddData 定义上传文件请求参数
type FileAddData struct {
	FileName  string   `json:"filename"`
	ClassName string   `json:"classname"`
	FileTag   []string `json:"filetag"`
	Size      int64    `json:"size"`
}

// FileAddReplyData 定义上传的文件返回数据的参数
type FileAddReplyData struct {
	SerectID string `json:"serect"`
}

// 返回状态码
var fileaddReplyStatus string

// 返回数据
var fileAddReplyData FileAddReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, fileaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	fileaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		fileaddReplyStatus = "data_format_error"
		return
	}

	// 解码请求数据中的Data段
	newFileData := gjson.ParseBytes(reqMsgInJSON.Data)

	// 尝试获取文件名称
	fileNameData := newFileData.Get("filename")
	// 如果文件名称不存在则返回错误码
	if !fileNameData.Exists() || fileNameData.String() == "" {
		fileaddReplyStatus = "fileName_cannot_be_empty"
		return
	}

	// 尝试获取班级名称
	classNameData := newFileData.Get("classname")
	// 若班级名称不存在则返回错误码
	if !classNameData.Exists() || classNameData.String() == "" {
		fileaddReplyStatus = "className_cannot_be_empty"
		return
	}

	// 尝试获取文件大小
	sizeData := newFileData.Get("size")
	// 如果文件大小不存在则返回错误码
	if !sizeData.Exists() || sizeData.Int() <= 0 {
		fileaddReplyStatus = "size_cannot_be_empty"
		return
	}

	// 权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 如果当前连接用户身份不为教师和学生,则权限错误
	if placeString != "teacher" && placeString != "student" {
		fileaddReplyStatus = "permission_error"
		return
	}

	// 检查当前用户是否在对应班级中
	ok = edumodel.CheckUserInClass(classNameData.String(), reqMsgInJSON.UID, placeString)
	if !ok {
		fileaddReplyStatus = "not_in_class"
		return
	}

	// 生成serect
	fileAddReplyData.SerectID = primitive.NewObjectID().Hex()

	// 拼接数据
	var newFileTag utils.FileTransmitTag
	newFileTag.FileName = fileNameData.String()
	newFileTag.ID = fileAddReplyData.SerectID
	newFileTag.Size = sizeData.Int()
	// 获取当前连接客户端地址
	newFileTag.ClientAddress = c.GetTCPConnection().RemoteAddr()
	newFileTag.ClassName = classNameData.String()
	newFileTag.UpdaterUID = reqMsgInJSON.UID
	// 生成上传时间
	newFileTag.UpdateTime = time.Now().In(utils.GlobalObject.TimeLocal)
	newFileTag.ServerToC = false
	newFileTag.ClientToS = true

	// 更新cache
	utils.SetFileTranCacheExpire(fileAddReplyData.SerectID, newFileTag)
	fileaddReplyStatus = "success"

}

// Handle 用于将请求的处理结果发回客户端
func (router *FileAddRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileAddRouter: ", fileaddReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if fileaddReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(fileaddReplyStatus, &fileAddReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(fileaddReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileAddRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
