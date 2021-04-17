package Socket

import (
	"fmt"
	"github.com/shenpan233/guTcpSocket/util/Pack"
	"github.com/shenpan233/guTcpSocket/util/UnPack"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type GuTcpClient struct {
	conn net.Conn
	//服务器IP
	mServerIp string
	mbin      *[]byte
	maps      *sync.Map
	mapsCount int32
	ssoSeq    int32
	wg        sync.WaitGroup
}

func (t *GuTcpClient) ConnectToServer(Ip, port string) bool {
	var err error
	t.conn, err = net.Dial("tcp", Ip+":"+port)
	if err != nil {
		fmt.Println("GuTcpClient", "Connect failed error:", err)
		return false
	}
	t.mServerIp = Ip
	fmt.Println("GuTcpClient", "connect Successfully")
	return true
}

func (t *GuTcpClient) CloseLink() {
	_ = t.conn.Close()
}

//SendDataAndGetStatus 发送并获取发送状态
func (t *GuTcpClient) SendDataAndGetStatus(data *[]byte) bool {
	_, err := t.conn.Write(*data)
	if err != nil {
		return false
	}
	return true
}

//SendAndGetData 发送并拉取封包
func (t *GuTcpClient) SendAndGetData(SsoSeq, waitTime int32, data *[]byte) []byte {
	_, err := t.conn.Write(*data)
	if err != nil {
		return nil
	}
	return t.pullData(SsoSeq, waitTime)
}

//UselessDataJoinCache 对于分析无用的封包加入缓存
func (t *GuTcpClient) UselessDataJoinCache(SsoSeq int32, data *[]byte) {
	if atomic.LoadInt32(&t.mapsCount) >= 60 {
		//重置
		atomic.StoreInt32(&t.mapsCount, 0)
		t.maps = &sync.Map{}
	}
	atomic.AddInt32(&t.mapsCount, 1)
	t.maps.Store(SsoSeq, *data)
}

//pullData 拉取封包
func (t *GuTcpClient) pullData(SsoSeq, waitTime int32) []byte {
	//异常处理，怕服务端不规则封包导致出问题
	defer func() {
		catch := recover()
		if catch != nil {
			return
		}
	}()
	for i := 0; i < (int(waitTime) / 100); i++ {
		data, ok := t.maps.LoadAndDelete(SsoSeq)
		if ok {
			atomic.AddInt32(&t.mapsCount, -1)
			return data.([]byte)
		}
		time.Sleep(100 * time.Millisecond) //每隔100毫秒拉取一次
	}
	return nil
}

//OnClientRecvHandlerFunc 接收参数回调程序
//	bin 已截取长度的封包
type OnClientRecvHandlerFunc func(bin []byte)

func (t *GuTcpClient) Listen(OnRecvHandlerFunc OnClientRecvHandlerFunc) {
	t.wg.Add(1)

	//异常捕获 你永远不知道发来的是什么
	defer func() {
		t.wg.Done()
		catch := recover()
		if catch != nil {
			t.mbin = nil
			go t.Listen(OnRecvHandlerFunc) //重启服务
			return
		}
	}()

	var data [5120]byte //5KB的缓冲区
	var (
		length, length2 int
	)
	var err error
	for {
		for i := 0; i < 20; i++ {
			length, err = t.conn.Read(data[:])
			if err != nil {
				break
			}
			if length == 0 {
				break
			}
			Pack.BufAdd(t.mbin, data[:length]) //粘包处理
			for len(*t.mbin) >= 4 {            //这个4是开头的长度
				length2 = int(UnPack.BufToInt32(t.mbin) - 4)
				if length >= 51200 { //50KB清空缓存
					*t.mbin = nil
				} else if length2 <= length {
					if len(*t.mbin) >= length2 {
						go OnRecvHandlerFunc(UnPack.BufGet(t.mbin, uint(length2)))
					}
				}

			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
