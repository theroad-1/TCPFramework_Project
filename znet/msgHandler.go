package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
	"strconv"
)

/*
 消息处理模块的实现
*/
type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

//初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//调度/执行对应的router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is NOT FOUND! please register")
		return
	}

	//2 根据msgID 调度对应的router业务
	handler.PostHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID =", msgID, "success!")
}

//启动一个worker工作池(开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池)
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize分别开启worker 每个worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 当前的worker对应的channel消息队列 开辟空间 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker ID = ", workerID, "is started ...")
	//不断阻塞等待对应消息队列的信息
	for {
		select {
		case request := <-taskQueue: //从消息队列获取到一个请求后
			mh.DoMsgHandler(request) //进行处理

		}
	}
}

//将消息交给taskQueue，由worker处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 将消息平均分配给不同的worker
	//根据客户端建立的connID来分配,cid是从0开始自增的。
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add connID=", request.GetConnection().GetConnID(),
		"request MsgID =", request.GetMsgID(), "to workerID", workerID)
	//2 将消息发送给对应worker的taskQueue
	mh.TaskQueue[workerID] <- request
}
