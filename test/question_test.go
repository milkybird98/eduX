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

	{
		fmt.Println("[TEST]test login")

		db := edunet.NewDataPack()

		var Data edurouter.LoginData
		Data.Pwd = "123"

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
		Data.Pwd = "123"

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

	{
		fmt.Println("[TEST]test question answer")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionAnswerData
		Data.QuestionID = "5e05a0f4afea43c4cf4f674a"
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
		fmt.Println("[TEST]test question get by class")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionGetByClassNameData
		Data.ClassName = "ts1001"

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
		if replystatus == "success" {
			fmt.Println("[TEST]question get by class pass")
		}
	}

	{
		fmt.Println("[TEST]test question get by sender UID")

		db := edunet.NewDataPack()

		var Data edurouter.QuestionGetBySenderUIDData
		Data.SenderUID = "U1001"
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

}
