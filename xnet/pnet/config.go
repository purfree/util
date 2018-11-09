package pnet

type handleSession func(s *Session)

var DefaultTCPConfig = tcpConfig{
	ReadDeadline:   30,
	WriteDeadline:  30,
	DataHeaderSize: 4,

	ClientConnAttemp:         5,
	ClientConnAttempinterval: 1,
}

type tcpConfig struct {
	//MaxConnNum    int // 最大连接数量,等于0无限制
	ReadDeadline  int // 读超时（秒）
	WriteDeadline int // 写超时（秒）

	DataHeaderSize int //数据头长度

	ClientConnAttemp         int //连接尝试次数
	ClientConnAttempinterval int //连接失败，重新连接的间隔，秒
}

type tempError interface {
	Temporary() bool
}
