package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

// FileDeleteRouter 负责文件删除业务的路由
type FileDeleteRouter struct {
	edunet.BaseRouter
}

// FileDeleteData 用于客户端指定需要删除的文件
type FileDeleteData struct {
	FileID string `json:"id"`
}

// 返回数据
var filedeleteReplyData string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileDeleteRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, filedeleteReplyData, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	filedeleteReplyData, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filedeleteReplyData = "data_format_error"
		return
	}

	// 尝试从数据Data段获取文件id
	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	// 若不存在则返回错误码
	if !innerIDData.Exists() || innerIDData.String() == "" {
		filedeleteReplyData = "fileid_cannot_be_empty"
		return
	}
	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()
	// 尝试从session中获取当前连接人员UID
	sessionUID, err := c.GetSession("UID")
	// 若不存在则返回错误码
	if err != nil {
		filedeleteReplyData = "session_error"
		return
	}

	// 试图将将uid的类型转换为字符串
	UID, ok := sessionUID.(string)
	if !ok {
		filedeleteReplyData = "session_error"
		return
	}

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 尝试根据文件id检查文件元数据是否在数据库中存在
	fileData := edumodel.GetFileByUUID(innerIDString)
	if fileData == nil {
		filedeleteReplyData = "file_not_found"
		return
	}

	// 如果不是管理员并且不是文件的上传者,则返回错误码
	if placeString != "manager" && UID != fileData.UpdaterUID {
		filedeleteReplyData = "permission_error"
		return
	}

	//数据库更新
	// 拼接文件地址
	path := "./file/" + innerIDString
	// 试图删除文件
	err = os.Remove(path)
	// 如果成功则更新数据库数据,删除对应文件元数据
	if err != nil {
		ok := edumodel.DeleteFileByUUID(innerIDString)
		// 如果成功则返回success,否则返回错误码
		if !ok {
			filedeleteReplyData = "model_fail"
		} else {
			filedeleteReplyData = "success"
		}
	} else {
		// 删除文件失败,返回错误码
		filedeleteReplyData = "remove_file_error"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *FileDeleteRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileDeleteRouter: ", filedeleteReplyData)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(filedeleteReplyData, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileDeleteRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
