package GuTcp

import (
	"net"
	"sync"
	"time"
)

type GuTcpServer struct {
	listener          net.Listener
	user              *sync.Map
	userId            *uint64
	heartBoatTime     time.Duration
	onRecvHandlerFunc OnServerRecvHandlerFunc
	onServerClose     OnServerClose
}

type clientMsg struct {
	conn      net.Conn
	cache     []byte
	heartBoat chan int
}

//OnServerRecvHandlerFunc 接收参数回调程序
// 	uid 用户标识id
//	bin 已截取长度的封包
type OnServerRecvHandlerFunc func(userId uint64, bin []byte)

//OnServerClose 客户连接关闭参数回调程序
// 	uid 用户标识id
type OnServerClose func(userId uint64)
