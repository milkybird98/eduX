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

// PersonInfoPutRouter 处理人员信息更新请求
type PersonInfoPutRouter struct {
	edunet.BaseRouter
}

// PersonInfoPutData 定义人员信息更新参数
type PersonInfoPutData struct {
	UID          string `json:"uid"`
	Name         string `json:"name"`
	Gender       int    `json:"gender,omitempty"`
	Birth        string `json:"birthday,omitempty"`
	Political    int    `json:"polit,omitempty"`
	Contact      string `json:"contact"`
	IsContactPub bool   `json:"public"`
	Email        string `json:"email,omitempty"`
	Localion     string `json:"local,omitempty"`
	Job          string `json:"job,omitempty"`
	Notmod       string `json:"notmod"`
	Com1A        string `json:"com1a"`
	Com1B        string `json:"com1b"`
	Com2A        string `json:"com2a"`
	Com2B        string `json:"com2b"`
	Com3A        string `json:"com3a"`
	Com3B        string `json:"com3b"`
	Com4A        string `json:"com4a"`
	Com4B        string `json:"com4b"`
}

// 返回状态码
var personputReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonInfoPutRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, personputReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	personputReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personputReplyStatus = "data_format_error"
		return
	}

	newPersonInfoData := gjson.ParseBytes(reqMsgInJSON.Data)
	// 试图获取人员uid数据
	UID := newPersonInfoData.Get("uid").String()
	// 若不存在则返回错误码
	if UID == "" {
		personputReplyStatus = "uid_cannot_be_empty"
		return
	}

	//权限检查
	c := request.GetConnection()

	// 查询要更新数据的用户是否存在
	userData := edumodel.GetUserByUID(UID)
	// 不存在则返回错误码
	if userData == nil {
		personputReplyStatus = "user_not_found"
		return
	}

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		personputReplyStatus = err.Error()
		return
	}

	// 如果是要修改他人数据
	if reqMsgInJSON.UID != UID {
		// 如果是学生则权限错误
		if placeString == "student" {
			personputReplyStatus = "permission_error"
			return
		} else if placeString == "teacher" {
			// 如果是教师则要求在同一班级
			ok := edumodel.CheckUserInClass(userData.Class, reqMsgInJSON.UID, placeString)
			if !ok {
				personputReplyStatus = "permission_error"
				return
			}
		}
	}

	// 拼接更新数据
	var newUserInfo edumodel.User
	newUserInfo.UID = UID
	newUserInfo.Name = newPersonInfoData.Get("name").String()
	newUserInfo.Gender = int(newPersonInfoData.Get("gender").Int())
	newUserInfo.Birth = newPersonInfoData.Get("birthday").String()
	newUserInfo.Political = int(newPersonInfoData.Get("polit").Int())
	newUserInfo.Contact = newPersonInfoData.Get("contact").String()
	newUserInfo.IsContactPub = newPersonInfoData.Get("public").Bool()
	newUserInfo.Email = newPersonInfoData.Get("email").String()
	newUserInfo.IsEmailPub = newPersonInfoData.Get("public").Bool()
	newUserInfo.Localion = newPersonInfoData.Get("local").String()
	newUserInfo.IsLocalionPub = newPersonInfoData.Get("public").Bool()

	if reqMsgInJSON.UID != UID {
		if newUserInfo.IsContactPub == false {
			newUserInfo.Contact = ""
			newUserInfo.Email = ""
			newUserInfo.Localion = ""
		}
	}

	newUserInfo.Job = newPersonInfoData.Get("job").String()

	newUserInfo.Com1A = newPersonInfoData.Get("com1a").String()
	newUserInfo.Com1B = newPersonInfoData.Get("com1b").String()
	newUserInfo.Com2A = newPersonInfoData.Get("com2a").String()
	newUserInfo.Com2B = newPersonInfoData.Get("com2b").String()
	newUserInfo.Com3A = newPersonInfoData.Get("com3a").String()
	newUserInfo.Com3B = newPersonInfoData.Get("com3b").String()
	newUserInfo.Com4A = newPersonInfoData.Get("com4a").String()
	newUserInfo.Com4B = newPersonInfoData.Get("com4b").String()

	// 更新数据库
	res := edumodel.UpdateUserByID(&newUserInfo, newPersonInfoData.Get("notmod").Bool())
	// 若成功则返回success,否则返回错误码
	if res {
		personputReplyStatus = "success"
	} else {
		personputReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonInfoPutRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonInfoPutRouter: ", personputReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(personputReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PersonInfoPutRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
