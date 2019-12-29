package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

type PersonInfoPutRouter struct {
	edunet.BaseRouter
}

type PersonInfoPutData struct {
	UID           string `json:"uid"`
	Name          string `json:"name"`
	Gender        int    `json:"gender,omitempty"`
	Birth         string `json:"birthday,omitempty"`
	Political     string `json:"polit,omitempty"`
	Contact       string `json:"contact"`
	IsContactPub  bool   `json:"isconpub"`
	Email         string `json:"email,omitempty"`
	IsEmailPub    bool   `json:"isemapub,omitempty"`
	Location      string `json:"locat,omitempty"`
	IsLocationPub bool   `json:"islocpub,omitempty"`
}

var personputReplyStatus string

func (router *PersonInfoPutRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, personputReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	personputReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personputReplyStatus = "data_format_error"
		return
	}

	newPersonInfoData := gjson.ParseBytes(reqMsgInJSON.Data)
	UID := newPersonInfoData.Get("uid").String()
	if UID == "" {
		personputReplyStatus = "uid_cannot_be_empty"
		return
	}

	userName := newPersonInfoData.Get("name").String()
	if userName == "" {
		personputReplyStatus = "name_cannot_be_empty"
		return
	}

	userContact := newPersonInfoData.Get("contact").String()
	if userContact == "" {
		personputReplyStatus = "contact_cannot_be_empty"
		return
	}

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		personputReplyStatus = "session_error"
		return
	}

	userData := edumodel.GetUserByUID(UID)
	if userData == nil {
		personputReplyStatus = "user_not_found"
		return
	}

	if sessionUID != UID {
		sessionPlace, err := c.GetSession("place")
		if err != nil {
			personputReplyStatus = "session_error"
			return
		}

		if sessionPlace == "student" {
			personputReplyStatus = "permission_error"
			return
		} else if sessionPlace == "teacher" {
			class := edumodel.GetClassByUID(reqMsgInJSON.UID, "teacher")
			if class == nil {
				personputReplyStatus = "not_in_class"
				return
			}

			if userData.Class != class.ClassName {
				personputReplyStatus = "permission_error"
				return
			}
		}
	}

	//修改个人信息
	var newUserInfo edumodel.User
	newUserInfo.UID = UID
	newUserInfo.Name = userName
	newUserInfo.Gender = int(newPersonInfoData.Get("gender").Int())
	newUserInfo.Birth = newPersonInfoData.Get("birthday").String()
	newUserInfo.Political = newPersonInfoData.Get("polit").String()
	newUserInfo.Contact = newPersonInfoData.Get("contact").String()
	newUserInfo.IsContactPub = newPersonInfoData.Get("isconpub").Bool()
	newUserInfo.Email = newPersonInfoData.Get("email").String()
	newUserInfo.IsEmailPub = newPersonInfoData.Get("isemapub").Bool()
	newUserInfo.Location = newPersonInfoData.Get("locat").String()
	newUserInfo.IsLocationPub = newPersonInfoData.Get("islocpub").Bool()

	res := edumodel.UpdateUserByID(&newUserInfo)
	if res {
		personputReplyStatus = "success"
	} else {
		personputReplyStatus = "model_fail"
	}
}

func (router *PersonInfoPutRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonInfoPutRouter: ", personputReplyStatus)
	jsonMsg, err := CombineReplyMsg(personputReplyStatus, nil)
	if err != nil {
		fmt.Println("PersonInfoPutRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
