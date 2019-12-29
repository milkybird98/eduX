package edurouter

import (
	"encoding/base64"
	"fmt"
	"time"

	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"

	"github.com/tidwall/gjson"
)

type LoginRouter struct {
	edunet.BaseRouter
}

type LoginData struct {
	Pwd string `json:"pwd"`
}

var loginReplyStatus string

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

	pwdInBtye := []byte(loginData.String())
	//去盐
	pwdInBtye = pwdInBtye[7:]
	pwdInBtye[3] -= 2
	pwdInBtye[5] -= 3
	pwdInBtye[7] -= 7
	pwdInBtye[8] -= 11
	pwdInBtye[10] -= 13

	pwdInByteDecode := make([]byte, 64)

	_, err := base64.StdEncoding.Decode(pwdInByteDecode, pwdInBtye)
	if err != nil {
		registerReplyStatus = "pwd_format_error"
		return
	}

	authData := edumodel.GetUserAuthByUID(reqMsgInJSON.UID)
	if authData == nil {
		loginReplyStatus = "login_fail"
		return
	}

	if string(pwdInByteDecode) == authData.Pwd {
		loginReplyStatus = "success"
		c := request.GetConnection()

		c.SetSession("isLogined", true)
		c.SetSession("UID", userData.UID)
		c.SetSession("place", userData.Place)
	} else {
		loginReplyStatus = "fail"
	}
}

func (router *LoginRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", LoginRouter: ", loginReplyStatus)
	jsonMsg, err := CombineReplyMsg(loginReplyStatus, nil)
	if err != nil {
		fmt.Println("LoginRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
