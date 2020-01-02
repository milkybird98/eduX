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

// QuestionGetBySenderUIDRouter 处理根据发送者uid获取问题列表的请求
type QuestionGetBySenderUIDRouter struct {
	edunet.BaseRouter
}

// QuestionGetBySenderUIDData 定义根据发送者uid请求问题时的参数
type QuestionGetBySenderUIDData struct {
	SenderUID   string `json:"uid"`
	Skip        int64  `json:"skip"`
	Limit       int64  `json:"limit"`
	DeferSolved bool   `json:"defer"`
	IsSolved    bool   `json:"issolved"`
}

// QuestionGetBySenderUIDReplyData 定义根据发送者uid请求问题时的返回参数
type QuestionGetBySenderUIDReplyData struct {
	QuestionList *[]edumodel.Question `json:"questions"`
}

// 返回状态码
var questiongetbysenderuidReplyStatus string

// 返回数据
var questiongetbysenderuidReplyData QuestionGetBySenderUIDReplyData

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *QuestionGetBySenderUIDRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据,并检查校验和
	reqMsgInJSON, questiongetbysenderuidReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	questiongetbysenderuidReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		questiongetbysenderuidReplyStatus = "data_format_error"
		return
	}

	// 获取发送者数据
	senderUIDData := gjson.GetBytes(reqMsgInJSON.Data, "uid")
	// 如果发送者uid不存在则报错
	if !senderUIDData.Exists() || senderUIDData.String() == "" {
		questiongetbysenderuidReplyStatus = "senderuid_cannot_be_empty"
		return
	}
	// 转换为字符串
	senderUID := senderUIDData.String()

	// 获取skip和limit数据,若不存在则使用默认值
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
		questiongetbysenderuidReplyStatus = err.Error()
		return
	}

	// 检查请求用户数据是否存在
	user := edumodel.GetUserByUID(senderUID)
	// 若未找到目标用户则报错
	if user == nil {
		questiongetbysenderuidReplyStatus = "user_not_found"
		return
	}

	// 如果用户不是管理员
	if placeString != "manager" {
		// 检查当前连接用户是否和查询用户在同一班级
		ok := edumodel.CheckUserInClass(user.Class, reqMsgInJSON.UID, placeString)
		// 若不在则权限错误
		if !ok {
			questiongetbysenderuidReplyStatus = "permission_error"
			return
		}
	}

	// 查询数据库,获取问题列表
	questionList := edumodel.GetQuestionBySenderUID(int(Skip), int(Limit), DetectSolved, IsSolved, senderUID)
	// 查询成功则返回问题数据,并设定状态为success,否则返回错误码
	if questionList != nil {
		questiongetbysenderuidReplyStatus = "success"
		questiongetbysenderuidReplyData.QuestionList = questionList
	} else {
		questiongetbysenderuidReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *QuestionGetBySenderUIDRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", QuestionGetBySenderUIDRouter: ", questiongetbysenderuidReplyStatus)

	var jsonMsg []byte
	var err error

	// 生成返回数据
	if questiongetbysenderuidReplyStatus == "success" {
		jsonMsg, err = CombineReplyMsg(questiongetbysenderuidReplyStatus, questiongetbysenderuidReplyData)
	} else {
		jsonMsg, err = CombineReplyMsg(questiongetbysenderuidReplyStatus, nil)
	}
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("QuestionGetBySenderUIDRouter : ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
