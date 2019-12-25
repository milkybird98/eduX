package edurouter

import (
	"crypto/md5"
	"eduX/eduiface"
	"eduX/utils"
	"encoding/base64"
	"encoding/json"

	"github.com/tidwall/gjson"
)

var checksumFlag bool

//ReqMsg 接收数据结构体
type ReqMsg struct {
	UID      string `json:"uid"`
	Data     []byte `json:"data"`
	CheckSum string `json:"checksum"`
}

//ResMsg 发送数据结构体
type ResMsg struct {
	Status   string `json:"status"`
	Data     []byte `json:"data"`
	Checksum string `json:"checksum"`
}

//CheckMsgFormat 检查接收数据格式是否正确
func CheckMsgFormat(request eduiface.IRequest) (*ReqMsg, string, bool) {
	var reqMsgInJSON ReqMsg
	reqMsgOrigin := request.GetData()

	checksumFlag = false

	if !gjson.Valid(string(reqMsgOrigin)) {
		return nil, "json_format_error", false
	}

	parseResult := gjson.ParseBytes(reqMsgOrigin)

	reqMsgInJSON.UID = parseResult.Get("uid").String()
	dataInMsg := parseResult.Get("data")
	if dataInMsg.Exists() {
		var err error
		reqMsgInJSON.Data, err = base64.StdEncoding.DecodeString(dataInMsg.String())
		if err != nil {
			return nil, "data_base64_format_error", false
		}
	} else {
		reqMsgInJSON.Data = nil
	}
	reqMsgInJSON.CheckSum = gjson.GetBytes(reqMsgOrigin, "checksum").String()

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJSON.UID))
	md5Ctx.Write([]byte(reqMsgInJSON.Data))

	if utils.SliceEqual([]byte(reqMsgInJSON.CheckSum), md5Ctx.Sum(nil)) != true {
		return nil, "check_sum_error", false
	}
	return &reqMsgInJSON, "", true
}

//CheckConnectionLogin 检查当前用户是否已登陆
func CheckConnectionLogin(request eduiface.IRequest, UID string) (string, bool) {
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		return "session_error", false
	}

	if sessionUID != UID {
		return "uid_not_match", false
	}

	value, err := c.GetSession("isLogined")
	if err != nil {
		return "session_error", false
	}

	if value == false {
		return "not_login", false
	}

	return "", true
}

//CombineReplyMsg 拼接,计算MD5并且序列化返回数据
func CombineReplyMsg(status string, dataInJSON interface{}) ([]byte, error) {
	var replyMsg ResMsg
	var err error

	replyMsg.Status = status

	if dataInJSON == nil {
		replyMsg.Data = nil
	} else {
		data, err := json.Marshal(dataInJSON)
		if err != nil {
			return nil, err
		}
		replyMsg.Data = data
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(replyMsg.Status))
	md5Ctx.Write(replyMsg.Data)
	replyMsg.Checksum = string(md5Ctx.Sum(nil))

	jsonMsg, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, err
	}

	return jsonMsg, nil
}
