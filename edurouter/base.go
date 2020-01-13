package edurouter

import (
	"bytes"
	"crypto/md5"
	"eduX/eduiface"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/tidwall/gjson"
)

/*
	用于处理通用的路由操作
*/

var checksumFlag bool

// ReqMsg 接收数据结构体
type ReqMsg struct {
	UID      string `json:"uid"`
	Data     []byte `json:"data"`
	CheckSum string `json:"checksum"`
}

// ResMsg 发送数据结构体
type ResMsg struct {
	Status   string `json:"status"`
	Data     []byte `json:"data"`
	CheckSum string `json:"checksum"`
}

//CheckMsgFormat 检查接收数据格式是否正确,并将数据解码
func CheckMsgFormat(request eduiface.IRequest) (*ReqMsg, string, bool) {
	// 获取原始的接收数据
	var reqMsgInJSON ReqMsg
	reqMsgOrigin := request.GetData()
	checksumFlag = false

	// 验证原始数据是否是合法json格式
	if !gjson.Valid(string(reqMsgOrigin)) {
		return nil, "json_format_error", false
	}

	// 解析原始数据
	parseResult := gjson.ParseBytes(reqMsgOrigin)

	// 试图从原始数据中获取人员uid
	reqMsgInJSON.UID = parseResult.Get("uid").String()
	// 如果uid不存在则报错返回
	if reqMsgInJSON.UID == "" {
		return nil, "req_uid_cannot_be_empty", false
	}

	// 从原始数据中获取Data段
	reqMsgDataData := parseResult.Get("data")

	// 如果Data段存在
	if reqMsgDataData.Exists() {
		var err error
		// 对Data做base64解码
		reqMsgInJSON.Data, err = base64.StdEncoding.DecodeString(reqMsgDataData.String())
		if err != nil {
			return nil, "data_base64_format_error", false
		}
	} else {
		reqMsgInJSON.Data = nil
	}

	// 试图从原始数据中获取校验和
	regMsgCheckSumData := gjson.GetBytes(reqMsgOrigin, "checksum").String()
	// 校验和不存在则报错返回
	if len(regMsgCheckSumData) == 0 {
		return nil, "checksum_cannot_be_empty", false
	}
	reqMsgInJSON.CheckSum = regMsgCheckSumData

	// 根据获取数据计算校验和
	md5Ctx := md5.New()

	md5Ctx.Write([]byte(reqMsgInJSON.UID))
	md5Ctx.Write([]byte(reqMsgInJSON.Data))

	// 如果校验和不一致则报错返回
	if string(reqMsgInJSON.CheckSum) != hex.EncodeToString(md5Ctx.Sum(nil)) {
		return nil, "check_sum_error", false
	}

	// 原始数据处理完毕
	return &reqMsgInJSON, "", true
}

//CheckConnectionLogin 检查当前用户是否已登陆
func CheckConnectionLogin(request eduiface.IRequest, UID string) (string, bool) {
	// 试图从session中获取UID
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	// 如果获取失败,则报错返回
	if err != nil {
		return "session_uid_not_found", false
	}

	// 如果发送请求UID于session中UID不一致,则报错返回
	if sessionUID != UID {
		return "uid_not_match", false
	}

	// 试图中session中获取登陆状态
	value, err := c.GetSession("isLogined")
	// 如果登陆状态不存在,则报错返回
	if err != nil {
		return "session_login_status_not_found", false
	}

	// 如果为登陆,则返回未登录
	if value == false {
		return "not_login", false
	}

	// 该人员已登录
	return "", true
}

//CombineReplyMsg 拼接状态和返回数据Data段,计算MD5校验和并且序列化数据
func CombineReplyMsg(status string, dataInJSON interface{}) ([]byte, error) {
	var replyMsg ResMsg
	var err error
	// 赋值返回状态
	if status == "" {
		panic(err.Error())
	}

	replyMsg.Status = status

	// 如果返回数据有Data段,则试图序列化Data段
	if dataInJSON == nil {
		replyMsg.Data = nil
	} else {
		data, err := json.Marshal(dataInJSON)
		// 序列化失败报错返回
		if err != nil {
			return nil, err
		}
		// 赋值序列化后的Data段
		replyMsg.Data = data
	}

	// 计算MD5校验和并赋值
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(replyMsg.Status))
	md5Ctx.Write(replyMsg.Data)
	replyMsg.CheckSum = hex.EncodeToString(md5Ctx.Sum(nil))

	// 将返回数据序列化
	jsonMsg, err := json.Marshal(replyMsg)
	if err != nil {
		return nil, err
	}

	return jsonMsg, nil
}

