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

// FileGetByTimeOrderRouter 处理按照时间顺序获取消息的请求
type FileGetByTimeOrderRouter struct {
	edunet.BaseRouter
}

// FileGetByTimeOrderData 定义按照时间顺序请求消息时的参数
type FileGetByTimeOrderData struct {
	Skip  int64 `json:"skip"`
	Limit int64 `json:"limit"`
}

// FileGetByTimeOrderReplyData 定义根据时间顺序请求消息时的返回参数
type FileGetByTimeOrderReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

// 返回状态码
var filegetbytimeorderReplyStatus string

// 返回数据
var filegetbytimeorderReplyData FileGetByTimeOrderReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileGetByTimeOrderRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, filegetbytimeorderReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	filegetbytimeorderReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filegetbytimeorderReplyStatus = "data_format_error"
		return
	}

	// 获取skip和limit值
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		filegetbytimeorderReplyStatus = err.Error()
		return
	}

	// 如果当前连接用户不是管理员则权限错误
	if placeString != "manager" {
		filegetbytimeorderReplyStatus = "permission_error"
		return
	}

	// 查询数据库
	fileList := edumodel.GetFileByTimeOrder(Skip, Limit)
	// 如果数据存在则返回success和数据,否则返回错误码
	if fileList != nil {
		filegetbytimeorderReplyStatus = "success"
		filegetbytimeorderReplyData.FileList = fileList
	} else {
		filegetbytimeorderReplyStatus = "model_fail"
	}

}

// Handle 用于将请求的处理结果发回客户端
func (router *FileGetByTimeOrderRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] Time: ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileGetByTimeOrderRouter: ", filegetbytimeorderReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if filegetbytimeorderReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbytimeorderReplyStatus, filegetbytimeorderReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filegetbytimeorderReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileGetByTimeOrderRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
