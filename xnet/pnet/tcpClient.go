package pnet

import (
	"context"
	"errors"
	"net"
	"time"
)

type tcpClient struct {
	cfg tcpConfig
	hs  handleSession // 会话处理函数，单条连接

	//aid      atomic.ID
	//sessions map[uint64]*Session

	status bool //true:运行
	//curConnNum int
}

func NewTCPClient() *tcpClient {
	return &tcpClient{
		cfg: DefaultTCPConfig,
		//sessions: make(map[uint64]*Session),
	}
}

func (p *tcpClient) SetConfig(cfg tcpConfig) {
	p.cfg = cfg
}

func (p *tcpClient) Start(address string, ctx context.Context, hs handleSession) error {
	defer func() {
		p.status = false
	}()
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
	count := 0
	for {
		count++
		if count > p.cfg.ClientConnAttemp {
			return errors.New("connect fail")
			break
		}
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			time.Sleep(time.Second * time.Duration(p.cfg.ClientConnAttempinterval))
			continue
		}

		s := newSession()
		s.start(conn, ctx, p.cfg)
		p.hs(s)
		break
	}
	return nil
}
