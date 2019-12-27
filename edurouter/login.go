package edurouter

import (
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

	reqPwd := loginData.String()

	if reqPwd == userData.Pwd {
		loginReplyStatus = "success"
		c := request.GetConnection()

		c.SetSession("isLogined", true)
		c.SetSession("UID", userData.UID)
		c.SetSession("place", userData.Place)
		c.SetSession("class", userData.Class)
	} else {
		loginReplyStatus = "fail"
	}
}

func (router *LoginRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] Time: ", time.Now(), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", LoginRouter: ", loginReplyStatus)
	jsonMsg, err := CombineReplyMsg(loginReplyStatus, nil)
	if err != nil {
		fmt.Println("LoginRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
