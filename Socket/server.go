package Socket

import (
	"errors"
	"fmt"
	"github.com/shenpan233/guTcpSocket/util/Pack"
	"github.com/shenpan233/guTcpSocket/util/UnPack"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type GuTcpServer struct {
	listener          net.Listener
	user              *sync.Map
	userId            *uint64
	heartBoatTime     int64
	onRecvHandlerFunc OnServerRecvHandlerFunc
}

type clientMsg struct {
	conn  net.Conn
	cache []byte
}

//Creat 创建tcp服务,成功返回true,失败返回false并在控制台打印错误信息
func (t *GuTcpServer) Creat(port int) bool {
	if t.user == nil {
		t.userId = new(uint64)
		t.user = new(sync.Map)
	}

	var err error
	t.listener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		go t.accept()
		return true
	} else {
		fmt.Println(err.Error())

		return false
	}
}

//SetHeartBoatTime 设置心跳时间
//  time 单位秒
func (t *GuTcpServer) SetHeartBoatTime(i int) {
	t.heartBoatTime = int64(time.Duration(i) * time.Second)
}

//BindingEvent 绑定事件
func (t *GuTcpServer) BindingEvent(OnRecvHandlerFunc OnServerRecvHandlerFunc) {
	t.onRecvHandlerFunc = OnRecvHandlerFunc
}

//accept 创建连接
func (t *GuTcpServer) accept() {
	for {
		conn, err := t.listener.Accept()
		if err == nil {
			atomic.AddUint64(t.userId, 1)
			uid := atomic.LoadUint64(t.userId)
			t.user.Store(uid, &clientMsg{
				conn: conn,
			})
			go t.recv(uid)
		}
	}
}

//OnServerRecvHandlerFunc 接收参数回调程序
// 	uid 用户标识id
//	bin 已截取长度的封包
type OnServerRecvHandlerFunc func(uid uint64, bin []byte)

//recv 接收数据
func (t *GuTcpServer) recv(uid uint64) {

	tmp, _ := t.user.Load(uid)
	conn := tmp.(*clientMsg)
	bin := [5024]byte{}
	defer func() {
		_ = conn.conn.Close()
		if recover() != nil {
			return
		}
	}()
	for {
		l, err := conn.conn.Read(bin[:])
		if err == nil {
			//粘包处理
			data := bin[:l]
			for len(data) >= 4 {
				dataLen := int(UnPack.BufToInt32(&data) - 4)
				if dataLen > len(data) { //粘包(数据不足)
					l, _ = conn.conn.Read(bin[:])
					Pack.Buf_Add(&data, bin[:l])
				}
				go t.onRecvHandlerFunc(uid, UnPack.BufGet(&data, uint(dataLen)))

			}
		} else {
			t.user.Delete(uid)
			//连接丢失
			return
		}
	}

}

//GroupSend 群发数据给所有客户端
func (t *GuTcpServer) GroupSend(bin []byte) {
	t.user.Range(func(key, _ interface{}) bool {
		go t.Send(key.(uint64), bin)
		return true
	})
}

//Send 发送数据给指定客户端
//  --------------------------
//  bin 会自动执行粘包编码Encode命令
//  err 返回错误信息,nil就是发送成功
func (t *GuTcpServer) Send(UserId uint64, bin []byte) (err error) {
	if len(bin) == 0 {
		err = errors.New("空数据无法发送")
		return
	}
	v, ok := t.user.Load(UserId)
	if ok {
		conn := v.(*clientMsg)
		Encode(&bin)
		_, err = conn.conn.Write(bin)
		fmt.Println(bin)
	} else {
		err = errors.New(fmt.Sprintf("无法通过客户端ID(%d)找到此客户端,请检查是否已连接或已断开", UserId))
	}
	return
}
