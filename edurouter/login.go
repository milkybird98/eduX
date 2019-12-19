package edurouter

import (
	"fmt"

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

var login_replyStatus string

func (this *LoginRouter) PreHandle(request eduiface.IRequest) {
	reqMsgInJSON, login_replyStatus, ok := CheckMsgFormat(request)
	if ok != true {
		fmt.Println("LoginRouter: ", login_replyStatus)
		return
	}

	userData := edumodel.GetUserByUID(reqMsgInJSON.UID)
	if userData == nil {
		login_replyStatus = "login_fail"
		return
	}

	loginData := gjson.GetBytes(reqMsgInJSON.Data, "pwd")
	if !loginData.Exists() {
		login_replyStatus = "data_format_error"
		return
	}

	reqPwd := loginData.String()

	if reqPwd == userData.Pwd {
		login_replyStatus = "success"
		c := request.GetConnection()

		c.SetSession("isLogined", true)
		c.SetSession("UID", userData.UID)
		c.SetSession("place", userData.Plcae)
		c.SetSession("class", userData.Class)
	} else {
		login_replyStatus = "fail"
	}

}

func (this *LoginRouter) Handle(request eduiface.IRequest) {
	fmt.Println("LoginRouter: ", login_replyStatus)
	jsonMsg, err := CombineReplyMsg(login_replyStatus, nil)
	if err != nil {
		fmt.Println("LoginRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
