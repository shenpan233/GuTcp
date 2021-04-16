package Pack

import (
	"AndroidQQ/src/util/tools"
)

func Int8_to_buf(bArr *[]byte, i int8) {
	*bArr = append(*bArr, (byte)(i>>0))
}

func Int16_to_buf(bArr *[]byte, i int16) {
	*bArr = append(*bArr, (byte)(i>>8), (byte)(i>>0))
}

func Int32_to_buf(bArr *[]byte, i int32) {
	*bArr = append(*bArr, (byte)(i>>24), (byte)(i>>16), (byte)(i>>8), (byte)(i>>0))
}
func Int64_to_buf(bArr *[]byte, i int64) {
	*bArr = append(*bArr, (byte)(i>>56), (byte)(i>>48), (byte)(i>>40), (byte)(i>>32), (byte)(i>>24), (byte)(i>>16), (byte)(i>>8), (byte)(i>>0))
}

func Int64_to_buf32(bArr *[]byte, i int64) {
	*bArr = append(*bArr, (byte)(i>>24), (byte)(i>>16), (byte)(i>>8), (byte)(i>>0))
}
func Hex_to_buf(bArr *[]byte, hexData string) {
	*bArr = append(*bArr, tools.HextoBin(hexData)...)
}
func Buf_Add(bArr *[]byte, Need []byte) {
	*bArr = append(*bArr, Need...)
}
func Buf_Add2(bArr *[]byte, Need *[]byte) {
	*bArr = append(*bArr, *Need...)
}
func BinToken_to_buf(bArr *[]byte, Data []byte) {
	Int16_to_buf(bArr, int16(len(Data)))
	*bArr = append(*bArr, Data...)
}
func BinToken_to_buf2(bArr *[]byte, Data *[]byte) {
	Int16_to_buf(bArr, int16(len(*Data)))
	*bArr = append(*bArr, *Data...)
}

func StrToken_To_Buf(bArr *[]byte, Data string) {
	BinToken_to_buf(bArr, []byte(Data))
}
