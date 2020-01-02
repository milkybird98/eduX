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

// FileCountRouter 用于处理管理员统计文件数量的请求
type FileCountRouter struct {
	edunet.BaseRouter
}

// FileCountData 定义了管理员请求文件数量时Data段的参数
type FileCountData struct {
	ClassName string    `json:"classname"`
	Date      time.Time `json:"time"`
}

// FileCountReplyData 定义了管理员统计文件数量请求的返回数据Data段参数
type FileCountReplyData struct {
	Number int `json:"num"`
}

// 返回状态码
var filecountReplyStatus string

// 返回数据
var filecountReplyData FileCountReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileCountRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, filecountReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	filecountReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filecountReplyStatus = "data_format_error"
		return
	}

	// 权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	userPlace, err := GetSessionPlace(c)
	if err != nil {
		filecountReplyStatus = err.Error()
		return
	}

	// 如果身份不为管理员则权限错误
	if userPlace != "manager" {
		filecountReplyStatus = "permission_error"
		return
	}

	// 试图从Data段中获取日期数据
	timeData := gjson.GetBytes(reqMsgInJSON.Data, "time")
	// 按照RFC3339标准解码时间数据
	targetTime, err := time.Parse(time.RFC3339, timeData.String())
	// 如果日期数据存在则指定统计日期,否则统计全部数据
	var isTimeRequired bool
	if err != nil || targetTime.IsZero() {
		isTimeRequired = false
	} else {
		isTimeRequired = true
	}

	// 尝试从Data段中获取班级名称
	className := gjson.GetBytes(reqMsgInJSON.Data, "classname").String()

	// 根据是否存在日期限定,进行数据库查询
	if isTimeRequired {
		filecountReplyData.Number = edumodel.GetFileNumberByDate(className, targetTime)
	} else {
		filecountReplyData.Number = edumodel.GetFileNumberAll(className)
	}

	// 如果查询成功在返回success,否则返回错误码
	if filecountReplyData.Number != -1 {
		filecountReplyStatus = "success"
	} else {
		filecountReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *FileCountRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileCountRouter: ", filecountReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if filecountReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filecountReplyStatus, filecountReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filecountReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileCountRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
