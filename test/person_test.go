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
	edumodel.ConnectMongo()
	edumodel.ConnectDatabase(nil)

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

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)
	}

	{
		db := edunet.NewDataPack()

		var student edurouter.PersonAddData
		student.Place = "student"
		student.Name = "测试姓名1"
		student.UID = "U1001"

		msgData, _ := edurouter.CombineSendMsg("M1001", student)

		msg := edunet.NewMsgPackage(1, msgData)

		data, _ := db.Pack(msg)

		_, err := conn.Write(data)
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

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status").String()
		fmt.Printf("reply status %s\n", replystatus)

	}

	{
		fmt.Println("test person info get by class router")

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

		buf := make([]byte, 512)
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
		replystatus := replyData.Get("status")
		fmt.Println("reply status: ", replystatus.String())
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
	}
	{

		fmt.Println("test person info get by uid router")

		db := edunet.NewDataPack()

		var person edurouter.PersonInfoGetData

		person.UID = "U1001"

		msgData, _ := edurouter.CombineSendMsg("M1001", person)

		msg := edunet.NewMsgPackage(3, msgData)

		data, _ := db.Pack(msg)

		_, err := conn.Write(data)
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
			fmt.Println("unpack error ", err)
			return
		}

		//根据 dataLen 读取 data，放在msg.Data中

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status")
		fmt.Println("reply status: ", replystatus.String())
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))

	}

	{
		fmt.Println("test person info put by uid router")

		db := edunet.NewDataPack()

		var person edurouter.PersonInfoPutData

		person.UID = "U1001"
		person.Name = "tt"
		person.Gender = 2

		msgData, _ := edurouter.CombineSendMsg("M1001", person)

		msg := edunet.NewMsgPackage(4, msgData)

		data, _ := db.Pack(msg)

		_, err := conn.Write(data)
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
			fmt.Println("unpack error ", err)
			return
		}

		//根据 dataLen 读取 data，放在msg.Data中

		replyMsg.SetData(buf[8:cnt])

		fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

		replyData := gjson.ParseBytes(replyMsg.GetData())
		replystatus := replyData.Get("status")
		fmt.Println("reply status: ", replystatus.String())
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
	}

}
