package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

// 数据长度标记位
const (
	msgLenFlag     = 4
	defaultBufSize = 4096
)

type Session struct {
	id   uint64
	nm   *NetworkMgr
	conn net.Conn

	reader *bufio.Reader
	writer *bufio.Writer
	data   bytes.Buffer
	msgLen uint32

	ClientInstance interface{} // 客户端在服务器的实例，如果有
}

func newSession() *Session {
	return &Session{}
}

func (p *Session) init(id uint64, nm *NetworkMgr, conn net.Conn) {
	p.conn = conn
	p.nm = p.nm
	p.id = id

	//p.data = new(bytes.Buffer)
	//TODO 提供参数可修改缓存空间，缓存空间默认4096
	p.reader = bufio.NewReader(p.conn)
	p.writer = bufio.NewWriter(p.conn)
}

func (p *Session) readFromData() (bool, error) {
	if p.data.Len() == 0 {
		return false, nil
	}
	var msgLen uint32 = p.msgLen
	if msgLen == 0 {
		flag := make([]byte, msgLenFlag)
		p.data.Read(flag)
		msgLen = binary.BigEndian.Uint32(flag)
		p.msgLen = msgLen
		if msgLen == 0 {
			return false, errors.New("非法数据")
		}
	}

	//fmt.Println(p.data.Len())
	if int64(msgLen) <= int64(p.data.Len()) {
		data := make([]byte, msgLen)
		p.data.Read(data)
		fmt.Println(string(data), p.data.Len(), p.data.Cap())
		//fmt.Println(p.data.Cap())
		p.msgLen = 0
		return true, nil
	}
	return false, nil
}

func (p *Session) readFromConn() error {
	buf := make([]byte, defaultBufSize)
	n, err := p.conn.Read(buf)
	if err != nil {
		return err
	}
	p.data.Write(buf[:n])
	return nil
}

func (p *Session) start() {
	go p.read()
	go p.write()

	//for {
	//	if b, err := p.readFromData(); err != nil {
	//		fmt.Println(err)
	//		return
	//	} else if b {
	//		continue
	//	}
	//
	//	p.readFromConn()
	//
	//	//b, _ := p.ReadMsg()
	//	//fmt.Println(string(b))
	//}

	//time.Sleep(time.Second * 10)
	//fmt.Println("server start")

	//buf := make([]byte, 4096)
	//for {
	//var data []byte
	//var err error
	//////读取指令头 返回输入流的前4个字节，不会移动读取位置
	//data, err = p.reader.Peek(MsgLenFlag)
	//if len(data) == 0 || err != nil {
	//	continue
	//}
	//
	////返回缓冲中现有的可读取的字节数
	//var byteSize = p.reader.Buffered()
	//fmt.Printf("读取字节长度：%d %d\n", byteSize)
	////生成一个字节数组，大小为缓冲中可读字节数
	//data = make([]byte, byteSize)
	////读取缓冲中的数据
	//p.reader.Read(data)
	//fmt.Printf("读取字节：%s\n", string(data))
	////保存到新的缓冲区
	//p.data.Write(data)
	//
	//if p.data.Len() < MsgLenFlag {
	//	//数据包缓冲区清空
	//	p.data.Reset()
	//	fmt.Printf("非法数据，无指令头...\n")
	//	continue
	//}
	//
	//bodyLength := binary.BigEndian.Uint32(data)
	//
	////判断数据包缓冲区的大小是否小于协议请求头中数据包大小
	////如果小于，等待读取下一个客户端数据包，否则对数据包解码进行业务逻辑处理
	//if bodyLength > uint32(p.data.Len())-MsgLenFlag {
	//	fmt.Printf("body体长度：%d,读取的body体长度：%d\n", bodyLength, p.data.Len())
	//	continue
	//}
	//for {
	//	if bodyLength <= uint32(p.data.Len())-MsgLenFlag {
	//		buf := make([]byte, MsgLenFlag+bodyLength)
	//		_, err := p.data.Read(buf)
	//		if err != nil {
	//			panic(err)
	//		}
	//		fmt.Sprintf("接受到的数据: %s", string(buf[MsgLenFlag:]))
	//	}
	//}

	//p.conn.Read()
	//
	//b, err := ioutil.ReadAll(p.conn)
	////n, err := p.data.ReadFrom(p.conn)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("读取字节：%d\n", b)

	//fmt.Printf("实际处理字节：%v\n", this.Data)
	//p = protocol.Decode(this.Data)
	//逻辑处理
	//go this.logicHandler(p)
	//数据包缓冲区清空
	//this.Data = []byte{}
	//}
	//}
}

func (p *Session) read() {
	for {
		dl := p.parseDataLen()
		data := make([]byte, dl)

		buf := bytes.Buffer{}
		for {
			if int64(buf.Len()) >= int64(dl) {
				break
			}
		}

		//buf := make([]byte, 1024)
		//n, err := p.reader.Read(buf)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(string(buf[:n]))

		//buf := make([]byte, 1024)
		//n, err := p.reader.Read(buf)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(string(buf[:n]))
	}
}

func (p *Session) write() {

}

// 获取命令长度
//func (p *Session) parseCMDLen() uint16 {
//	clb := make([]byte, cmdLen)
//	_, err := io.ReadFull(p.conn, clb)
//	if err != nil {
//		// TODO 未知原因导致的错误
//		fmt.Println(err)
//	}
//	cl := binary.BigEndian.Uint16(clb)
//	return cl
//}

func (p *Session) parseDataLen() uint32 {
	// 获取数据长度
	dlb := make([]byte, dataLen)
	_, err := io.ReadFull(p.conn, dlb)
	if err != nil {
		// TODO 未知原因导致的错误
		fmt.Println(err)
	}
	dl := binary.BigEndian.Uint32(dlb)
	return dl
}

func (p *Session) parseCMD() {

}

func (p *Session) ReadMsg() ([]byte, error) {
	ml := make([]byte, msgLenFlag)
	_, err := io.ReadFull(p.conn, ml)
	if err != nil {
		panic(err)
	}
	d := binary.BigEndian.Uint32(ml)
	b := make([]byte, d)
	n, err := io.ReadFull(p.conn, b)
	if err != nil {
		panic(err)
	}
	// TODO 解析器，根据传入的参数将字节数组转换为对应的类型，例如转换成protobuf
	return b[:n], nil
}

func (p *Session) WriteMsg(b []byte) (int, error) {
	data := make([]byte, len(b)+msgLenFlag)
	binary.BigEndian.PutUint32(data, uint32(len(b)))
	copy(data[msgLenFlag:], b)
	return p.conn.Write(data)
}
