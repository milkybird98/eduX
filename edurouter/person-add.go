package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

// PersonAddRouter 处理添加人员请求
type PersonAddRouter struct {
	edunet.BaseRouter
}

// PersonAddData 定义添加人员请求参数
type PersonAddData struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Place string `json:"place"`
}

// 返回状态码
var personAddReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	var reqDataInJSON PersonAddData
	reqMsgInJSON, personAddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	personAddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personAddReplyStatus = "data_format_error"
		return
	}

	// 获取数据
	newStudentData := gjson.ParseBytes(reqMsgInJSON.Data)
	reqDataInJSON.UID = newStudentData.Get("uid").String()
	reqDataInJSON.Name = newStudentData.Get("name").String()
	reqDataInJSON.Place = newStudentData.Get("place").String()

	//权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 如果用户身份不是管理员,则权限错误
	if placeString != "manager" {
		personAddReplyStatus = "permission_error"
		return
	}

	// 检查是否有同名用户存在
	userData := edumodel.GetUserByUID(reqDataInJSON.UID)
	// 若存在则返回错误码
	if userData != nil {
		personAddReplyStatus = "same_uid_exist"
		return
	}

	//数据库操作
	var newUser edumodel.User
	var newUserAuth edumodel.UserAuth

	newUser.UID = reqDataInJSON.UID
	newUser.Name = reqDataInJSON.Name
	newUser.Place = reqDataInJSON.Place

	newUserAuth.UID = reqDataInJSON.UID
	// 密码加密入库
	newUserAuth.Pwd = base64.StdEncoding.EncodeToString([]byte(newUser.UID))

	// 更新人员数据库和认证数据库
	res := edumodel.AddUser(&newUser) && edumodel.AddUserAuth(&newUserAuth)
	if res {
		personAddReplyStatus = "success"
	} else {
		personAddReplyStatus = "model_fail"
	}

}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonAddRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonAddRouter: ", personAddReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(personAddReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PersonAddRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
