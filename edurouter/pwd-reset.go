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

// PwdResetData 定了重设密码时的请求参数
type PwdResetData struct {
	UID       string `json:"uid"`
	OriginPwd string `json:"oripwd"`
	Serect    string `json:"serect"`
	NewPwd    string `json:"newpwd"`
}

// 返回状态码
var pwdresetReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
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
		registerReplyStatus = "data_format_error"
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
	fmt.Println(string(reqMsgInJSON.Data))
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
	var passwordResetNews *edumodel.News
	passwordResetNews = nil
	// 如果当前连接用户是学生
	if user.Place == "student" {
		// 如果填写了原密码则判断原密码是否一致
		if pwdInString != "" {
			// 如果密码不一致则报错返回
			if pwdInString != userAuth.Pwd {
				registerReplyStatus = "password_wrong"
				return
			}
		} else if serectInString != "" { // 如果填写了serect则判断serect是否存在
			cache, err := utils.GetResetPasswordCache(serectInString)
			// 如果serect不存在或者对应uid不一致,则报错返回
			if err != nil || cache == nil || cache.UID != reqMsgInJSON.UID {
				registerReplyStatus = "serect_not_found"
				return
			}
		} else { // 否则授权失败
			registerReplyStatus = "auth_fail"
			return
		}
	} else if user.Place == "teacher" || user.Place == "manager" { // 如果用户是教师或者管理员
		if uidInString != reqMsgInJSON.UID { // 如果修改他人密码
			userAuth := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
			// 检查密码填写是否正确
			if pwdInString != userAuth.Pwd {
				registerReplyStatus = "password_wrong"
				return
			}
			// 重置密码
			newPwdInString = base64.StdEncoding.EncodeToString([]byte(uidInString))

			// 想用户发送提醒消息
			var newNews edumodel.News
			newNews.SenderUID = reqMsgInJSON.UID
			newNews.AudientUID = []string{uidInString}
			newNews.SendTime = time.Now()
			newNews.NewsType = 2
			newNews.Title = "密码重置成功"
			newNews.Text = "你好,你的密码已经重置完成,请及时修改密码,以防他人恶意登陆."

			passwordResetNews = &newNews

		} else { // 修改自己的密码
			if pwdInString != "" { // 如果填写了原密码则判断原密码是否一致
				if pwdInString != userAuth.Pwd {
					registerReplyStatus = "password_wrong"
					return
				}
			} else if serectInString != "" { // 如果填写了serect则判断serect是否存在
				cache, err := utils.GetResetPasswordCache(serectInString)
				if err != nil || cache == nil || cache.UID != reqMsgInJSON.UID {
					registerReplyStatus = "serect_not_found"
					return
				}
			} else { // 否则授权失败
				registerReplyStatus = "auth_fail"
				return
			}
		}
	}

	// 更新用户授权数据库
	ok = edumodel.UpdateUserAuthByUID(uidInString, newPwdInString, "", "", "", "", "", "")
	if ok {
		// 如果有新的消息需要发送
		if passwordResetNews != nil {
			// 更新消息数据库
			ok = edumodel.AddNews(passwordResetNews)
			if ok {
				registerReplyStatus = "success"
			} else {
				registerReplyStatus = "model_fail"
			}
		} else {
			registerReplyStatus = "success"
		}
	} else {
		registerReplyStatus = "model_fail"
	}

}

// Handle 用于将请求的处理结果发回客户端
func (router *PwdResetRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PwdResetRouter: ", registerReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(registerReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("PwdResetRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
