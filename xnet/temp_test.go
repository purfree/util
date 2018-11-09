package temp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/purfree/util/xnet/pnet"
)

func TestNet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	startServer(ctx)
	startClient(ctx)
	Wait(func() {
		cancel()
	})
}

func startServer(ctx context.Context) {
	t := pnet.NewTCPServer()
	err := t.Start(":9090", ctx, func(s *pnet.Session) {
		go func() {
			data, err := s.Receive()
			if err != nil {
				panic(err)
			}
			fmt.Println("server:", string(data))
		}()
	})
	if err != nil {
		panic(err)
	}
}

func startClient(ctx context.Context) {
	c := pnet.NewTCPClient()
	err := c.Start("127.0.0.1:9090", ctx, func(s *pnet.Session) {
		go func() {
			//b := make([]byte, 0)
			//for i := 0; i < 10000000; i++ {
			//	b = append(b, 5)
			//}
			b := []byte("hello world")
			fmt.Println("client:", string(b))
			s.Send(b)
		}()

	})
	if err != nil {
		panic(err)
	}
}

func Wait(back func()) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigc
	back()
}
