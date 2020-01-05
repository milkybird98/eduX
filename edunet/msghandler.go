package edunet

import (
	"eduX/eduiface"
	"eduX/utils"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]eduiface.IRouter //存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32                      //业务工作Worker池的数量
	TaskQueue      []chan eduiface.IRequest    //Worker负责取任务的消息队列
	WorkerIndex    uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]eduiface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue:   make([]chan eduiface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerIndex: 0,
	}
}

func (mh *MsgHandle) NextIndex() uint32 {
	if mh.WorkerIndex >= (mh.WorkerPoolSize - 1) {
		mh.WorkerIndex = 0
	} else {
		mh.WorkerIndex++
	}
	return mh.WorkerIndex
}

//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request eduiface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := mh.NextIndex()
	fmt.Println("[MSGHAND] Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[request.GetConnection().GetConnID()%mh.WorkerPoolSize] <- request
}

//马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request eduiface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("[MSGHAND][WARNING] api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router eduiface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("[MSGHAND][ERROR] repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("[MSGHAND] Add api msgId = ", msgId)
}

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan eduiface.IRequest) {
	fmt.Println("[MSGHAND][IMPORTANT] Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	fmt.Println("[MSGHAND][IMPORTANT] Worker Poll Starting...")
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan eduiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
	fmt.Println("[MSGHAND][IMPORTANT] Worker Poll Started, Workerpoll Size: ", mh.WorkerPoolSize)
}
