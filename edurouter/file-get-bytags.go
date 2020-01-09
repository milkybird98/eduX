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

// FileGetByTagsRouter 处理根据文件标签获取文件的请求
type FileGetByTagsRouter struct {
	edunet.BaseRouter
}

// FileGetByTagsData 定义根据标签获取文件请求的参数
type FileGetByTagsData struct {
	Tags  []string `json:"filetag"`
	Skip  int64    `json:"skip"`
	Limit int64    `json:"limit"`
}

// FileGetByTagsReplyData 定义根据标签获取文件返回数据的参数
type FileGetByTagsReplyData struct {
	FileList *[]edumodel.File `json:"files"`
}

// 返回状态码
var filegetbytagsReplyStatus string

// 返回数据
var filegetbytagsReplyData FileGetByTagsReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *FileGetByTagsRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, filegetbytagsReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	filegetbytagsReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filegetbytagsReplyStatus = "data_format_error"
		return
	}

	// 从Data段获取文件标签数据
	tagsData := gjson.GetBytes(reqMsgInJSON.Data, "filetag")
	// 如果文件标签不存在则返回错误码
	if !tagsData.Exists() || !tagsData.IsArray() || len(tagsData.Array()) == 0 {
		filegetbytagsReplyStatus = "tags_cannot_be_empty"
		return
	}

	// 将文件标签数据转换为切片
	var tagInString []string
	for _, tag := range tagsData.Array() {
		if tag.String() != "" {
			tagInString = append(tagInString, tag.String())
		}
	}

	// 如果切片长度为0,即文件标签全为空串,在返回错误码
	if len(tagInString) == 0 {
		filegetbytagsReplyStatus = "tags_cannot_be_empty"
		return
	}

	// 获取skip和limit数据
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	//权限检查
	c := request.GetConnection()

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		filegetbytagsReplyStatus = err.Error()
		return
	}

	var className string
	// 如果连接用户身份为教师或者学生
	if placeString == "teacher" || placeString == "student" {
		// 获取当前用户已加入班级
		class := edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)
		// 若未加入班级则返回错误码
		if class == nil {
			filegetbytagsReplyStatus = "model_fail"
			return
		}
		className = class.ClassName
		// 如果用户是管理员,则返回全部班级的数据
	} else if placeString == "manager" {
		className = ""
	}

	//数据查询

	// 查询根据tags和班级名查询数据库
	fileList := edumodel.GetFileByTags(int(Skip), int(Limit), tagInString, className)
	// 如果成功则返回success,否则返回错误码
	if fileList != nil {
		filegetbytagsReplyStatus = "success"
		filegetbytagsReplyData.FileList = fileList
	} else {
		filegetbytagsReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *FileGetByTagsRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileGetByTagsRouter: ", filegetbytagsReplyStatus)

	var jsonMsg []byte
	var err error
	// 生成返回数据
	if filegetbytagsReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(filegetbytagsReplyStatus, filegetbytagsReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(filegetbytagsReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("FileGetByTagsRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
