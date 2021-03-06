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

func TestServerQuestionOpertaion(t *testing.T) {
	edumodel.ConnectMongo()
	edumodel.ConnectDatabase(nil)

	//创建一个server句柄
	s := edunet.NewServer()

	s.AddRouter(0, &edurouter.LoginRouter{})
	s.AddRouter(1, &edurouter.QuestionAddRouter{})
	s.AddRouter(2, &edurouter.QuestionAnswerRouter{})
	s.AddRouter(3, &edurouter.QuestionDeleteRouter{})
	s.AddRouter(4, &edurouter.QuestionGetByClassNameRouter{})
	s.AddRouter(5, &edurouter.QuestionGetBySenderUIDRouter{})
	s.AddRouter(6, &edurouter.QuestionCountRouter{})

	go ClientTestQO(t)

	s.Serve()
}

func ClientTestQO(t *testing.T) {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	connStu, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	connTea, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	connMgr, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	{
		fmt.Println("[TEST]test login")

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

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(0, msgData)

		data, _ := db.Pack(msg)

		_, err = connStu.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := connStu.Read(buf)
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
		if replystatus == "success" {
			fmt.Println("[TEST]login pass")
		}
	}
	{
		fmt.Println("[TEST]test login")

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
		if replystatus == "success" {
			fmt.Println("[TEST]login pass")
		}
	}

	{
		fmt.Println("[TEST]test login")

		db := edunet.NewDataPack()

		var loginData edurouter.LoginData
		loginData.Pwd = base64.StdEncoding.EncodeToString([]byte("12312312"))

		PwdInByte := []byte(loginData.Pwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		loginData.Pwd = string(PwdInByte)
		loginData.Pwd = ("1234567") + loginData.Pwd

		msgData, _ := edurouter.CombineSendMsg("M1001", loginData)

		msg := edunet.NewMsgPackage(0, msgData)

		data, _ := db.Pack(msg)

		_, err = connMgr.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 4096)
		cnt, err := connMgr.Read(buf)
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
		if replystatus == "success" {
			fmt.Println("[TEST]login pass")
		}
	}

	{
		fmt.Println("[TEST]test question add")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionAddData
		Data.Title = "test"
		Data.Text = "testsssadfawe"

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(1, msgData)

		data, _ := db.Pack(msg)

		_, err = connStu.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := connStu.Read(buf)
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
		if replystatus == "success" {
			fmt.Println("[TEST]question add pass")
		}
	}

	var notansid string

	{
		fmt.Println("[TEST]test question get by class")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionGetByClassNameData
		Data.ClassName = "ts1001"
		Data.DeferSolved = true
		Data.IsSolved = false

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(4, msgData)

		data, _ := db.Pack(msg)

		_, err = connStu.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 40960)
		cnt, err := connStu.Read(buf)
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
		notansid = gjson.GetBytes(data, "questions").Array()[0].Get("ID").String()
		if replystatus == "success" {
			fmt.Println("[TEST]question get by class pass")
		}
	}

	{
		fmt.Println("[TEST]test question answer")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionAnswerData
		Data.QuestionID = notansid
		Data.Answer = "答 案ansasdfass)*(as23&%&#^w   e r   "

		msgData, _ := edurouter.CombineSendMsg("T1001", Data)

		msg := edunet.NewMsgPackage(2, msgData)

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
		if replystatus == "success" {
			fmt.Println("[TEST]question answer pass")
		}
	}

	{
		fmt.Println("[TEST]test question get by sender UID")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionGetBySenderUIDData
		Data.SenderUID = "U1003"
		Data.DeferSolved = false
		Data.IsSolved = true
		Data.Limit = 3

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(5, msgData)

		data, _ := db.Pack(msg)

		_, err = connStu.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 40960)
		cnt, err := connStu.Read(buf)
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
			fmt.Println("[TEST]question get by sender UID pass")
		}
	}

	{
		fmt.Println("[TEST]test question count")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionCountData
		Data.Date = time.Now()
		Data.ClassName = "ts1001"

		msgData, _ := edurouter.CombineSendMsg("M1001", Data)

		msg := edunet.NewMsgPackage(6, msgData)

		data, _ := db.Pack(msg)

		_, err = connMgr.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 40960)
		cnt, err := connMgr.Read(buf)
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
			fmt.Println("[TEST]question count pass")
		}
	}

	{
		fmt.Println("[TEST]test question count")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionCountData

		msgData, _ := edurouter.CombineSendMsg("M1001", Data)

		msg := edunet.NewMsgPackage(6, msgData)

		data, _ := db.Pack(msg)

		_, err = connMgr.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 40960)
		cnt, err := connMgr.Read(buf)
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
			fmt.Println("[TEST]question count pass")
		}
	}

}
