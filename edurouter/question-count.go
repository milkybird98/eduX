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

// QuestionCountRouter 处理统计问题数据请求
type QuestionCountRouter struct {
	edunet.BaseRouter
}

// QuestionCountData 定义请求问题统计数据的参数
type QuestionCountData struct {
	ClassName string    `json:"classname"`
	Date      time.Time `json:"time"`
	IsSolved  bool      `json:"issolved"`
}

// QuestionCountReplyData 定义问题统计数据返回的参数
type QuestionCountReplyData struct {
	Number int `json:"num"`
}

// 返回状态码
var questioncountReplyStatus string

// 返回数据
var questioncountReplyData QuestionCountReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *QuestionCountRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, questioncountReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	questioncountReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questioncountReplyStatus = "data_format_error"
		return
	}

	// 权限检查
	c := request.GetConnection()

	// 获取当前用户身份
	userPlace, err := GetSessionPlace(c)
	if err != nil {
		questioncountReplyStatus = err.Error()
		return
	}

	// 如果当前用户不是管理员则权限错误
	if userPlace != "manager" {
		questioncountReplyStatus = "permission_error"
		return
	}

	// 获取查询的时间参数
	timeData := gjson.GetBytes(reqMsgInJSON.Data, "time")
	// 解码时间数据
	targetTime, err := time.Parse(time.RFC3339, timeData.String())
	var isTimeRequired bool
	// 如果成功解码出时间则限定统计时间
	if err != nil || targetTime.IsZero() {
		isTimeRequired = false
	} else {
		isTimeRequired = true
	}

	// 获取班级名
	className := gjson.GetBytes(reqMsgInJSON.Data, "classname").String()
	// 获取是否限定已解决问题标志位
	IsSolved := gjson.GetBytes(reqMsgInJSON.Data, "issolved").Bool()

	// 查询数据库
	if isTimeRequired {
		if IsSolved {
			questioncountReplyData.Number = edumodel.GetQuestionAnsweredNumberByDate(className, targetTime)
		} else {
			questioncountReplyData.Number = edumodel.GetQuestionNumberByDate(className, targetTime)
		}
	} else {
		if IsSolved {
			questioncountReplyData.Number = edumodel.GetQuestionAnsweredNumberByDate(className, targetTime)
		} else {
			questioncountReplyData.Number = edumodel.GetQuestionNumberByDate(className, targetTime)
		}
	}

	// 如果获取到有效数据则返回success,否则提示数据库错误
	if questioncountReplyData.Number != -1 {
		questioncountReplyStatus = "success"
	} else {
		questioncountReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *QuestionCountRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionCountRouter: ", questioncountReplyStatus)

	var jsonMsg []byte
	var err error

	// 生成返回数据
	if questioncountReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(questioncountReplyStatus, questioncountReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(questioncountReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("QuestionCountRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
