package edurouter

import (
	"fmt"
	"time"

	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/utils"

	"github.com/tidwall/gjson"
)

// LoginRouter 处理用户登陆请求
type LoginRouter struct {
	edunet.BaseRouter
}

// LoginData 定义登陆请求数据的参数
type LoginData struct {
	Pwd string `json:"pwd"`
}

// 回复状态码
var loginReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *LoginRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, loginReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	userData := edumodel.GetUserByUID(reqMsgInJSON.UID)
	if userData == nil {
		loginReplyStatus = "login_fail"
		return
	}

	// 尝试从请求数据Data段获取密码数据
	loginData := gjson.GetBytes(reqMsgInJSON.Data, "pwd")
	// 若数据不存在,则返回错误码
	if !loginData.Exists() || loginData.String() == "" {
		loginReplyStatus = "data_format_error"
		return
	}

	// 密码去盐
	pwdInByte := []byte(loginData.String())
	pwdInDecode, err := PwdRemoveSalr(pwdInByte)
	// 若密码解码失败,则返回错误码
	if err != nil {
		loginReplyStatus = "pwd_format_error"
		return
	}

	// 从数据库查询授权数据
	authData := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	// 如果授权数据不存在则返回错误码
	if authData == nil {
		loginReplyStatus = "login_fail"
		return
	}

	// 如果认证成功,则返回success
	if pwdInDecode == authData.Pwd {
		loginReplyStatus = "success"
		c := request.GetConnection()

		// 修改session
		c.SetSession("isLogined", true)
		c.SetSession("UID", userData.UID)
		c.SetSession("place", userData.Place)
	} else {
		// 否则返回错误码
		loginReplyStatus = "fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *LoginRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", LoginRouter: ", loginReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(loginReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("LoginRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
