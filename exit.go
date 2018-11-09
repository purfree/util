package util

import (
	"os"
	"os/signal"
	"syscall"
)

func Wait(back func()) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigc
	back()
}
