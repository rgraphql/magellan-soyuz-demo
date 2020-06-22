package main

import (
	"context"
	"net/http"
	"os"

	"github.com/rgraphql/magellan/schema"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	listen = ":8093"
)

func main() {
	app := cli.NewApp()
	app.Name = "server"
	app.Usage = "magellan soyuz demo server"
	app.HideVersion = true
	app.Action = runServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "listen",
			EnvVar:      "LISTEN",
			Usage:       "listen string, default `LISTEN`",
			Value:       listen,
			Destination: &listen,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err.Error())
	}
}

func runServer(_ *cli.Context) error {
	// Construct the websocket server.
	ctx := context.Background()
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	le := logrus.NewEntry(log)

	scm, err := schema.Parse(schemaStr)
	if err != nil {
		return err
	}
	server := NewServer(ctx, le, scm)
	mux := http.NewServeMux()
	mux.Handle("/ws", server)
	le.Infof("starting listener on %s", listen)
	return http.ListenAndServe(listen, mux)
}
