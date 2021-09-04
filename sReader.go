package GuTcp

import (
	"github.com/shenpan233/GuTcp/util/Pack"
	"github.com/shenpan233/GuTcp/util/UnPack"
	"sync/atomic"
)

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
					Pack.BufAdd(&data, bin[:l])
				}
				if dataLen > 0 {
					if t.onRecvHandlerFunc != nil {
						conn.heartBoat <- 1
						go t.onRecvHandlerFunc(uid, UnPack.BufGet(&data, uint(dataLen)))

					}
				}

			}
		} else {
			if t.onServerClose != nil {
				go t.onServerClose(uid)
			}
			t.user.Delete(uid)
			//连接丢失
			return
		}
	}

}

//accept 创建连接
func (t *GuTcpServer) accept() {
	for {
		conn, err := t.listener.Accept()
		if err == nil {
			atomic.AddUint64(t.userId, 1)
			uid := atomic.LoadUint64(t.userId)
			msg := &clientMsg{
				conn:      conn,
				heartBoat: make(chan int),
			}
			t.user.Store(uid, msg)
			go t.recv(uid)
			go t.heartBoating(msg)
		}
	}

}
