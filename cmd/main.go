package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AlpsMonaco/portforward"
)

func main() {
	f, err := portforward.NewForward("tcp", ":65432", "192.168.1.202:33123")
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	f.Close()
}
