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

// ClassDelRouter 课程删除路由
type ClassDelRouter struct {
	edunet.BaseRouter
}

// ClassDelData 课程删除接口数据结构
type ClassDelData struct {
	ClassName string `json:"class"`
}

var classdelReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *ClassDelRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, classdelReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	classdelReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		classdelReplyStatus = "data_format_error"
		return
	}

	delClassData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	if !delClassData.Exists() {
		classdelReplyStatus = "data_format_error"
		return
	}

	delClassName := delClassData.String()

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		classdelReplyStatus = "session_error"
		return
	}

	if sessionPlace != "manager" {
		classdelReplyStatus = "permission_error"
		return
	}

	//数据库操作
	class := edumodel.GetClassByName(delClassName)
	if class == nil {
		classdelReplyStatus = "class_not_found"
		return
	}

	ok = edumodel.DeleteClassByName(delClassName)
	if ok == true {
		classdelReplyStatus = "success"
	} else {
		classdelReplyStatus = "model_fail"
	}
}

//Handle 返回课程删除结果
// Handle 用于将请求的处理结果发回客户端
func (router *ClassDelRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", ClassDelRouter: ", classdelReplyStatus)
	jsonMsg, err := CombineReplyMsg(classdelReplyStatus, nil)
	if err != nil {
		fmt.Println("ClassDelRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
