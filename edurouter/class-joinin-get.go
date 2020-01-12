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

// ClassJoinInGetRouter 用于获取当前用户已加入的班级
type ClassJoinInGetRouter struct {
	edunet.BaseRouter
}

// ClassJoinInGetReplyData 定义发当前用户已加入班级时Data段的参数
type ClassJoinInGetData struct {
	UID string `json:"useruid"`
}

// ClassJoinInGetReplyData 定义返回当前用户已加入班级时Data段的参数
type ClassJoinInGetReplyData struct {
	ClassName   string    `json:"class"`
	AlterName   string    `json:"alter"`
	TeacherList []string  `json:"teachers"`
	StudentList []string  `json:"students"`
	CreateTime  time.Time `json:"time"`
}

// 返回状态
var classjoiningetReplyStatus string

// 返回数据
var classjoiningetReplyData ClassJoinInGetReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *ClassJoinInGetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, classjoiningetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	classjoiningetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	var user *edumodel.User
	uid := gjson.GetBytes(reqMsgInJSON.Data, "useruid").String()

	if uid == "" {
		uid = reqMsgInJSON.UID
	}

	user = edumodel.GetUserByUID(uid)
	if user == nil {
		classjoiningetReplyStatus = "user_not_found"
		return
	}

	// 从数据库中获取当前用户已加入班级数据
	class := edumodel.GetClassByUID(user.UID, user.Place)

	// 如果班级存在则将班级数据返回
	if class == nil {
		classjoiningetReplyStatus = "not_join_class"
	} else {
		classjoiningetReplyStatus = "success"
		classjoiningetReplyData.AlterName = class.AlterName
		classjoiningetReplyData.ClassName = class.ClassName
		classjoiningetReplyData.StudentList = class.StudentList
		classjoiningetReplyData.TeacherList = class.TeacherList
		classjoiningetReplyData.CreateTime = class.CreateDate
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *ClassJoinInGetRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ",ClassJoinInGetRouter: ", classjoiningetReplyStatus)
	var jsonMsg []byte
	var err error
	// 生成返回数据
	if classjoiningetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(classjoiningetReplyStatus, classjoiningetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(classjoiningetReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println(err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
