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

type NewsAddRouter struct {
	edunet.BaseRouter
}

type NewsAddData struct {
	IsAnnounce bool     `json:"isannounce"`
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	AudientUID []string `json:"audients"`
}

var newsaddReplyStatus string

func (router *NewsAddRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, newsaddReplyStatus, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	newsaddReplyStatus, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		newsaddReplyStatus = "data_format_error"
		return
	}

	newNewsData := gjson.ParseBytes(reqMsgInJSON.Data)

	titleData := newNewsData.Get("title")
	if !titleData.Exists() {
		newsaddReplyStatus = "title_cannot_be_empty"
		return
	}

	textData := newNewsData.Get("text")
	if !textData.Exists() {
		newsaddReplyStatus = "title_cannot_be_empty"
		return
	}

	audientData := newNewsData.Get("audients")

	isannounceData := newNewsData.Get("isannounce")
	var isannounce bool
	if !isannounceData.Exists() {
		isannounce = false
	} else {
		isannounce = isannounceData.Bool()
	}

	// 权限检查
	c := request.GetConnection()
	sessionPlace, err := c.GetSession("place")
	if err != nil {
		newsaddReplyStatus = "session_error"
		return
	}

	if sessionPlace != "teacher" && sessionPlace != "manager" {
		newsaddReplyStatus = "permission_error"
		return
	}

	placeString, ok := sessionPlace.(string)
	if ok != true {
		filegetbytagsReplyStatus = "session_place_data_error"
		return
	}

	class := edumodel.GetClassByUID(reqMsgInJSON.UID, placeString)
	if class == nil {
		filegetbytagsReplyStatus = "model_fail"
		return
	}

	className := class.ClassName
	if className == "" {
		filegetbytagsReplyStatus = "not_in_class"
		return
	}

	// 拼接数据
	var newNews edumodel.News
	newNews.SenderUID = reqMsgInJSON.UID
	newNews.SendTime = time.Now().In(utils.GlobalObject.TimeLocal)
	newNews.Title = titleData.String()
	newNews.Text = textData.String()
	newNews.IsAnnounce = isannounce

	if audientData.Exists() && audientData.IsArray() && len(audientData.Array()) > 0 {
		if sessionPlace == "manager" {
			for _, audient := range audientData.Array() {
				if audient.String() != "" {
					newNews.AudientUID = append(newNews.AudientUID, audient.String())
				} else {
					newsaddReplyStatus = "audient_UID_cannot_be_empty"
					return
				}
			}
		} else if sessionPlace == "teacher" {
			for _, audient := range audientData.Array() {
				if audient.String() != "" {
					if edumodel.CheckUserInClass(className, audient.String(), "student") {
						newNews.AudientUID = append(newNews.AudientUID, audient.String())
					} else {
						newsaddReplyStatus = "permission_error_cannot_send_news_to_another_class"
						return
					}
				} else {
					newsaddReplyStatus = "audient_UID_cannot_be_empty"
					return
				}
			}
		}
	} else {
		if sessionPlace == "manager" {
			newNews.AudientUID[0] = "all"
		} else if sessionPlace == "teacher" {
			class := edumodel.GetClassByUID(reqMsgInJSON.UID, "teacher")
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

// Handle 返回处理结果
func (router *NewsAddRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ", time.Now().In(utils.GlobalObject.TimeLocal).Format(utils.GlobalObject.TimeFormat), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", NewsAddRouter: ", newsaddReplyStatus)
	jsonMsg, err := CombineReplyMsg(newsaddReplyStatus, nil)
	if err != nil {
		fmt.Println("NewsAddRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
