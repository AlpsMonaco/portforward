package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlpsMonaco/portforward"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("error param")
		PrintHelp()
		os.Exit(1)
	}

	src := os.Args[1]
	dst := os.Args[2]

	f, err := portforward.NewForward("tcp", src, dst)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	f.Close()
}

func PrintHelp() {
	fmt.Println(`
	usage: portforward [src] [dst]
	src could be ":80","127.0.0.1:80","0.0.0.0:80"
	dst could be "target.com:12345" "172.31.0.4:889"
	`)
}
