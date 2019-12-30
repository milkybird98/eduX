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

type LoginRouter struct {
	edunet.BaseRouter
}

type LoginData struct {
	Pwd string `json:"pwd"`
}

var loginReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *LoginRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, loginReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	userData := edumodel.GetUserByUID(reqMsgInJSON.UID)
	if userData == nil {
		loginReplyStatus = "login_fail"
		return
	}

	loginData := gjson.GetBytes(reqMsgInJSON.Data, "pwd")
	if !loginData.Exists() {
		loginReplyStatus = "data_format_error"
		return
	}

	pwdInByte := []byte(loginData.String())
	pwdInDecode, err := PwdRemoveSalr(pwdInByte)
	if err != nil {
		loginReplyStatus = "pwd_format_error"
		return
	}

	authData := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	if authData == nil {
		loginReplyStatus = "login_fail"
		return
	}

	//fmt.Println([]byte(pwdInDecode))
	//fmt.Println([]byte(authData.Pwd))
	if pwdInDecode == authData.Pwd {
		loginReplyStatus = "success"
		c := request.GetConnection()

		c.SetSession("isLogined", true)
		c.SetSession("UID", userData.UID)
		c.SetSession("place", userData.Place)
	} else {
		loginReplyStatus = "fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *LoginRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", LoginRouter: ", loginReplyStatus)
	jsonMsg, err := CombineReplyMsg(loginReplyStatus, nil)
	if err != nil {
		fmt.Println("LoginRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