//CombineSendMsg 拼接状态和发送数据Data段,计算MD5校验和并且序列化数据
func CombineSendMsg(UID string, dataInJSON interface{}) ([]byte, error) {
	var sendMsg ReqMsg
	var err error
	// 赋值发送UID
	sendMsg.UID = UID

	// 如果发送数据有Data段,则试图序列化Data段
	if dataInJSON == nil {
		sendMsg.Data = nil
	} else {
		data, err := json.Marshal(dataInJSON)
		// 序列化失败报错返回
		if err != nil {
			return nil, err
		}
		// 赋值序列化后的Data段
		sendMsg.Data = data
	}

	// 计算MD5校验和并赋值
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(sendMsg.UID))
	md5Ctx.Write(sendMsg.Data)
	sendMsg.CheckSum = hex.EncodeToString(md5Ctx.Sum(nil))

	// 将发送数据序列化
	jsonMsg, err := json.Marshal(sendMsg)
	if err != nil {
		return nil, err
	}

	return jsonMsg, nil
}

// PwdRemoveSalr 将收到的加盐后的密码数据去盐
func PwdRemoveSalr(pwdWithSalt []byte) (string, error) {
	// 进行基本的数据校验
	if pwdWithSalt == nil {
		return "", errors.New("pwd_cannot_be_empty")
	}

	if len(pwdWithSalt) <= 7+7 {
		return "", errors.New("pwd_too_short")
	}

	//去盐
	pwdWithoutSalt := pwdWithSalt[7:]
	pwdWithoutSalt[2] -= 2
	pwdWithoutSalt[3] -= 3
	pwdWithoutSalt[5] -= 7
	pwdWithoutSalt[6] -= 11

	// 去盐结果做base64解码
	pwdInByteDecode := make([]byte, 64)
	_, err := base64.StdEncoding.Decode(pwdInByteDecode, pwdWithoutSalt)
	// 如果解码失败说明格式错误,报错并返回
	if err != nil {
		return "", errors.New("pwd_format_error")
	}

	// 去除尾部的零
	zeroIndex := bytes.IndexByte(pwdInByteDecode, 0)
	newPwdInString := string(pwdInByteDecode[0:zeroIndex])

	// 返回解码后的密码结果
	return newPwdInString, nil
}

// GetSkipAndLimit 用于获取请求数据中的skip和limit项目,在不存在时使用默认值
func GetSkipAndLimit(msgData []byte) (int64, int64) {
	var Skip int64
	// 从原始数据中获取skip数据
	skipData := gjson.GetBytes(msgData, "skip")
	// 如果不存在则使用默认值0
	if skipData.Exists() && skipData.Int() >= 0 {
		Skip = skipData.Int()
	} else {
		Skip = 0
	}

	var Limit int64
	// 从原始数据中获取limit数据
	limitData := gjson.GetBytes(msgData, "limit")
	// 如果不存在则使用默认值10
	if limitData.Exists() && limitData.Int() > 0 {
		Limit = limitData.Int()
	} else {
		Limit = 10
	}

	return Skip, Limit
}

// GetSessionPlace 用于从当前连接的session中获取place数据,并转换为字符串
func GetSessionPlace(c eduiface.IConnection) (string, error) {
	// 试图从session中获取place
	sessionPlace, err := c.GetSession("place")
	// 如果不存在则报错返回
	if err != nil {
		return "", errors.New("session_place_not_found")
	}

	// 试图将其转换为字符串
	placeString, ok := sessionPlace.(string)
	// 如果转换失败则报错返回
	if ok != true {
		return "", errors.New("session_place_data_error")
	}

	if placeString != "manager" && placeString != "teacher" && placeString != "student" {
		return "", errors.New("session_place_format_error")
	}

	// 返回从session中得到的身份数据
	return placeString, nil
}
