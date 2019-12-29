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

type RegisterRouter struct {
	edunet.BaseRouter
}

type RegisterData struct {
	UID string `json:"uid"`
	Pwd []byte `json:"pwd"`
}

var registerReplyStatus string

func (router *RegisterRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool

	var registerTimer utils.RegisterTimerTag
	registerTimer.IP = request.GetConnection().GetTCPConnection().RemoteAddr()

	timer, _ := utils.GetRegisterTimerCache(registerTimer.IP.String())
	if timer != nil {
		registerReplyStatus = "try_register_too_fast"
		return
	}

	utils.SetRegisterTimerCacheExpire(registerTimer.IP.String(), registerTimer)

	reqMsgInJSON, registerReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personputReplyStatus = "data_format_error"
		return
	}

	registerData := gjson.ParseBytes(reqMsgInJSON.Data)

	uidData := registerData.Get("uid")
	if !uidData.Exists() || uidData.String() == "" {
		registerReplyStatus = "uid_cannot_be_empty"
		return
	}
	uidInString := uidData.String()

	pwdData := registerData.Get("pwd")
	if !pwdData.Exists() || pwdData.String() == "" {
		registerReplyStatus = "pwd_cannot_be_empty"
		return
	}
	pwdInBtye := []byte(pwdData.String())

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

	user := edumodel.GetUserByUID(uidInString)
	if user != nil {
		registerReplyStatus = "same_uid_exist"
		return
	}

	var newUser edumodel.User
	newUser.UID = uidInString
	newUser.Place = "student"

	var newUserAuth edumodel.UserAuth
	newUserAuth.UID = uidInString
	newUserAuth.Pwd = string(pwdInByteDecode)

	res := edumodel.AddUser(&newUser) && edumodel.AddUserAuth(&newUserAuth)
	if res {
		registerReplyStatus = "success"
		c := request.GetConnection()
		c.SetSession("isLogined", true)
		c.SetSession("UID", newUser.UID)
		c.SetSession("place", newUser.Place)
	} else {
		registerReplyStatus = "model_fail"
	}

}

func (router *RegisterRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", [ROUTER] Time: ", time.Now(), " Client address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), "RegisterRouter: ", registerReplyStatus)
	jsonMsg, err := CombineReplyMsg(registerReplyStatus, nil)
	if err != nil {
		fmt.Println("RegisterRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
