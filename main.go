package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/derekbit/ubdsrv-backend/pkg/server"
)

const (
	SockAddr    = "/tmp/backend.sock"
	BackendFile = "/data/ubdsrv-backend"
)

func StartCommand() cli.Command {
	return cli.Command{
		Name: "start",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sock-addr",
				Value: SockAddr,
				Usage: "Path to unix domain socket",
			},
			cli.StringFlag{
				Name:  "backend-file",
				Value: BackendFile,
				Usage: "Path to backend file",
			},
			cli.Uint64Flag{
				Name:  "size",
				Value: 1073741824,
				Usage: "Size of backend file",
			},
		},
		Action: func(c *cli.Context) {
			if err := server.StartServer(c); err != nil {
				logrus.Fatalf("Error starting manager: %v", err)
			}
		},
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	a := cli.NewApp()
	a.Usage = "ubdsrv-backend"

	a.Commands = []cli.Command{
		StartCommand(),
	}

	if err := a.Run(os.Args); err != nil {
		logrus.Fatalf("Critical error: %v", err)
	}
}

/*
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
*/
