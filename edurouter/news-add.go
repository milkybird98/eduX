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
	IsAnnounce bool      `json:"type"`     // 是否是公告
	Title      string    `json:"title"`    // 消息标题
	Text       string    `json:"text"`     // 消息正文
	AudientUID []string  `json:"audients"` // 消息接收者
	TargetTime time.Time `json:"targettime,omitempty"`
	NewsType   int64     `json:""`
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
	if !titleData.Exists() || titleData.String() == "" {
		newsaddReplyStatus = "title_cannot_be_empty"
		return
	}

	// 从Data段中获取新消息正文
	textData := newNewsData.Get("text")
	// 如果消息正文不存在,则返回错误码
	if !textData.Exists() || textData.String() == "" {
		newsaddReplyStatus = "title_cannot_be_empty"
		return
	}

	// 从Data中获取消息"听众"
	audientData := newNewsData.Get("audients")

	// 从Data段获取公告标志位,判断是否是公告
	newsTypeData := newNewsData.Get("type")
	// 如果不存在,则认为默认是非公告
	if !newsTypeData.Exists() || newsTypeData.Int() < 1 || newsTypeData.Int() > 5 {
		newsaddReplyStatus = "type_cannot_be_empty"
		return
	}

	// 试图从Data段中获取日期数据
	timeData := gjson.GetBytes(reqMsgInJSON.Data, "targettime")
	var targetTime time.Time
	var err error
	// 解码时间数据
	if timeData.Exists() && timeData.String() != "" {
		targetTime, err = time.Parse(time.RFC3339, timeData.String())
		// 如果成功解码出时间则限定统计时间
		if err != nil || targetTime.IsZero() {
			targetTime = time.Now()
		}
	} else {
		targetTime = time.Now()
	}

	// 权限检查
	c := request.GetConnection()
	// 试图从session中获取身份数据
	placeString, err := GetSessionPlace(c)
	// 若不存在则返回
	if err != nil {
		newsaddReplyStatus = err.Error()
		return
	}

	// 如果当前用户不是教师或管理员
	if placeString != "teacher" && placeString != "manager" {
		newsaddReplyStatus = "permission_error"
		return
	}

	var class *edumodel.Class
	if placeString == "teacher" {
		// 尝试获取班级数据
		class = edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)
		// 如果班级不存在则报错
		if class == nil {
			newsaddReplyStatus = "not_in_class"
			return
		}
	}

	// 拼接数据
	var newNews edumodel.News
	newNews.SenderUID = reqMsgInJSON.UID
	// 获取消息发送时间
	newNews.SendTime = time.Now()
	newNews.Title = titleData.String()
	newNews.Text = textData.String()
	newNews.NewsType = newsTypeData.Int()

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
					if newNews.NewsType == 3 || newNews.NewsType == 5 {
						newNews.AudientUID = append(newNews.AudientUID, audient.String())
					} else {
						if edumodel.CheckUserInClass(class.ClassName, audient.String(), "student") {
							newNews.AudientUID = append(newNews.AudientUID, audient.String())
						} else {
							newsaddReplyStatus = "permission_error_cannot_send_news_to_another_class"
							return
						}
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
			if newNews.NewsType == 4 {
				// 将听众设为所有人
				newNews.AudientUID = []string{"all"}
			} else {
				newsaddReplyStatus = "audient_cannot_be_empty"
				return
			}
		} else if placeString == "teacher" { // 如果当前用户是管理员
			if newNews.NewsType == 3 || newNews.NewsType == 5 {
				newNews.AudientUID = []string{class.ClassName}
			} else {
				newsaddReplyStatus = "audient_cannot_be_empty"
				return
			}
		}
	}
	newNews.TargetTime = targetTime

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
	fmt.Println("[ROUTERS] ", time.Now().Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsAddRouter: ", newsaddReplyStatus)
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
