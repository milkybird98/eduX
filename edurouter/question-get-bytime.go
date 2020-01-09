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

// QuestionGetByTimeNameRouter 处理根据班级名称获取问题的请求
type QuestionGetByTimeNameRouter struct {
	edunet.BaseRouter
}

// QuestionGetByTimeNameData 定义根据问题获取班级名称时的请求参数
type QuestionGetByTimeNameData struct {
	Skip        int64 `json:"skip"`
	Limit       int64 `json:"limit"`
	DeferSolved bool  `json:"defer"`
	IsSolved    bool  `json:"issolved"`
}

// QuestionGetByTimeReplyData 定义根据班级名称查询问题数据的返回参数
type QuestionGetByTimeReplyData struct {
	QuestionList *[]edumodel.Question `json:"questions"`
}

// 返回状态码
var questiongetbytimenameReplyStatus string

// 返回数据
var questiongetbytimenameReplyData QuestionGetByTimeReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *QuestionGetByTimeNameRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, questiongetbytimenameReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	questiongetbytimenameReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questiongetbytimenameReplyStatus = "data_format_error"
		return
	}

	// 获取班级名称数据
	classNameData := gjson.GetBytes(reqMsgInJSON.Data, "class")
	// 如果班级名称不存在则返回错误码
	if !classNameData.Exists() || classNameData.String() == "" {
		questiongetbytimenameReplyStatus = "classname_cannot_be_empty"
		return
	}
	// 将班级名称数据转换为字符串
	className := classNameData.String()

	// 获取skip和limit值,若不存在则使用默认值
	Skip, Limit := GetSkipAndLimit(reqMsgInJSON.Data)

	// 检查是否需要进行回答判定
	var DetectSolved bool
	// 如果需要判定项不存在,则默认不进行判定
	detectSolvedData := gjson.GetBytes(reqMsgInJSON.Data, "defer")
	if detectSolvedData.Exists() {
		DetectSolved = detectSolvedData.Bool()
	} else {
		DetectSolved = false
	}

	// 是否只查询已解答或是未解答的问题
	var IsSolved bool
	// 如果参数未指定,则默认查询未解答问题
	issolvedData := gjson.GetBytes(reqMsgInJSON.Data, "issolved")
	if issolvedData.Exists() {
		IsSolved = issolvedData.Bool()
	} else {
		IsSolved = false
	}

	//权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		questiongetbytimenameReplyStatus = err.Error()
		return
	}

	// 获取班级数据
	class := edumodel.GetClassByName(className)
	// 若班级不存在则报错返回
	if class == nil {
		questiongetbytimenameReplyStatus = "class_not_found"
		return
	}

	// 如果用户不是管理员
	if placeString != "manager" {
		// 检查当前用户是否在期望查询的班级中
		ok := edumodel.CheckUserInClass(className, reqMsgInJSON.UID, placeString)
		// 若不在则权限错误
		if !ok {
			questiongetbytimenameReplyStatus = "permission_error"
			return
		}
	}

	// 查询数据库,获取问题列表
	questionList := edumodel.GetQuestionByTimeOrder(Skip, Limit, DetectSolved, IsSolved)
	// 查询成功则返回问题数据,并设定状态为success,否则返回错误码
	if questionList != nil {
		questiongetbytimenameReplyStatus = "success"
		questiongetbytimenameReplyData.QuestionList = questionList
	} else {
		questiongetbytimenameReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *QuestionGetByTimeNameRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionGetByTimeNameRouter: ", questiongetbytimenameReplyStatus)

	var jsonMsg []byte
	var err error

	// 生成返回数据
	if questiongetbytimenameReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(questiongetbytimenameReplyStatus, questiongetbytimenameReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(questiongetbytimenameReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("QuestionGetByTimeNameRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
