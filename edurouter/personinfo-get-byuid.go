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

type PersonInfoGetRouter struct {
	edunet.BaseRouter
}

type PersonInfoGetData struct {
	UID string `json:"uid"`
}

type PersonInfoGetReplyData struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	ClassName string `json:"class"`
	Gender    int    `json:"gender,omitempty"`
	Birth     string `json:"birthday,omitempty"`
	Political string `json:"polit,omitempty"`
	Contact   string `json:"contact"`
	Email     string `json:"email,omitempty"`
	Location  string `json:"locat,omitempty"`
}

var persongetReplyStatus string
var persongetReplyData PersonInfoGetReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonInfoGetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, persongetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}
	time.Now().In(utils.GlobalObject.TimeLocal).String()
	persongetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	personInfoGetData := gjson.GetBytes(reqMsgInJSON.Data, "uid")
	if !personInfoGetData.Exists() {
		persongetReplyStatus = "data_format_error"
		return
	}

	personUID := personInfoGetData.String()
	if personUID == "" {
		persongetReplyStatus = "uid_cannot_be_empty"
		return
	}

	c := request.GetConnection()
	var userData *edumodel.User

	sessionPlace, err := c.GetSession("place")
	if err != nil {
		persongetReplyStatus = "64session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		persongetReplyStatus = "session_place_data_error"
		return
	}

	userData = edumodel.GetUserByUID(personUID)
	if userData == nil {
		persongetReplyStatus = "user_not_found"
		return
	}

	if reqMsgInJSON.UID != personUID {
		if placeString == "teacher" || placeString == "student" {
			ok := edumodel.CheckUserInClass(userData.Class, reqMsgInJSON.UID, placeString)
			if !ok {
				persongetReplyStatus = "permission_error"
				return
			}
		}
	}

	persongetReplyData.UID = userData.UID
	persongetReplyData.Name = userData.Name
	persongetReplyData.ClassName = userData.Class
	persongetReplyData.Gender = userData.Gender
	persongetReplyData.Birth = userData.Birth
	persongetReplyData.Political = userData.Political
	if reqMsgInJSON.UID != personUID {
		if userData.IsContactPub {
			persongetReplyData.Contact = userData.Contact
		} else {
			persongetReplyData.Contact = "未公开"
		}
		if userData.IsEmailPub {
			persongetReplyData.Email = userData.Email
		} else {
			persongetReplyData.Email = "未公开"
		}
		if userData.IsLocationPub {
			persongetReplyData.Location = userData.Location
		} else {
			persongetReplyData.Location = "未公开"
		}
	} else {
		persongetReplyData.Contact = userData.Contact
		persongetReplyData.Email = userData.Email
		persongetReplyData.Location = userData.Location
	}

	persongetReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonInfoGetRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonInfoGetRouter: ", persongetReplyStatus)
	var jsonMsg []byte
	var err error
	if persongetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(persongetReplyStatus, persongetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(persongetReplyStatus, nil)

	}
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
