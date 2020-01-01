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

// FileGetByClassNameRouter 处理获取当前班级文件的请求
type FileGetByClassNameRouter struct {
	edunet.BaseRouter
}

// FileGetByClassNameData 定义获取班级文件元数据请求Data段参数
type FileGetByClassNameData struct {
	ClassName string `json:"class"`
	Skip      int64  `json:"skip"`
	Limit     int64  `json:"limit"`
}

// FileGetByClassReplyData 定义请求班级文件元数据的返回数据Data段参数
type FileGetByClassReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

// 返回状态码
var filegetbyclassnameReplyStatus string

// 返回数据
var filegetbyclassnameReplyData FileGetByClassReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileGetByClassNameRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 尝试解码原始数据,并且检查校验和是否正确
	reqMsgInJSON, filegetbyclassnameReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接用户是否登陆
	filegetbyclassnameReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filegetbyclassnameReplyStatus = "data_format_error"
		return
	}

	// 尝试从Data段获取班级数据
	classNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	className := classNameData.String()

	// 尝试从Data段获取skip和limit值
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		filegetbyclassnameReplyStatus = err.Error()
		return
	}

	// 如果当前用户为管理员,则检查需要获取的班级是否存在
	if placeString == "manager" {
		class := edumodel.GetClassByName(className)
		// 如果班级不存在则返回错误码
		if class == nil {
			filegetbyclassnameReplyStatus = "class_not_found"
			return
		}
		// 如果用户为教师或学生,则检查是否已加入班级
	} else if placeString == "student" || placeString == "teacher" {
		class := edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)
		// 如果未加入班级则返回错误码
		if class == nil {
			filegetbyclassnameReplyStatus = "not_in_class"
			return
		}
		className = class.ClassName
	}

	// 查询数据库,获取对应班级文件元数据列表
	fileList := edumodel.GetFileByClassName(int(Skip), int(Limit), className)
	// 获取成功则返回success,否则返回错误码
	if fileList != nil {
		filegetbyclassnameReplyStatus = "success"
		filegetbyclassnameReplyData.FileList = fileList
	} else {
		filegetbyclassnameReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *FileGetByClassNameRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileGetByClassNameRouter: ", filegetbyclassnameReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if filegetbyclassnameReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbyclassnameReplyStatus, filegetbyclassnameReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filegetbyclassnameReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileGetByClassNameRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
