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

func TestServerFileOpertaion(t *testing.T) {
	edumodel.ConnectMongo()
	edumodel.ConnectDatabase(nil)

	//创建一个server句柄
	s := edunet.NewServer()
	sFile := edunet.NewServer()

	s.AddRouter(0, &edurouter.LoginRouter{})
	s.AddRouter(1, &edurouter.FileAddRouter{})
	s.AddRouter(2, &edurouter.FileDownloadRouter{})
	s.AddRouter(3, &edurouter.FileGetByClassNameRouter{})
	s.AddRouter(4, &edurouter.FileGetBySenderUIDRouter{})
	s.AddRouter(5, &edurouter.FileGetByTagsRouter{})
	s.AddRouter(6, &edurouter.FileCountRouter{})
	s.AddRouter(7, &edurouter.FileDeleteRouter{})

	go ClientTestFO(t)

	go s.Serve()
	go sFile.ServeFile()

	for a := 2; a < 10; a++ {
		a = 1
	}

}

func ClientTestFO(t *testing.T) {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:23333")
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

		var loginData edurouter.LoginData
		loginData.Pwd = base64.StdEncoding.EncodeToString([]byte("12312312"))

		PwdInByte := []byte(loginData.Pwd)
		PwdInByte[2] += 2
		PwdInByte[3] += 3
		PwdInByte[5] += 7
		PwdInByte[6] += 11

		loginData.Pwd = string(PwdInByte)
		loginData.Pwd = ("1234567") + loginData.Pwd

		msgData, _ := edurouter.CombineSendMsg("U1003", loginData)

		msg := edunet.NewMsgPackage(0, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 4096)
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
		if replystatus == "success" {
			fmt.Println("[TEST]login pass")
		}
	}

	{
		fmt.Println("[TEST]test file add")

		db := edunet.NewDataPack()

		var Data edurouter.FileAddData
		Data.ClassName = "ts1001"
		Data.FileName = "a new fes"
		Data.FileTag = []string{"blank"}
		Data.Size = 10

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(1, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 4096)
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
			fmt.Println("[TEST]add file pass")
		}
	}

	{
		fmt.Println("[TEST]test file download")

		db := edunet.NewDataPack()

		var Data edurouter.FileDownloadData
		Data.UUID = "5e12452c6670537b52c2cda9"

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(2, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 4096)
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
			fmt.Println("[TEST]get file download serect")
		}
	}

	{
		fmt.Println("[TEST]test file list get by classname")

		db := edunet.NewDataPack()

		var Data edurouter.FileGetByClassNameData
		Data.ClassName = "ts1001"

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(3, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 4096)
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
			fmt.Println("[TEST]get file list pass")
		}
	}

	{
		fmt.Println("[TEST]test file list get by senduid")

		db := edunet.NewDataPack()

		var Data edurouter.FileGetBySenderUIDData
		Data.Sender = "U1003"

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(4, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 4096)
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
			fmt.Println("[TEST]get file list pass")
		}

	}

	{
		fmt.Println("[TEST]test file list get by tags")

		db := edunet.NewDataPack()

		var Data edurouter.FileGetByTagsData
		Data.Tags = []string{"blank"}

		msgData, _ := edurouter.CombineSendMsg("U1003", Data)

		msg := edunet.NewMsgPackage(5, msgData)

		data, _ := db.Pack(msg)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 4096)
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
			fmt.Println("[TEST]get file list pass")
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
		fmt.Println("[TEST]test file count")

		db := edunet.NewDataPack()

		var Data edurouter.FileCountData

		Data.ClassName = "ts1001"
		Data.Date = time.Now()

		msgData, _ := edurouter.CombineSendMsg("M1001", Data)

		msg := edunet.NewMsgPackage(6, msgData)

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
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[TEST]get file list pass")
		}
	}

	{
		fmt.Println("[TEST]test file delete")

		db := edunet.NewDataPack()

		var Data edurouter.FileDeleteData
		Data.FileID = "5e0e1c935197e45e740ef031"

		msgData, _ := edurouter.CombineSendMsg("M1001", Data)

		msg := edunet.NewMsgPackage(7, msgData)

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
		data, _ = base64.StdEncoding.DecodeString(replyData.Get("data").String())
		fmt.Println("reply data: ", string(data))
		if replystatus == "success" {
			fmt.Println("[TEST]get file list pass")
		}
	}
}
