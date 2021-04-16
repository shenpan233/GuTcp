package Socket

import (
	"GuTcpSocket/util/Pack"
)

// Encode 将消息编码
func Encode(message *[]byte) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(*message) + 4)
	tmp := *message
	*message = nil
	Pack.Int32_to_buf(message, length)
	Pack.Buf_Add(message, tmp)

	return
}
