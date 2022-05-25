package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/andreykyz/nats_jsq/config"
	"github.com/andreykyz/nats_jsq/pubsub"
	log "github.com/sirupsen/logrus"
)

var subscriber string

func init() {
	flag.StringVar(&subscriber, "name", "", "Subscriber name, publisher if blank")
}

func main() {

	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())

	pubsubListener := pubsub.NewPubSub(ctx, cfg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	if subscriber == "" {
		go pubsubListener.SendMsgChan(ctx)
		<-signalChan
		cancel()
	} else {
		for {
			select {
			case <-signalChan:
				cancel()
			case msg := <-pubsubListener.GetMsgChan(ctx):
				log.Infof("%s %s", subscriber, msg)
			}
		}
	}
}
