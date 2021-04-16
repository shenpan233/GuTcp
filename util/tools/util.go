package tools

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetRandomBin(i int) (data []byte) {
	rand.Seed(time.Now().Unix())
	for i2 := 0; i2 < i; i2++ {
		data = append(data, byte(rand.Intn(255)))
	}
	return
}
func GetRand32() int32 {
	rand.Seed(GetServerCurTime())
	return int32(rand.Int())
}

func RandUint32(min, max uint32) uint32 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return uint32(rand.Int31n(int32(max-min)) + int32(min))
}

func GetServerCurTime() int64 {
	return time.Now().Unix()
}
func HextoBin(HexData string) []byte {
	tmp, _ := hex.DecodeString(strings.ReplaceAll(HexData, " ", ""))
	return tmp
}
func BintoHex(Bin []byte) string {
	return strings.ToUpper(hex.EncodeToString(Bin))
}
func BintoHex2(Bin *[]byte) string {
	return strings.ToUpper(hex.EncodeToString(*Bin))
}
func ToMd5Bytes(data []byte) (ret []byte) {
	tmp := md5.Sum(data)
	ret = tmp[:]
	return
}
func IpToInt(Ip string) int {
	split := strings.Split(Ip, ".")
	if len(split) != 4 {
		return 0
	}
	var intIp []int
	for _, s := range split {
		tmp, _ := strconv.Atoi(s)
		intIp = append(intIp, tmp)
	}
	//JAVA反编译↓
	//Long.parseLong(split[0]) + (Long.parseLong(split[3]) << 24) + (Long.parseLong(split[2]) << 16) + (Long.parseLong(split[1]) << 8);

	return intIp[0] + (intIp[3] << 24) + (intIp[2] << 16) + (intIp[1] << 8)
}
