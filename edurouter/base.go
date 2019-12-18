package edurouter

import (
	"crypto/md5"
	"eduX/eduiface"
	"eduX/utils"
	"encoding/json"
)

var checksumFlag bool

type ReqMsg struct {
	uid      string
	data     []byte
	checksum []byte
}

type ResMsg struct {
	status   string
	data     []byte
	checksum []byte
}

func CheckMsgFormat(request eduiface.IRequest) (*ReqMsg, string, bool) {
	var reqMsgInJSON ReqMsg
	reqMsgOrigin := request.GetData()

	checksumFlag = false

	err := json.Unmarshal(reqMsgOrigin, &reqMsgInJSON)
	if err != nil {
		return nil, "json_format_error", false
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(reqMsgInJSON.uid))
	md5Ctx.Write(reqMsgInJSON.data)

	if utils.SliceEqual(reqMsgInJSON.checksum, md5Ctx.Sum(nil)) != true {
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

	replyMsg.status = status

	if dataInJSON == nil {
		replyMsg.data = nil
	} else {
		data, err := json.Marshal(dataInJSON)
		if err != nil {
			return nil, err
		}
		replyMsg.data = data
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(replyMsg.status))
	md5Ctx.Write(replyMsg.data)
	replyMsg.checksum = md5Ctx.Sum(nil)

	jsonMsg, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, err
	}

	return jsonMsg, nil
}
