package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/derekbit/ubdsrv-backend/pkg/message"
)

const (
	TypeRead = iota
	TypeWrite
	TypeResponse
)

func StartServer(c *cli.Context) error {
	sockAddr := c.String("sock-addr")
	backendFile := c.String("backend-file")
	size := c.Int64("size")

	logrus.Infof("sock-addr: %v", sockAddr)
	logrus.Infof("backend-file: %v", backendFile)
	logrus.Infof("size: %v", size)

	f, err := os.OpenFile(backendFile, os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return errors.Wrapf(err, "failed to create backend file %v", backendFile)
	}

	if err := f.Truncate(size); err != nil {
		return errors.Wrapf(err, "failed to truncate %v to %v", backendFile, size)
	}

	if err := os.RemoveAll(sockAddr); err != nil {
		return errors.Wrapf(err, "failed to remove sock file %v", sockAddr)
	}

	l, err := net.Listen("unix", sockAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to listen %v", sockAddr)
	}
	defer l.Close()

	logrus.Info("Waiting for connection...")

	for {
		conn, err := l.Accept()
		if err != nil {
			return errors.Wrapf(err, "failed to accept connection")
		}

		go handleRequest(conn, f)
	}
}

func handleRequest(conn net.Conn, f *os.File) {
	var err error

	defer f.Close()
	defer conn.Close()

	for {
		req := &message.Message{}

		err = req.Read(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("=== closed by client")
				break
			}

			logrus.Error(err)
			break
		}

		// Response
		res := &message.Message{
			MagicVersion: req.MagicVersion,
			Type:         TypeResponse,
			Seq:          req.Seq,
			Offset:       req.Offset,
			Size:         req.Size,
		}

		switch req.Type {
		case TypeRead:
			res.Data = make([]byte, req.Size)

			_, err = f.ReadAt(res.Data, res.Offset)
			if err != nil {
				logrus.WithError(err).Errorf("Failed to read offset=%v size=%v", res.Offset, len(res.Data))
			}
		case TypeWrite:
			_, err = f.WriteAt(req.Data, req.Offset)
			if err != nil {
				logrus.WithError(err).Errorf("Failed to write offset=%v size=%v", res.Offset, len(res.Data))
			}
		}

		if err != nil {
			break
		}

		err = res.Write(conn)
		if err != nil {
			logrus.Error(err)
			break
		}
	}
}
