package GuTcp

import "sync"

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

//Create
// 创建tcp服务,成功返回true,失败返回false并在控制台打印错误信息
//	port 端口号
func (t *GuTcpServer) Create(port int) {
	if t.user == nil {
		t.userId = new(uint64)
		t.user = new(sync.Map)
	}
	var err error
	t.listener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		go t.accept()
		select {} //阻塞,主进程退出他才会退出
	} else {
		fmt.Println(err.Error())
	}

}

//SetHeartBoatTime 设置心跳时间
//  time 单位秒
func (t *GuTcpServer) SetHeartBoatTime(i int) {
	t.heartBoatTime = time.Duration(i) * time.Second
}

//BindingEvent 绑定事件
//	OnRecvHandlerFunc 数据接收事件
func (t *GuTcpServer) BindingEvent(OnRecvHandlerFunc OnServerRecvHandlerFunc, OnServerClose OnServerClose) {
	t.onRecvHandlerFunc = OnRecvHandlerFunc
	t.onServerClose = OnServerClose
}

//heartBoating 心跳事件
func (t *GuTcpServer) heartBoating(conn *clientMsg) {
	for {
		select {
		case <-time.After(t.heartBoatTime):
			_ = conn.conn.Close()
			return
		case <-(*&conn.heartBoat):
			break
		}
	}
}
