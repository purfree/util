package net

import (
	"errors"
	"net"
)

type tcpClient struct {
	RemoteAddr    string
	HandleSession func(s *Session)

	addr *net.TCPAddr
}

func newTCPClient() *tcpClient {
	return &tcpClient{}
}

func (p *tcpClient) init() error {
	if p.RemoteAddr == "" {
		return errors.New("please input local address")
	}
	addr, err := net.ResolveTCPAddr("tcp", p.RemoteAddr)
	if err != nil {
		return err
	}
	p.addr = addr

	if p.HandleSession == nil {
		errors.New("please set HandleSession")
	}

	return nil
}

func (p *tcpClient) start() error {
	if err := p.init(); err != nil {
		return err
	}
	go p.run()
	return nil
}

func (p *tcpClient) run() {
	conn, err := net.DialTCP("tcp", nil, p.addr)
	if err != nil {
		// TODO 错误处理，有些错误类型可以执行重连或其他操作
		panic(err)
	}

	// TODO 连接数量限制

	go func(c net.Conn) {
		// TODO 缓冲池
		//sid := p.genSessionID()
		//s := newSession(conn)
		//p.HandleSession(s)
	}(conn)
}

//func (p *TCPServer) genSessionID() uint64 {
//	return p.Add()
//}

func (p *tcpClient) Stop() {

}
