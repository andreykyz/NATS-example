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
	flag.Parse()
}

func main() {

	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())

	pubsubListener := pubsub.NewPubSub(ctx, cfg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	if subscriber == "" {
		log.Info("Run publisher")
		go pubsubListener.SendMsgChan(ctx)
		<-signalChan
		cancel()
	} else {
		log.Infof("Run subscriber %s", subscriber)
		msgChan := pubsubListener.GetMsgChan(ctx)
		for {
			select {
			case <-signalChan:
				cancel()
			case msg := <-msgChan:
				log.Infof("Subscriber %s recive %s", subscriber, msg)
				if msg == nil {
					log.Fatal("msg is nil")
				}
			}
		}
	}
}
