package edurouter

import (
	"crypto/md5"
	"eduX/eduiface"
	"eduX/utils"
	"encoding/json"
	"encoding/base64"

	"github.com/tidwall/gjson"
)

var checksumFlag bool

type ReqMsg struct {
	UID      string `json:"uid"`
	Data     []byte	`json:"data"`
	CheckSum string	`json:"checksum"`
}

type ResMsg struct {
	Status   string	`json:"status"`
	Data     []byte	`json:"data"`
	Checksum string	`json:"checksum"`
}

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
		reqMsgInJSON.Data,err = base64.StdEncoding.DecodeString(dataInMsg.String())
		if err!=nil{
			return nil, "data_base64_format_error", false
		}
	}else{
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

func CheckConnectionLogin(request eduiface.IRequest) (string, bool) {
	c := request.GetConnection()
	value, err := c.GetSession("isLogined")
	if err != nil {
		return "session_error", false
	}

	if value == false {
		return "not_login", false
	}

	return "", true
}

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
