package main

import (
	"context"

	"github.com/StudioSol/async"
	log "github.com/sirupsen/logrus"
	"github.com/tasylab/fantasy"
	"github.com/tasylab/fantasy/backends/bigcache"
)

func main() {
	log.Info("STARTING FANTASY")
	f := fantasy.New(bigcache.New())
	err := async.Run(context.TODO(),
		func(_ context.Context) error {
			log.WithField("ADDRESS", "127.0.0.1:3687").Info("STARTING HTTP SERVER")
			return f.HTTPServer("127.0.0.1:3687")
		},
		func(_ context.Context) error {
			log.WithField("ADDRESS", "127.0.0.1:3688").Info("STARTING TCP SERVER")
			return f.TCPServer("127.0.0.1:3688")
		},
	)
	if err != nil {
		log.WithError(err).Fatal("CRASHED")
	}
}
