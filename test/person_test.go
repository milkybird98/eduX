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

func TestServerPersonOperation(t *testing.T) {
	if !edumodel.ConnectMongo() {
		t.FailNow()
	}
	if !edumodel.ConnectDatabase(nil) {
		t.FailNow()
	}

	//创建一个server句柄
	s := edunet.NewServer()

	s.AddRouter(0, &edurouter.LoginRouter{})
	s.AddRouter(1, &edurouter.PersonAddRouter{})
	s.AddRouter(2, &edurouter.PersonInfoGetByClassRouter{})
	s.AddRouter(3, &edurouter.PersonInfoGetRouter{})
	s.AddRouter(4, &edurouter.PersonInfoPutRouter{})

	//	客户端测试
	go ClientTestSA(t)

	//2 开启服务
	s.Serve()
}

func ClientTestSA(t *testing.T) {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	{
		fmt.Println("[test] login")

		db := edunet.NewDataPack()

		var loginData edurouter.LoginData
		loginData.Pwd = "123"

		msgData, _ := edurouter.CombineSendMsg("M1001", loginData)

		msg := edunet.NewMsgPackage(0, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 2048)
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
		fmt.Println("[test] person info add router")

		db := edunet.NewDataPack()

		var student edurouter.PersonAddData
		student.Place = "student"
		student.Name = "测试姓名1"
		student.UID = "U1005"

		msgData, _ := edurouter.CombineSendMsg("M1001", student)

		msg := edunet.NewMsgPackage(1, msgData)

		data, _ := db.Pack(msg)

		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 2048)
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
			fmt.Println("[test] person add pass")
		}

	}
	{
		fmt.Println("[test] person info put by uid router")

		db := edunet.NewDataPack()

		var person edurouter.PersonInfoPutData

		person.UID = "U1005"
		person.Name = "新名字"
		person.Gender = 2
		person.Birth = ""
		person.Political = "群众"
		person.Contact = "123456789"
		person.IsContactPub = true
		person.Email = "asd@asd.com"
		person.IsEmailPub = false

		msgData, _ := edurouter.CombineSendMsg("M1001", person)

		msg := edunet.NewMsgPackage(4, msgData)

		data, _ := db.Pack(msg)

		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 2048)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println("unpack error ", err)
			return
		}

		//根据 dataLen 读取 data，放在msg.Data中

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Println("reply status: ", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] person info put pass")
		}
	}

	{
		fmt.Println("[test] person info get by class router")

		db := edunet.NewDataPack()

		var person edurouter.PersonInfoGetByClassData

		person.ClassName = "ts1001"

		msgData, _ := edurouter.CombineSendMsg("M1001", person)

		msg := edunet.NewMsgPackage(2, msgData)

		data, _ := db.Pack(msg)

		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 2048)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println("unpack error ", err)
			return
		}

		//根据 dataLen 读取 data，放在msg.Data中

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Println("reply status: ", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] person info get by class pass")
		}
	}
	{
		fmt.Println("[test] person info get by uid router")

		db := edunet.NewDataPack()

		var person edurouter.PersonInfoGetData

		person.UID = "U1005"

		msgData, _ := edurouter.CombineSendMsg("M1001", person)

		msg := edunet.NewMsgPackage(3, msgData)

		data, _ := db.Pack(msg)

		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 2048)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		replyMsg, err := db.Unpack(buf)
		if err != nil {
			fmt.Println("unpack error ", err)
			return
		}

		//根据 dataLen 读取 data，放在msg.Data中

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Println("reply status: ", replystatus)
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[test] person info get by uid pass")
		}
	}

}
