package portforward

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

var handleError func(err error)
var errZeroRead = errors.New("read size is zero")
var maxZeroRead int32 = 10

func init() {
	handleError = func(err error) {
		fmt.Println(err)
	}
}

func SetErrorHandler(f func(err error)) {
	handleError = f
}

func closeConn(conn ...net.Conn) {
	for _, c := range conn {
		if err := c.Close(); err != nil {
			handleError(err)
		}
	}
}

func handleStream(srcConn, dstConn net.Conn) error {
	var err error
	n, err := io.Copy(srcConn, dstConn)
	if n == 0 {
		return errZeroRead
	}

	return err
}

func NewForward(network, bindAddr, dstAddr string) (*Forward, error) {
	var f Forward
	if err := f.Bind(network, bindAddr, dstAddr); err != nil {
		return nil, err
	}

	return &f, nil
}

// Forward
// usage :
// xf := new(Forward)
// xf.Bind("tcp",":65432","192.168.1.200:3389")
// Bind() will not blocking the following code.
// this will forward port 65432 's request.

type Forward struct {
	network  string
	listener net.Listener
	to       string
	isQuit   bool
}

func (xf *Forward) Bind(network, bindAddr, dstAddr string) error {
	var err error
	xf.listener, err = net.Listen(network, bindAddr)
	if err != nil {
		return err
	}

	xf.network = network
	xf.to = dstAddr

	go xf.forward()
	return nil
}

func (xf *Forward) Close() {
	xf.isQuit = true
	if err := xf.listener.Close(); err != nil {
		handleError(err)
	}
}

func (xf *Forward) forward() {
	for {
		conn, err := xf.listener.Accept()
		if err != nil {
			handleError(err)
			break
		}

		go xf.handleConn(conn)
	}
}

func (xf *Forward) handleConn(srcConn net.Conn) {
	defer closeConn(srcConn)

	dstConn, err := net.Dial(xf.network, xf.to)
	if err != nil {
		handleError(err)
		return
	}

	defer closeConn(dstConn)
	var isConnErr bool
	var zeroReadCount int32

	addZeroRead := func() int32 {
		return atomic.AddInt32(&zeroReadCount, 1)
	}

	go func() {
		for {
			err := handleStream(srcConn, dstConn)
			if err != nil {
				handleError(err)
				if err == errZeroRead {
					if addZeroRead() > maxZeroRead {
						isConnErr = true
						break
					}

				} else {
					handleError(err)
					isConnErr = true
					break
				}
			}
		}
	}()

	go func() {
		for {
			err := handleStream(dstConn, srcConn)
			if err != nil {
				handleError(err)
				if err == errZeroRead {
					if addZeroRead() > maxZeroRead {
						isConnErr = true
						break
					}

				} else {
					handleError(err)
					isConnErr = true
					break
				}
			}
		}
	}()

	for {
		if xf.isQuit {
			break
		}

		if isConnErr {
			break
		}

		time.Sleep(3 * time.Second)
	}

}
