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

// PersonInfoGetRouter 处理获取人员信息请求
type PersonInfoGetRouter struct {
	edunet.BaseRouter
}

// PersonInfoGetData 定义了请求人员信息时的参数
type PersonInfoGetData struct {
	UID string `json:"uid"`
}

// PersonInfoGetReplyData 定义了请求人员信息的返回参数
type PersonInfoGetReplyData struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	ClassName string `json:"class"`
	Gender    int    `json:"gender,omitempty"`
	Birth     string `json:"birthday,omitempty"`
	Political int    `json:"polit,omitempty"`
	Contact   string `json:"contact"`
	Email     string `json:"email,omitempty"`
	Location  string `json:"locat,omitempty"`
}

// 返回状态码
var persongetReplyStatus string

// 返回数据
var persongetReplyData PersonInfoGetReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *PersonInfoGetRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, persongetReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	persongetReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 从Data段获取人员uid数据
	personInfoGetData := gjson.GetBytes(reqMsgInJSON.Data, "uid")
	// 若不存在则返回错误码
	if !personInfoGetData.Exists() || personInfoGetData.String() == "" {
		persongetReplyStatus = "data_format_error"
		return
	}
	personUID := personInfoGetData.String()

	//权限验证

	c := request.GetConnection()
	var userData *edumodel.User

	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		persongetReplyStatus = err.Error()
		return
	}

	// 查询数据库
	userData = edumodel.GetUserByUID(personUID)
	// 若查询用户不存在,返回错误码
	if userData == nil {
		persongetReplyStatus = "user_not_found"
		return
	}

	// 如果用户不是管理员且不在请求用户的班级,则权限错误
	if reqMsgInJSON.UID != personUID {
		if placeString == "teacher" || placeString == "student" {
			ok := edumodel.CheckUserInClass(userData.Class, reqMsgInJSON.UID, placeString)
			if !ok {
				persongetReplyStatus = "permission_error"
				return
			}
		}
	}

	// 保护隐私数据
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

	// 设定返回状态码
	persongetReplyStatus = "success"
}

// Handle 用于将请求的处理结果发回客户端
func (router *PersonInfoGetRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", PersonInfoGetRouter: ", persongetReplyStatus)
	var jsonMsg []byte
	var err error
	// 生成返回数据
	if persongetReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(persongetReplyStatus, persongetReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(persongetReplyStatus, nil)

	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println(err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
