package pnet

import (
	"context"
	"errors"
	"net"
)

//var sp *sync.Pool

//func init() {
//sp = &sync.Pool{New: func() interface{} {
//	return newSession()
//}}
//}

type tcpServer struct {
	cfg tcpConfig
	hs  handleSession // 会话处理函数，单条连接

	//aid      atomic.ID
	//sessions map[uint64]*Session

	status bool //true:运行
	//curConnNum int
}

func NewTCPServer() *tcpServer {
	return &tcpServer{
		cfg: DefaultTCPConfig,
		//sessions: make(map[uint64]*Session),
	}
}

func (p *tcpServer) SetConfig(cfg tcpConfig) {
	p.cfg = cfg
}

func (p *tcpServer) Start(address string, ctx context.Context, hs handleSession) error {
	if p.status {
		return errors.New("运行中")
	}

	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	if hs == nil {
		return errors.New("请输入处理函数")
	}
	p.hs = hs
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	p.status = true
	go p.run(l, ctx)
	return nil
}

func (p *tcpServer) run(l *net.TCPListener, ctx context.Context) {
	defer func() {
		p.status = false
		//p.curConnNum = 0
	}()

	for {
		select {
		case <-ctx.Done():
			break
		default:
			conn, err := l.Accept()
			if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
				//log.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				//log.Debug("Read error", "err", err)
				panic(err)
			}

			//if p.cfg.MaxConnNum != 0 {
			//	if p.curConnNum >= p.cfg.MaxConnNum {
			//		// 已达最大连接数量，不再接受连接
			//		continue
			//	}
			//	p.curConnNum++
			//}

			go func() {
				//s := sp.Get().(*Session)
				s := newSession()
				s.start(conn, ctx, p.cfg)
				p.hs(s)
			}()
		}
	}
}
