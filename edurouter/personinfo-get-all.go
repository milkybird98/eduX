package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"fmt"
	"time"
)

// PersonInfoGetAllRouter 处理获取全班人员数据请求
type PersonInfoGetAllRouter struct {
	edunet.BaseRouter
}

// PersonInfoGetAllData 定义获取全班人员数据时的参数
type PersonInfoGetAllData struct {
}

// PersonInfoGetAllReplyData 定义获取全班人员数据的返回参数
type PersonInfoGetAllReplyData struct {
	UserList []SimpleUserInfo `json:"userlist"`
}
type SimpleUserInfo struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// 返回状态码
var persongetallReplyStatus string

// 返回数据
var persongetallReplyData PersonInfoGetAllReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonInfoGetAllRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	persongetallReplyData.UserList = []SimpleUserInfo{}
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, persongetallReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	persongetallReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 权限验证
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		persongetallReplyStatus = err.Error()
		return
	}

	// 如果当前用户是教师或者学生
	if placeString == "teacher" || placeString == "student" {
		persongetallReplyStatus = "permission_error"
		return
	}

	//查询数据库
	// 通过班级数据库获取人员数据
	userManyData := edumodel.GetUserSimpleAll()
	if userManyData == nil || len(*userManyData) <= 0 {
		persongetallReplyStatus = "data_not_found"
		return
	}

	// 将人员隐私数据进行保护
	for _, person := range *userManyData {
		var newInfo SimpleUserInfo
		newInfo.Name = person.Name
		newInfo.UID = person.UID

		// 添加到返回数据中
		persongetallReplyData.UserList = append(persongetallReplyData.UserList, newInfo)
	}

	// 设定返回状态为success
	persongetallReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonInfoGetAllRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonInfoGetAllRouter: ", persongetallReplyStatus)
	var jsonMsg []byte
	var err error

	// 生成返回数据
	if persongetallReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(persongetallReplyStatus, persongetallReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(persongetallReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println(", PersonInfoGetAllRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
