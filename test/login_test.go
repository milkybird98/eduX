package test

import (
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/edurouter"
	"encoding/base64"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

func TestServerLogin(t *testing.T) {
	//创建一个server句柄
	edumodel.ConnectMongo()

	edumodel.ConnectDatabase(nil)

	s := edunet.NewServer()

	s.AddRouter(0, &edurouter.LoginRouter{})
	s.AddRouter(1, &edurouter.PwdSetQuestionRouter{})
	s.AddRouter(2, &edurouter.PwdGetQuestionRouter{})
	s.AddRouter(3, &edurouter.PwdResetRouter{})
	s.AddRouter(4, &edurouter.PwdForgetRouter{})

	//	客户端测试
	go ClientTestLI(t)

	//2 开启服务
	s.Serve()
}

func ClientTestLI(t *testing.T) {

	fmt.Println("Client Test ... start")

	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	connTea, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	var serect string
	{
		fmt.Println("[test] login test")
		db := edunet.NewDataPack()

		var Data edurouter.LoginData
		Data.Pwd = base64.StdEncoding.EncodeToString([]byte("MTExMTEx"))

		PwdInByte := []byte(Data.Pwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		Data.Pwd = string(PwdInByte)
		Data.Pwd = ("1234567") + Data.Pwd

		msgData, _ := edurouter.CombineSendMsg("A100110001", Data)

		msg := edunet.NewMsgPackage(0, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] login pass")
		}
	}

	{
		fmt.Println("[test] login test")
		db := edunet.NewDataPack()

		var Data edurouter.LoginData
		Data.Pwd = base64.StdEncoding.EncodeToString([]byte("12312312"))

		PwdInByte := []byte(Data.Pwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		Data.Pwd = string(PwdInByte)
		Data.Pwd = ("1234567") + Data.Pwd

		msgData, _ := edurouter.CombineSendMsg("T1001", Data)

		msg := edunet.NewMsgPackage(0, msgData)

		data, _ := db.Pack(msg)

		_, err = connTea.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := connTea.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] login pass")
		}
	}

	{
		fmt.Println("[test] set question test")

		db := edunet.NewDataPack()

		var Data edurouter.PwdSetQuestionData
		Data.Pwd = base64.StdEncoding.EncodeToString([]byte("12312312"))

		PwdInByte := []byte(Data.Pwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		Data.Pwd = string(PwdInByte)
		Data.Pwd = ("1234567") + Data.Pwd

		Data.QuestionA = "q1"
		Data.QuestionB = "qsdsf"
		Data.QuestionC = "safawq"
		Data.AnswerA = "123"
		Data.AnswerB = "123"
		Data.AnswerC = "123"

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(1, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] set question pass")
		}
	}

	{
		fmt.Println("[test] get question test")

		db := edunet.NewDataPack()

		msgData, _ := edurouter.CombineSendMsg("U1003", nil)

		msg := edunet.NewMsgPackage(2, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] get question pass")
		}
	}

	{
		fmt.Println("[test] reset password test -- by origin password")

		db := edunet.NewDataPack()

		var Data edurouter.PwdResetData
		Data.UID = "U1003"

		Data.OriginPwd = base64.StdEncoding.EncodeToString([]byte("12312312"))

		PwdInByte := []byte(Data.OriginPwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		Data.OriginPwd = string(PwdInByte)
		Data.OriginPwd = ("1234567") + Data.OriginPwd

		Data.NewPwd = Data.OriginPwd

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(3, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] reset password pass")
		}
	}

	{
		fmt.Println("[test] reset password test -- by teacher")

		db := edunet.NewDataPack()

		var Data edurouter.PwdResetData
		Data.UID = "U1003"

		Data.OriginPwd = base64.StdEncoding.EncodeToString([]byte("12312312"))

		PwdInByte := []byte(Data.OriginPwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		Data.OriginPwd = string(PwdInByte)
		Data.OriginPwd = ("1234567") + Data.OriginPwd

		msgData, _ := edurouter.CombineSendMsg("T1001", Data)

		msg := edunet.NewMsgPackage(3, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] reset password pass")
		}
	}

	{
		fmt.Println("[test] reset password test -- answer question")

		db := edunet.NewDataPack()

		var Data edurouter.PwdForgetData
		Data.AnswerA = "123"
		Data.AnswerB = "123"
		Data.AnswerC = "123"

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(4, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] reset password pass")
		}

		serect = gjson.GetBytes(data, "serect").String()
	}

	{
		fmt.Println("[test] reset password test -- by serect")

		db := edunet.NewDataPack()

		var Data edurouter.PwdResetData
		Data.UID = "U1003"

		pwd := base64.StdEncoding.EncodeToString([]byte("12312312"))

		PwdInByte := []byte(pwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		Data.NewPwd = ("1234567") + string(PwdInByte)

		Data.Serect = serect

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(3, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] reset password pass")
		}
	}
}
