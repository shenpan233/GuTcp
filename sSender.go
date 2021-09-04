package GuTcp

import (
	"errors"
	"fmt"
)

//GroupSend 群发数据给所有客户端
func (t *GuTcpServer) GroupSend(bin []byte) {
	t.user.Range(func(key, _ interface{}) bool {
		go t.Send(key.(uint64), bin)
		return true
	})
}

//Send 发送数据给指定客户端
//  --------------------------
//	userId 用户id,OnServerRecvHandlerFunc方法会返回
//  bin 会自动执行粘包处理Encode命令
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
	} else {
		err = errors.New(fmt.Sprintf("无法通过客户端ID(%d)找到此客户端,请检查是否已连接或已断开", UserId))
	}
	return
}
