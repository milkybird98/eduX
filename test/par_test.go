package test

import (
	"eduX/edunet"
	"eduX/edurouter"
	"encoding/base64"
	"fmt"
	"net"
	"testing"

	"github.com/tidwall/gjson"
)

func TestParallel(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go login()
	}

	select {}
}

func login() {
	conn, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for i := 0; i < 3; i++ {
		{
			//fmt.Println("[TEST]test login")

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

			//fmt.Printf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

			replyData := gjson.ParseBytes(replyMsg.GetData())
			replystatus := replyData.Get("status").String()
			//fmt.Printf("reply status %s\n", replystatus)
			if replystatus == "success" {
				fmt.Println("[TEST]login pass")
			}
		}
	}
}
