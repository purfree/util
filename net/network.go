package net

/*
 数据格式
 \ul\data\
 dl: 数据长度,4byte uint32 4G
 data: 数据
*/

// 数据长度标记位
const (
	dataLen = 4
)

type NetworkMgr struct {
	Route
	TCPServer []*tcpServer
	TCPClient []*tcpClient
}

func NewNetworkMgr() *NetworkMgr {
	return &NetworkMgr{}
}

// 必须先设置服务端参数
func (p *NetworkMgr) StartTCPServer(conf *TCPServerConfig) error {
	ts := newTCPServer()
	if conf == nil {
		ts.TCPServerConfig = DefaultTCPServer
	} else {
		ts.TCPServerConfig = conf
	}

	if err := ts.start(p); err != nil {
		return err
	}
	p.TCPServer = append(p.TCPServer, ts)
	return nil
}

// 必须先设置客户端参数
//func (p *NetworkMgr) StartTCPClient() error {
//	return p.TCPClient.start()
//}
