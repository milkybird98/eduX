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

type RegisterRouter struct {
	edunet.BaseRouter
}

type RegisterData struct {
	UID string `json:"uid"`
	Pwd []byte `json:"pwd"`
}

var registerReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
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
	pwdInByte := []byte(pwdData.String())
	pwdInDecode, err := PwdRemoveSalr(pwdInByte)

	if err != nil {
		loginReplyStatus = "pwd_format_error"
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
	newUserAuth.Pwd = pwdInDecode

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

// Handle 用于将请求的处理结果发回客户端
func (router *RegisterRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", [ROUTER] Time: ", time.Now().In(utils.GlobalObject.TimeLocal), " Client address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), "RegisterRouter: ", registerReplyStatus)
	jsonMsg, err := CombineReplyMsg(registerReplyStatus, nil)
	if err != nil {
		fmt.Println("RegisterRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
