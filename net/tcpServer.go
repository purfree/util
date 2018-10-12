package net

import (
	"errors"
	"net"
	"sync"

	"github.com/purfree/common/atomic"
)

var DefaultTCPServer = &TCPServerConfig{}

type TCPServerConfig struct {
	LocalAddr  string // 监听地址
	MaxConnNum uint64 // 最大连接数量,等于0无限制
}

type tcpServer struct {
	*TCPServerConfig
	HandleSession func(s *Session) // 会话处理函数，单条连接

	CurConnNum uint64 // 当前连接数量

	addr *net.TCPAddr

	nm *NetworkMgr

	aid      atomic.ID
	sessions map[uint64]*Session
	spool    *sync.Pool
}

func newTCPServer() *tcpServer {
	return &tcpServer{
		sessions: make(map[uint64]*Session),
		spool: &sync.Pool{New: func() interface{} {
			return newSession()
		}},
	}
}

func (p *tcpServer) init() error {
	if p.LocalAddr == "" {
		return errors.New("please input local address")
	}
	addr, err := net.ResolveTCPAddr("tcp", p.LocalAddr)
	if err != nil {
		return err
	}
	p.addr = addr

	if p.HandleSession == nil {
		errors.New("please set HandleSession")
	}

	return nil
}

func (p *tcpServer) start(nm *NetworkMgr) error {
	if err := p.init(); err != nil {
		return err
	}
	p.nm = nm
	go p.run()
	return nil
}

func (p *tcpServer) run() error {
	l, err := net.ListenTCP("tcp", p.addr)
	if err != nil {
		return err
	}
	for {
		//TODO 超时控制
		conn, err := l.Accept()
		if err != nil {
			// TODO 错误处理，有些错误类型可以执行重连或其他操作???
			panic(err)
		}

		if p.MaxConnNum != 0 {
			if p.CurConnNum >= p.MaxConnNum {
				// 已达最大连接数量，不再接受连接
				continue
			}
			p.CurConnNum++
		}

		// 启动会话
		sid := p.genSessionID()
		s := p.spool.Get().(*Session)
		//s.conn = conn
		//s.nm = p.nm
		//s.id = sid
		s.init(sid, p.nm, conn)
		p.sessions[sid] = s
		go s.start()

	}
}
func (p *tcpServer) genSessionID() uint64 {
	return p.aid.Add()
}

func (p *tcpServer) Stop() {

}
