package edutest

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

func TestServerClassOpertaion(t *testing.T) {
	edumodel.ConnectMongo()
	edumodel.ConnectDatabase(nil)

	//创建一个server句柄
	s := edunet.NewServer()

	s.AddRouter(0, &edurouter.LoginRouter{})
	s.AddRouter(1, &edurouter.ClassAddRouter{})
	s.AddRouter(2, &edurouter.ClassDelRouter{})
	s.AddRouter(3, &edurouter.ClassJoinInGetRouter{})
	s.AddRouter(4, &edurouter.ClassListGetRouter{})
	s.AddRouter(5, &edurouter.ClassStudentAddRouter{})
	s.AddRouter(6, &edurouter.ClassStudentDelRouter{})

	go ClientTestCO(t)

	s.Serve()
}

func ClientTestCO(t *testing.T) {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	{
		fmt.Println("[TEST]test login")

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
		fmt.Println("[TEST]test class add")

		db := edunet.NewDataPack()

		var Class edurouter.ClassAddData
		Class.ClassName = "ts1003"
		Class.TeacherUID = ""

		msgData, _ := edurouter.CombineSendMsg("M1001", Class)

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
		if replystatus == "init_teacher_cannot_be_empty" {
			fmt.Println("[TEST]class add pass")
		}
	}
	{
		fmt.Println("[TEST]test class add")

		db := edunet.NewDataPack()

		var Class edurouter.ClassAddData
		Class.ClassName = "ts1003"
		Class.TeacherUID = "T1001"

		msgData, _ := edurouter.CombineSendMsg("M1001", Class)

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
		if replystatus == "same_class_exist" {
			fmt.Println("[TEST]class add pass")
		}
	}
	{
		fmt.Println("[TEST]test class add")

		db := edunet.NewDataPack()

		var Class edurouter.ClassAddData
		Class.ClassName = ""
		Class.TeacherUID = "test"

		msgData, _ := edurouter.CombineSendMsg("M1001", Class)

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
		if replystatus == "classname_cannot_be_empty" {
			fmt.Println("[TEST]class add pass")
		}
	}
	{
		fmt.Println("[TEST]test class add")

		db := edunet.NewDataPack()

		var Class edurouter.ClassAddData
		Class.ClassName = "ts1004"
		Class.TeacherUID = "T1001"

		msgData, _ := edurouter.CombineSendMsg("M1001", Class)

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
		if replystatus == "success" {
			fmt.Println("[TEST]class add pass")
		}
	}
	{
		fmt.Println("[TEST]test class delete")

		db := edunet.NewDataPack()

		var Class edurouter.ClassDelData
		Class.ClassName = "ts1004"

		msgData, _ := edurouter.CombineSendMsg("M1001", Class)

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
			fmt.Println("[TEST]get join list pass")
		}
	}
	{
		fmt.Println("[TEST]test class join list get")

		db := edunet.NewDataPack()

		msgData, _ := edurouter.CombineSendMsg("M1001", nil)

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
			fmt.Println("[TEST]class join list pass")
		}
	}
	{
		fmt.Println("[TEST]test class list get")

		db := edunet.NewDataPack()

		var Data edurouter.ClassListGetData
		Data.Skip = 0
		Data.Limit = 2

		msgData, _ := edurouter.CombineSendMsg("M1001", Data)

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
			fmt.Println("[TEST]class list pass")
		}
	}
	{
		fmt.Println("[TEST]test class student add")

		db := edunet.NewDataPack()

		var Data edurouter.ClassStudentAddData
		Data.ClassName = "ts1001"
		Data.StudentListInUID = []string{"U1001", "U1002", "U1003"}

		msgData, _ := edurouter.CombineSendMsg("M1001", Data)

		msg := edunet.NewMsgPackage(5, msgData)

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
			fmt.Println("[TEST]class student add pass")
		}
	}
	{
		fmt.Println("[TEST]test class student delete")

		db := edunet.NewDataPack()

		var Data edurouter.ClassStudentDelData
		Data.ClassName = "ts1001"
		Data.StudentListInUID = []string{"U1001", "U1002"}

		msgData, _ := edurouter.CombineSendMsg("M1001", Data)

		msg := edunet.NewMsgPackage(6, msgData)

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
			fmt.Println("[TEST]class student del pass")
		}
	}
}
