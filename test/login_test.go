package edutest

import (
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/edurouter"
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

	//	客户端测试
	go ClientTestLI(t)

	//2 开启服务
	s.Serve()
}

func ClientTestLI(t *testing.T) {

	t.Log("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:23333")
	if err != nil {
		t.Log("client start err, exit!")
		return
	}

	db := edunet.NewDataPack()

	var loginData edurouter.LoginData
	loginData.Pwd = "123"

	msgData, _ := edurouter.CombineSendMsg("M1001", loginData)

	msg := edunet.NewMsgPackage(0, msgData)

	data, _ := db.Pack(msg)

	_, err = conn.Write(data)
	if err != nil {
		t.Log("write error err ", err)
		return
	}

	buf := make([]byte, 512)
	cnt, err := conn.Read(buf)
	if err != nil {
		t.Log("read buf error ")
		return
	}

	replyMsg, err := db.Unpack(buf)
	if err != nil {
		t.Log(err)
		return
	}

	t.Logf("server call back msgID = %d, msgLength = %d, originLength = %d\n", replyMsg.GetMsgId(), replyMsg.GetDataLen(), cnt)

	replyData := gjson.ParseBytes(replyMsg.GetData())
	replystatus := replyData.Get("status").String()
	t.Logf("reply status %s\n", replystatus)
}
