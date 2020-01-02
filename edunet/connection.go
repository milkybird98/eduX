package edunet

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer eduiface.IServer
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler eduiface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte

	//链接属性
	session map[string]interface{}
	//保护链接属性修改的锁
	sessionLock sync.RWMutex
}

//创建连接的方法
func NewConntion(server eduiface.IServer, conn *net.TCPConn, sessionID uint32, msgHandler eduiface.IMsgHandle) *Connection {
	//初始化Conn属性
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       sessionID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		session:      make(map[string]interface{}),
	}

	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func NewFileConntion(server eduiface.IServer, conn *net.TCPConn, sessionID uint32) *Connection {
	//初始化Conn属性
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       sessionID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		session:      make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c)
	return c
}

/*
	初始化Session,设定默认值
*/
func initSession(session map[string]interface{}) {
	session["isLogined"] = false
}

/*
	文件传输Goroutine,和客户端进行文件传输
*/
func (c *Connection) StartTransmiter() {
	fmt.Println("[File Transmiter Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn File Transmiter exit!]")
	defer c.Stop()

	fmt.Println("[CONNECT] start file transmite operation")

	serectSlice := make([]byte, 24)
	if _, err := io.ReadFull(c.GetTCPConnection(), serectSlice); err != nil {
		fmt.Println("[CONNECT][ERROR] read serect error ", err)
		return
	}

	fileTag, err := utils.GetFileTranCache(string(serectSlice))
	if err != nil {
		fmt.Println("[CONNECT][ERROR] file not in transmit list ", serectSlice, " client IP addr ", c.RemoteAddr())
		return
	}

	/*
		clientIP := c.RemoteAddr()
		if clientIP != fileTag.ClientAddress {
			fmt.Println("[CONNECT][ERROR] request ip not match, want ", fileTag, " but ", clientIP)
			return
		}
	*/

	data := "ready"
	if _, err := c.Conn.Write([]byte(data)); err != nil {
		fmt.Println("[CONNECT][ERROR] Send Data error: ", err)
		return
	}

	fmt.Println("[CONNECT] file transmite operation ready")

	workPath, err := os.Getwd()
	if err != nil {
		fmt.Println("[CONNECT][ERROR] get work directory fail: ", err)
		return
	}

	if fileTag.ServerToC {

		filePath := workPath + "/file/" + string(fileTag.ID)

		if res, err := utils.PathExists(filePath); res != true {
			fmt.Println("[CONNECT][WARNING] Wanted file not exist,error: ", err)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("[CONNECT][ERROR] Open file error: ", err)
			return
		}
		defer file.Close()

		size, err := io.CopyN(c.GetTCPConnection(), file, fileTag.Size)
		if err != nil {
			fmt.Println("[CONNECT][ERROR] File transmite error: ", err)
			return
		}

		if size != fileTag.Size {
			fmt.Println("[CONNECT][WARNING] File size not match, want: ", fileTag.Size, " fact: ", size)
		}

		fmt.Println("[CONNECT] file transmite operation success!")
		return

	}

	if fileTag.ClientToS {
		filePath := workPath + "/file/" + string(fileTag.ID)

		_, err := os.Stat(filePath)
		if res, _ := utils.PathExists(filePath); res == true {
			fmt.Println("[CONNECT][WARNING] Same serect file already exist: ", string(fileTag.ID))
			return
		}

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("[CONNECT][ERROR] Create file error:, ", err)
			return
		}
		defer file.Close()

		fmt.Println(fileTag.Size)

		size, err := io.CopyN(file, c.GetTCPConnection(), fileTag.Size)
		if err != nil {
			fmt.Println("[CONNECT][WARNING] File transmite error: ", err, ", start removing file...")
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println("[CONNECT][ERROR] Remove file error: ", err)
			} else {
				fmt.Println("[CONNECT] Remove file suceess")
			}
			return
		}

		if size != fileTag.Size {
			fmt.Println("[CONNECT][WARNING] File size not match, want: ", fileTag.Size, " fact: ", size, ", start removing file...")
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println("[CONNECT][ERROR] Remove file error: ", err)
			} else {
				fmt.Println("[CONNECT] Remove file suceess")
			}
			return
		}

		var newFile edumodel.File
		newFile.ClassName = fileTag.ClassName
		newFile.FileName = fileTag.FileName
		newFile.ID, err = primitive.ObjectIDFromHex(fileTag.ID)
		newFile.FileTag = fileTag.FileTags
		if err != nil {
			fmt.Println("[CONNECT][WARNING] File ID Decode error: ", err, ", start removing file...")
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println("[CONNECT][ERROR] Remove file error: ", err)
			} else {
				fmt.Println("[CONNECT] Remove file suceess")
			}
			return
		}
		newFile.Size = uint64(fileTag.Size)
		newFile.UpdateTime = fileTag.UpdateTime
		newFile.UpdaterUID = fileTag.UpdaterUID

		ok := edumodel.AddFile(&newFile)
		if !ok {
			return
		}

		fmt.Println("[CONNECT] file transmite operation success!")
		return
	}
}

/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[CONNECT] Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

/*
	读消息Goroutine，用于从客户端中读取数据
*/
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()

	for {
		// 创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("[CONNECT] read msg head error ", err)
			break
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("[CONNECT] unpack error ", err)
			break
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("[CONNECT] read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) StartFileTransmit() {
	//开始准备和客户端进行文件传输
	go c.StartTransmiter()
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//初始化session
	initSession(c.session)
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	fmt.Println("[CONNECT] Conn Stop()...ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer
	c.ExitBuffChan <- true

	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("[CONNECT] Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("[CONNECT] Pack error msg id = ", msgId)
		return errors.New("[CONNECT] Pack error msg ")
	}

	//写回客户端
	c.msgChan <- msg

	return nil
}

//设置链接属性
func (c *Connection) SetSession(key string, value interface{}) {
	c.sessionLock.Lock()
	defer c.sessionLock.Unlock()

	c.session[key] = value
}

//获取链接属性
func (c *Connection) GetSession(key string) (interface{}, error) {
	c.sessionLock.RLock()
	defer c.sessionLock.RUnlock()

	if value, ok := c.session[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("[CONNECT] no session found")
	}
}

//移除链接属性
func (c *Connection) RemoveSession(key string) {
	c.sessionLock.Lock()
	defer c.sessionLock.Unlock()

	delete(c.session, key)
}
