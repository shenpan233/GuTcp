package Pack

func Int32ToBuf(bArr *[]byte, i int32) {
	*bArr = append(*bArr, (byte)(i>>24), (byte)(i>>16), (byte)(i>>8), (byte)(i>>0))
}
func BufAdd(bArr *[]byte, Need []byte) {
	*bArr = append(*bArr, Need...)
}
