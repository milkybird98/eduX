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

// FileGetBySenderUIDRouter 处理获取某一用户上传文件元数据的列表请求
type FileGetBySenderUIDRouter struct {
	edunet.BaseRouter
}

// FileGetBySenderUIDData 定义根据上传者uid获取文件元数据时的参数列表
type FileGetBySenderUIDData struct {
	Sender string `json:"sender"`
	Skip   int64  `json:"skip"`
	Limit  int64  `json:"limit"`
}

// FileGetBySenderUIDReplyData 定了请求特定发送者文件元数据时的返回参数
type FileGetBySenderUIDReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

// 返回状态码
var filegetbysenderuidReplyStatus string

// 返回数据
var filegetbysenderuidReplyData FileGetBySenderUIDReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileGetBySenderUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 尝试解码收到的数据,并验证校验和是否正确
	reqMsgInJSON, filegetbysenderuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接对应用户是否登陆
	filegetbysenderuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据数据Data格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filegetbysenderuidReplyStatus = "data_format_error"
		return
	}

	// 试图从Data中获取skip和limit数据
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	// 权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 尝试获取发送者uid数据
	senderUIDData := gjson.GetBytes(reqMsgInJSON.Data, "sender")
	// 如果不存在,则返回错误码
	if !senderUIDData.Exists() && senderUIDData.String() != "" {
		questiongetbyclassnameReplyStatus = "sender_uid_cannot_be_empty"
		return
	}
	senderUID := senderUIDData.String()

	// 检查请求uid对应用户是否存在
	user := edumodel.GetUserByUID(senderUID)
	// 若不存在则返回错误码
	if user == nil {
		filegetbysenderuidReplyStatus = "user_not_found"
		return
	}

	// 若查询用户未未加入班级则返回错误码
	if user.Class == "" {
		filegetbysenderuidReplyStatus = "not_in_class"
		return
	}

	// 如果当前用户不是管理员,则检查当前用户是否和请求用户在同一班级
	if placeString != "manager" {
		ok = edumodel.CheckUserInClass(user.Class, reqMsgInJSON.UID, placeString)
		// 若不在则出现权限错误
		if !ok {
			questiongetbyclassnameReplyStatus = "permission_error_not_in_same_class"
			return
		}
	}

	// 获取文件列表
	fileList := edumodel.GetFileBySenderUID(int(Skip), int(Limit), senderUID)
	// 如果成功则返回success,否则返回错误码
	if fileList != nil {
		filegetbysenderuidReplyStatus = "success"
		filegetbysenderuidReplyData.FileList = fileList
	} else {
		filegetbysenderuidReplyStatus = "model_fail"
	}

}

// Handle 用于将请求的处理结果发回客户端
func (router *FileGetBySenderUIDRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileGetBySenderUIDRouter: ", filegetbysenderuidReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if filegetbysenderuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbysenderuidReplyStatus, filegetbysenderuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filegetbysenderuidReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileGetBySenderUIDRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
