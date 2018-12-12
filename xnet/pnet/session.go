package pnet

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"net"
	"time"
)

type rwChan struct {
	c    chan struct{}
	data []byte
	err  error
}

type Session struct {
	c net.Conn
	r *bufio.Reader
	w *bufio.Writer

	cr chan *rwChan
	cw chan *rwChan

	cancel context.CancelFunc

	dataHeaderSize int
}

func newSession() *Session {
	return &Session{}
}

func (p *Session) start(c net.Conn, ctx context.Context, cfg tcpConfig) {
	c.SetReadDeadline(time.Now().Add(time.Second * time.Duration(cfg.ReadDeadline)))
	c.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(cfg.WriteDeadline)))

	p.r = bufio.NewReader(c)
	p.w = bufio.NewWriter(c)

	p.cr = make(chan *rwChan, 16)
	p.cw = make(chan *rwChan, 16)

	p.dataHeaderSize = cfg.DataHeaderSize

	ctx2, cancel := context.WithCancel(context.Background())

	p.cancel = cancel

	p.run(ctx, ctx2)
}

func (p *Session) run(ctx, ctx2 context.Context) {
	go p.read(ctx, ctx2)
	go p.write(ctx, ctx2)
}

func (p *Session) read(ctx, ctx2 context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ctx2.Done():
			return
		case r := <-p.cr:
			// read the header
			headBuf := make([]byte, p.dataHeaderSize)
			if _, err := io.ReadFull(p.r, headBuf); err != nil {
				r.err = err
				r.c <- struct{}{}
				continue
			}
			msgLen := binary.BigEndian.Uint32(headBuf)

			dataBuf := make([]byte, msgLen)
			if _, err := io.ReadFull(p.r, dataBuf); err != nil {
				r.err = err
				r.c <- struct{}{}
				continue
			}
			r.data = dataBuf
			r.c <- struct{}{}
		}
	}
}

func (p *Session) write(ctx, ctx2 context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ctx2.Done():
			return
		case w := <-p.cw:
			b := make([]byte, len(w.data)+p.dataHeaderSize)
			binary.BigEndian.PutUint32(b, uint32(len(w.data)))
			copy(b[p.dataHeaderSize:], w.data)

			_, err := p.w.Write(b)
			p.w.Flush()
			w.err = err
			w.c <- struct{}{}
		}
	}
}

func (p *Session) Receive() ([]byte, error) {
	r := &rwChan{
		c: make(chan struct{}),
	}
	p.cr <- r
	<-r.c
	return r.data, r.err
}

func (p *Session) Send(data []byte) error {
	w := &rwChan{
		c:    make(chan struct{}),
		data: data,
	}
	p.cw <- w
	<-w.c
	return w.err
}

func (p *Session) Close() error {
	p.cancel()
	return p.c.Close()
}
