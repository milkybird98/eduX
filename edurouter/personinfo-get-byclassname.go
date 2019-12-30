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

type PersonInfoGetByClassRouter struct {
	edunet.BaseRouter
}

type PersonInfoGetByClassData struct {
	ClassName string `json:"class"`
}

type PersonInfoGetByClassReplyData struct {
	UserList []PersonInfoGetReplyData `json:"userlist"`
}

var persongetbyclassReplyStatus string
var persongetbyclassReplyData PersonInfoGetByClassReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonInfoGetByClassRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, persongetbyclassReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	persongetbyclassReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	reqClassNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	if !reqClassNameData.Exists() {
		persongetbyclassReplyStatus = "data_format_error"
		return
	}

	reqClassName := reqClassNameData.String()

	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		persongetbyclassReplyStatus = "59session_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		persongetbyclassReplyStatus = "session_place_data_error"
		return
	}

	if placeString == "student" {
		persongetbyclassReplyStatus = "permission_error"
		return
	} else if placeString == "teacher" {
		class := edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)
		if class == nil {
			persongetbyclassReplyStatus = "model_fail"
			return
		}

		className := class.ClassName
		if className == "" {
			persongetbyclassReplyStatus = "not_in_class"
			return
		}

		if reqClassName != className {
			persongetbyclassReplyStatus = "permission_error"
			return
		}
	}

	userManyData := edumodel.GetUserByClass(reqClassName)
	if userManyData == nil || len(*userManyData) <= 0 {
		persongetbyclassReplyStatus = "data_not_found"
		return
	}
	persongetbyclassReplyStatus = "success"

	for _, person := range *userManyData {
		var personData PersonInfoGetReplyData
		personData.UID = person.UID
		personData.Name = person.Name
		personData.ClassName = person.Class
		personData.Gender = person.Gender
		personData.Birth = person.Birth
		personData.Political = person.Political
		if person.IsContactPub {
			personData.Contact = person.Contact
		} else {
			personData.Contact = "未公开"
		}
		if person.IsEmailPub {
			personData.Email = person.Email
		} else {
			personData.Email = "未公开"
		}
		if person.IsLocationPub {
			personData.Location = person.Location
		} else {
			personData.Location = "未公开"
		}

		persongetbyclassReplyData.UserList = append(persongetbyclassReplyData.UserList, personData)
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonInfoGetByClassRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonInfoGetByClassRouter: ", persongetbyclassReplyStatus)
	var jsonMsg []byte
	var err error

	if persongetbyclassReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(persongetbyclassReplyStatus, persongetbyclassReplyData)

	} else {
		jsonMsg, err = CombineReplyMsg(persongetbyclassReplyStatus, nil)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
