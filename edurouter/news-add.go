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

// NewsAddRouter 处理用户添加新消息的请求
type NewsAddRouter struct {
	edunet.BaseRouter
}

// NewsAddData 定义了添加新消息的参数
type NewsAddData struct {
	IsAnnounce bool     `json:"isannounce"` // 是否是公告
	Title      string   `json:"title"`      // 消息标题
	Text       string   `json:"text"`       // 消息正文
	AudientUID []string `json:"audients"`   // 消息接收者
}

// 返回状态码
var newsaddReplyStatus string

// PreHandle 用于进行原始数据校验,权限验证,身份验证,数据获取和数据库操作
func (router *NewsAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	// 试图解码原始数据
	reqMsgInJSON, newsaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	// 检查当前连接是否已登录
	newsaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	// 验证请求数据Data段格式是否正确
	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newsaddReplyStatus = "data_format_error"
		return
	}

	// 解码请求数据Data段
	newNewsData := gjson.ParseBytes(reqMsgInJSON.Data)

	// 从Data段中获取消息标题
	titleData := newNewsData.Get("title")
	// 如果消息标题不存在,则返回错误码
	if !titleData.Exists() {
		newsaddReplyStatus = "title_cannot_be_empty"
		return
	}

	// 从Data段中获取新消息正文
	textData := newNewsData.Get("text")
	// 如果消息正文不存在,则返回错误码
	if !textData.Exists() {
		newsaddReplyStatus = "title_cannot_be_empty"
		return
	}

	// 从Data中获取消息"听众"
	audientData := newNewsData.Get("audients")

	// 从Data段获取公告标志位,判断是否是公告
	isannounceData := newNewsData.Get("isannounce")
	var isannounce bool
	// 如果不存在,则认为默认是非公告
	if !isannounceData.Exists() {
		isannounce = false
	} else {
		isannounce = isannounceData.Bool()
	}

	// 权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		classdelReplyStatus = err.Error()
		return
	}

	// 如果当前用户不是教师或管理员
	if placeString != "teacher" && placeString != "manager" {
		newsaddReplyStatus = "permission_error"
		return
	}

	// 尝试获取班级数据
	class := edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)
	// 如果班级不存在则报错
	if class == nil {
		filegetbytagsReplyStatus = "not_in_class"
		return
	}

	// 拼接数据
	var newNews edumodel.News
	newNews.SenderUID = reqMsgInJSON.UID
	// 获取消息发送时间
	newNews.SendTime = time.Now().In(utils.GlobalObject.TimeLocal)
	newNews.Title = titleData.String()
	newNews.Text = textData.String()
	newNews.IsAnnounce = isannounce

	// 如果听众数据存在
	if audientData.Exists() && audientData.IsArray() && len(audientData.Array()) > 0 {
		// 如果当前用户是管理员
		if placeString == "manager" {
			// 将接收数据转换为字符串切片
			for _, audient := range audientData.Array() {
				if audient.String() != "" {
					newNews.AudientUID = append(newNews.AudientUID, audient.String())
				}
			}
			// 如果全部都是无效数据则报错返回
			if len(newNews.AudientUID) <= 0 {
				newsaddReplyStatus = "audient_UID_cannot_be_empty"
				return
			}
		} else if placeString == "teacher" { // 如果当前用户是教师
			// 将接收数据转换为字符串切片
			for _, audient := range audientData.Array() {
				// 检查添加用户是否在教师管理的班级中,若不在则返回错误码
				if audient.String() != "" {
					if edumodel.CheckUserInClass(class.ClassName, audient.String(), "student") {
						newNews.AudientUID = append(newNews.AudientUID, audient.String())
					} else {
						newsaddReplyStatus = "permission_error_cannot_send_news_to_another_class"
						return
					}
				} else {
					// 若存在无效数据则返回错误码
					newsaddReplyStatus = "audient_UID_cannot_be_empty"
					return
				}
			}
		}
	} else { //若干不存在听众数据
		if placeString == "manager" { //如果当前用户是管理员
			// 将听众设为所有人
			newNews.AudientUID[0] = "all"
		} else if placeString == "teacher" { // 如果当前用户是管理员
			class := edumodel.GetClassByUID(reqMsgInJSON.UID, "teacher")
			// 将听众设为所在班级的全部学生
			if class == nil {
				newsaddReplyStatus = "not_join_class"
				return
			}
			newNews.AudientUID = class.StudentList
		}
	}

	// 更新数据库
	ok = edumodel.AddNews(&newNews)
	if ok {
		newsaddReplyStatus = "success"
	} else {
		newsaddReplyStatus = "model_fail"
	}
}

// Handle 用于将请求的处理结果发回客户端
func (router *NewsAddRouter) Handle(request eduiface.IRequest) {
	// 打印请求处理Log
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsAddRouter: ", newsaddReplyStatus)
	// 生成返回数据
	jsonMsg, err := CombineReplyMsg(newsaddReplyStatus, nil)
	// 如果生成失败则报错返回
	if err != nil {
		fmt.Println("NewsAddRouter: ", err)
		return
	}

	// 发送返回数据
	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
