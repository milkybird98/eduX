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

// PersonInfoGetByClassRouter 处理获取全班人员数据请求
type PersonInfoGetByClassRouter struct {
	edunet.BaseRouter
}

// PersonInfoGetByClassData 定义获取全班人员数据时的参数
type PersonInfoGetByClassData struct {
	ClassName string `json:"class"`
}

// PersonInfoGetByClassReplyData 定义获取全班人员数据的返回参数
type PersonInfoGetByClassReplyData struct {
	UserList []PersonInfoGetReplyData `json:"userlist"`
}

// 返回状态码
var persongetbyclassReplyStatus string

// 返回数据
var persongetbyclassReplyData PersonInfoGetByClassReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonInfoGetByClassRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	persongetbyclassReplyData.UserList = []PersonInfoGetReplyData{}
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, persongetbyclassReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	persongetbyclassReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 从Data段获取班级名称数据
	reqClassNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	// 不存在则返回错误码
	if !reqClassNameData.Exists() || reqClassNameData.String() == "" {
		persongetbyclassReplyStatus = "data_format_error"
		return
	}
	reqClassName := reqClassNameData.String()

	// 权限验证

	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		persongetbyclassReplyStatus = err.Error()
		return
	}

	// 如果当前用户是教师或者学生
	if placeString == "teacher" || placeString == "student" {
		// 检测当前用户是否在查询班级中
		ok := edumodel.CheckUserInClass(reqClassName, reqMsgInJSON.UID, placeString)
		// 如果不在则权限错误
		if !ok {
			persongetbyclassReplyStatus = "permission_error"
			return
		}
	}

	//查询数据库
	// 通过班级数据库获取人员数据
	userManyData := edumodel.GetUserByClass(reqClassName)
	if userManyData == nil || len(*userManyData) <= 0 {
		persongetbyclassReplyStatus = "data_not_found"
		return
	}

	// 将人员隐私数据进行保护
	for _, person := range *userManyData {
		var personData PersonInfoGetReplyData
		personData.UID = person.UID
		personData.Place = person.Place
		personData.Name = person.Name
		personData.ClassName = person.Class
		personData.Gender = person.Gender
		personData.Birth = person.Birth
		personData.Political = person.Political
		personData.IsPublic = person.IsContactPub
		personData.Job = person.Job
		if reqMsgInJSON.UID != person.UID {
			if personData.IsPublic {
				personData.Contact = person.Contact
				personData.Localion = person.Localion
				personData.Email = person.Email
			} else {
				personData.Contact = "未公开"
				personData.Email = "未公开"
				personData.Localion = "未公开"
			}
		} else {
			personData.Contact = person.Contact
			personData.Localion = person.Localion
			personData.Email = person.Email
		}

		personData.Com1A = person.Com1A
		personData.Com1B = person.Com1B
		personData.Com2A = person.Com2A
		personData.Com2B = person.Com2B
		personData.Com3A = person.Com3A
		personData.Com3B = person.Com3B
		personData.Com4A = person.Com4A
		personData.Com4B = person.Com4B

		// 添加到返回数据中
		persongetbyclassReplyData.UserList = append(persongetbyclassReplyData.UserList, personData)
	}

	// 设定返回状态为success
	persongetbyclassReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonInfoGetByClassRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonInfoGetByClassRouter: ", persongetbyclassReplyStatus)
	var jsonMsg []byte
	var err error

	// 生成返回数据
	if persongetbyclassReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(persongetbyclassReplyStatus, persongetbyclassReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(persongetbyclassReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println(", PersonInfoGetByClassRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
