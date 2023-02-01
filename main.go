package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/derekbit/backend-server/pkg/message"
)

const (
	SockAddr = "/tmp/backend.sock"
)

const (
	TypeRead = iota
	TypeWrite
	TypeResponse
)

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		req := &message.Message{}
		err := req.Read(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("=== closed by client")
				break
			}

			logrus.Error(err)
			break
		}

		//fmt.Println("[READ ] ", req)

		res := &message.Message{
			MagicVersion: req.MagicVersion,
			Type:         TypeResponse,
			Data:         make([]byte, req.Size),
			Seq:          req.Seq,
			Offset:       req.Offset,
			Size:         req.Size,
		}

		err = res.Write(conn)
		if err != nil {
			logrus.Error(err)
			break
		}

		//fmt.Println("[WRITE] ", res)
	}
}

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		logrus.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		logrus.WithError(err).Fatal("listen")
	}
	defer l.Close()

	logrus.Info("Waiting for connection...")
	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.WithError(err).Fatal("accept")
		}

		go handleRequest(conn)
	}
}
