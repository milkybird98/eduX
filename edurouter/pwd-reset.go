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

// PwdResetRouter 用于学生设置自己的密码,教师设置自己密码以及教师重置学生密码
type PwdResetRouter struct {
	edunet.BaseRouter
}

type PwdResetData struct {
	UID       string `json:"uid"`
	OriginPwd string `json:"oripwd"`
	Serect    string `json:"serect"`
	NewPwd    string `json:"newpwd"`
}

var pwdresetReplyStatus string

func (router *PwdResetRouter) PreHandle(request eduiface.IRequest) {
	// 数据验证
	var reqMsgInJSON *ReqMsg
	var ok bool

	// 检查数据格式
	reqMsgInJSON, registerReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查data段格式
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personputReplyStatus = "data_format_error"
		return
	}

	// 解析data段
	resetPwdData := gjson.ParseBytes(reqMsgInJSON.Data)

	// 解析uid段,不能为空
	uidData := resetPwdData.Get("uid")
	if !uidData.Exists() || uidData.String() == "" {
		registerReplyStatus = "uid_cannot_be_empty"
		return
	}
	uidInString := uidData.String()

	var newPwdInString string
	// 更新自己的密码
	if uidInString == reqMsgInJSON.UID {
		// 获取新密码
		newPwdData := resetPwdData.Get("newpwd")
		if !newPwdData.Exists() || newPwdData.String() == "" {
			registerReplyStatus = "new_password_cannot_be_empty"
			return
		}

		// 新密码去盐
		var err error
		newPwdInString, err = PwdRemoveSalr([]byte(newPwdData.String()))
		if err != nil {
			registerReplyStatus = err.Error()
			return
		}
	}

	// 获取旧密码或serect
	pwdData := resetPwdData.Get("oripwd")
	serectData := resetPwdData.Get("serect")

	// 二者必须有其一
	if (!pwdData.Exists() || pwdData.String() == "") && (!serectData.Exists() || serectData.String() == "") {
		registerReplyStatus = "origin_password_or_serect_must_have_one"
		return
	}

	// 如果旧密码存在则去盐
	var pwdInString string
	pwdInString = ""
	if pwdData.Exists() && pwdData.String() != "" {
		var err error
		pwdInString, err = PwdRemoveSalr([]byte(pwdData.String()))
		if err != nil {
			registerReplyStatus = err.Error()
			return
		}
	}

	// 如果serect存在则获取
	var serectInString string
	serectInString = ""
	if serectData.Exists() && serectData.String() != "" {
		serectInString = serectData.String()
	}

	// 权限检查
	// 当前连接用户数据检查
	user := edumodel.GetUserByUID(reqMsgInJSON.UID)
	if user == nil {
		registerReplyStatus = "user_not_found"
		return
	}

	// 如果当前连接用户为学生,则只能更新自己的密码
	if user.Place == "student" {
		if uidInString != reqMsgInJSON.UID {
			registerReplyStatus = "permission_error"
			return
		}
		// 如果当前连接用户为教师
	} else if user.Place == "teacher" {
		// 更新学生的密码
		if uidInString != reqMsgInJSON.UID {
			// 获取要更改对象所在班级
			class := edumodel.GetClassByUID(uidInString, "student")
			if class == nil {
				registerReplyStatus = "user_not_found"
				return
			}

			// 检查是否在同一班级
			ok := edumodel.CheckUserInClass(class.ClassName, reqMsgInJSON.UID, "teacher")
			if !ok {
				registerReplyStatus = "permission_error"
				return
			}
		}
	}

	// 获取要更改对象的验证数据
	userAuth := edumodel.GetUserAuthByUID(uidInString)
	if userAuth == nil {
		registerReplyStatus = "user_not_found"
		return
	}

	// 身份验证
	if user.Place == "student" {
		if pwdInString != "" {
			if pwdInString != userAuth.Pwd {
				registerReplyStatus = "password_wrong"
				return
			}
		} else if serectInString != "" {
			cache, err := utils.GetResetPasswordCache(serectInString)
			if err != nil || cache == nil || cache.UID != reqMsgInJSON.UID {
				registerReplyStatus = "serect_not_found"
				return
			}
		} else {
			registerReplyStatus = "auth_fail"
			return
		}
	} else if user.Place == "teacher" {
		if uidInString != reqMsgInJSON.UID {
			teacherAuth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
			if pwdInString != teacherAuth.Pwd {
				registerReplyStatus = "password_wrong"
				return
			} else {
				newPwdInString = base64.StdEncoding.EncodeToString([]byte(uidInString))
			}
		} else {
			if pwdInString != "" {
				if pwdInString != userAuth.Pwd {
					registerReplyStatus = "password_wrong"
					return
				}
			} else if serectInString != "" {
				cache, err := utils.GetResetPasswordCache(serectInString)
				if err != nil || cache == nil || cache.UID != reqMsgInJSON.UID {
					registerReplyStatus = "serect_not_found"
					return
				}
			} else {
				registerReplyStatus = "auth_fail"
				return
			}
		}
	}

	ok = edumodel.UpdateUserAuthByUID(uidInString, newPwdInString, "", "", "", "", "", "")
	if ok {
		registerReplyStatus = "success"
	} else {
		registerReplyStatus = "model_fail"
	}
}

func (router *PwdResetRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdResetRouter: ", registerReplyStatus)
	jsonMsg, err := CombineReplyMsg(registerReplyStatus, nil)
	if err != nil {
		fmt.Println("PwdResetRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}