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

type PersonAddRouter struct {
	edunet.BaseRouter
}

type PersonAddData struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Place string `json:"place"`
}

var personAddReplyStatus string

func (router *PersonAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	var reqDataInJSON PersonAddData
	reqMsgInJSON, personAddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	personAddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		personAddReplyStatus = "data_format_error"
		return
	}

	newStudentData := gjson.ParseBytes(reqMsgInJSON.Data)
	reqDataInJSON.UID = newStudentData.Get("uid").String()
	reqDataInJSON.Name = newStudentData.Get("name").String()
	reqDataInJSON.Place = newStudentData.Get("place").String()

	//权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		personAddReplyStatus = "seesion_plcae_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		filegetbytagsReplyStatus = "session_place_data_error"
		return
	}

	if placeString != "manager" {
		personAddReplyStatus = "permission_error"
		return
	}

	userData := edumodel.GetUserByUID(reqDataInJSON.UID)
	if userData != nil {
		personAddReplyStatus = "same_uid_exist"
		return
	}

	//数据库操作
	var newUser edumodel.User
	var newUserAuth edumodel.UserAuth

	newUser.UID = reqDataInJSON.UID
	newUser.Name = reqDataInJSON.Name
	newUser.Place = reqDataInJSON.Place

	newUserAuth.UID = reqDataInJSON.UID
	newUserAuth.Pwd = base64.StdEncoding.EncodeToString([]byte(newUser.UID))

	res := edumodel.AddUser(&newUser) && edumodel.AddUserAuth(&newUserAuth)
	if res {
		personAddReplyStatus = "success"
	} else {
		personAddReplyStatus = "model_fail"
	}

}

func (router *PersonAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonAddRouter: ", personAddReplyStatus)
	jsonMsg, err := CombineReplyMsg(personAddReplyStatus, nil)
	if err != nil {
		fmt.Println("PersonAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
