package edurouter

import (
	"bytes"
	"crypto/md5"
	"eduX/eduiface"
	"eduX/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
)

var checksumFlag bool

//ReqMsg 接收数据结构体
type ReqMsg struct {
	UID      string `json:"uid"`
	Data     []byte `json:"data"`
	CheckSum []byte `json:"checksum"`
}

//ResMsg 发送数据结构体
type ResMsg struct {
	Status   string `json:"status"`
	Data     []byte `json:"data"`
	CheckSum []byte `json:"checksum"`
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
	if reqMsgInJSON.UID == "" {
		return nil, "uid_cannot_be_empty", false
	}

	reqMsgDataData := parseResult.Get("data")

	if reqMsgDataData.Exists() {
		var err error
		reqMsgInJSON.Data, err = base64.StdEncoding.DecodeString(reqMsgDataData.String())
		if err != nil {
			return nil, "data_base64_format_error", false
		}
	} else {
		reqMsgInJSON.Data = nil
	}

	regMsgCheckSumData := gjson.GetBytes(reqMsgOrigin, "checksum").String()
	if len(regMsgCheckSumData) == 0 {
		return nil, "checksum_cannot_be_empty", false
	}

	var err error
	reqMsgInJSON.CheckSum, err = base64.StdEncoding.DecodeString(regMsgCheckSumData)
	if err != nil {
		return nil, "data_base64_format_error", false
	}

	md5Ctx := md5.New()
	//fmt.Println(reqMsgInJSON.UID)
	md5Ctx.Write([]byte(reqMsgInJSON.UID))
	//fmt.Println(string(reqMsgInJSON.Data))
	md5Ctx.Write([]byte(reqMsgInJSON.Data))

	//fmt.Println([]byte(md5Ctx.Sum(nil)))
	//fmt.Println([]byte(reqMsgInJSON.CheckSum))

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
		//replyMsg.Data = make([]byte, 256)
		//base64.StdEncoding.Encode(replyMsg.Data, data)
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(replyMsg.Status))
	md5Ctx.Write(replyMsg.Data)
	replyMsg.CheckSum = md5Ctx.Sum(nil)

	jsonMsg, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jsonMsg))
	return jsonMsg, nil
}

func CombineSendMsg(UID string, dataInJSON interface{}) ([]byte, error) {
	var sendMsg ReqMsg
	var err error

	sendMsg.UID = UID

	if dataInJSON == nil {
		sendMsg.Data = nil
	} else {
		data, err := json.Marshal(dataInJSON)
		if err != nil {
			return nil, err
		}

		sendMsg.Data = data
		//	sendMsg.Data = make([]byte, 256)
		//	base64.StdEncoding.Encode(sendMsg.Data, data)
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(sendMsg.UID))
	md5Ctx.Write(sendMsg.Data)
	sendMsg.CheckSum = md5Ctx.Sum(nil)

	jsonMsg, err := json.Marshal(sendMsg)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jsonMsg))
	return jsonMsg, nil
}

func PwdRemoveSalr(pwdWithSalt []byte) (string, error) {
	if pwdWithSalt == nil {
		return "", errors.New("pwd_cannot_be_empty")
	}

	if len(pwdWithSalt) <= 7+7 {
		return "", errors.New("pwd_too_short")
	}

	//去盐
	pwdWithoutSalt := pwdWithSalt[7:]
	//fmt.Println(string(pwdWithoutSalt))
	pwdWithoutSalt[2] -= 2
	pwdWithoutSalt[3] -= 3
	pwdWithoutSalt[5] -= 7
	pwdWithoutSalt[6] -= 11

	//fmt.Println(string(pwdWithoutSalt))
	pwdInByteDecode := make([]byte, 64)

	_, err := base64.StdEncoding.Decode(pwdInByteDecode, pwdWithoutSalt)
	if err != nil {
		return "", errors.New("pwd_format_error")
	}

	zeroIndex := bytes.IndexByte(pwdInByteDecode, 0)
	newPwdInString := string(pwdInByteDecode[0:zeroIndex])
	return newPwdInString, nil
}
