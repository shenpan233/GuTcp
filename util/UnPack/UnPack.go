package UnPack

func BufToInt8(data *[]byte) int8 {
	bArr := BufGet(data, 1)
	if bArr == nil {
		return 0
	}
	return int8(bArr[0] & 255)
}

func BufToInt16(data *[]byte) int16 {
	bArr := BufGet(data, 2)
	if bArr == nil {
		return 0
	}
	return int16(((int(bArr[0]) << 8) & 65280) + ((int(bArr[1]) << 0) & 255))
}

func BufToInt32(data *[]byte) int32 {
	bArr := BufGet(data, 4)
	if bArr == nil {
		return 0
	}
	return int32((int(bArr[0])<<24)&-16777216 + ((int(bArr[1]) << 16) & 16711680) + ((int(bArr[2]) << 8) & 65280) + (int(bArr[3]<<0) & 255))

}

//func BufToInt64(data *[]byte) int64 {
//	bArr := BufGet(data, 8)
//	if bArr == nil {
//		return 0
//	}
//	return int64(0 + ((int(bArr[0]) << 56) & -72057594037927936) + (int((bArr[1])<<48) & 71776119061217280) + (((int(bArr[2]) << 40) & 280375465082880) + ((int(bArr[3]) << 32) & 1095216660480) + ((int(bArr[4]) << 24) & 4278190080) + ((int(bArr[5]) << 16) & 16711680) + ((int(bArr[6]) << 8) & 65280) + ((int(bArr[7]) << 0) & 255)))
//
//}

func BufGet(data *[]byte, a uint) []byte {
	if data == nil || len(*data) == 0 || uint(len(*data)) < a {
		return nil
	}
	tmp := (*data)[:a]
	*data = (*data)[a:]
	return tmp

}
