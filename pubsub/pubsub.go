package pubsub

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/andreykyz/nats_jsq/config"

	nats "github.com/nats-io/nats.go"
)

//go:generate msgp
type Message struct {
	Time    string `msg:"time,omitempty"`
	Counter int    `msg:"counter,omitempty"`
}

func (m *Message) String() string {
	return fmt.Sprintf("%d %s", m.Counter, m.Time)
}

//msgp:ignore PubSub
type PubSub struct {
	subj           string
	consumer       string
	consumerFilter string
	js             nats.JetStreamContext
	msg            chan *Message
}

func NewPubSub(ctx context.Context, c *config.Config) *PubSub {
	ps := &PubSub{subj: c.SubjectFilter, consumer: c.Consumer, consumerFilter: c.ConsumerFilter}
	ps.msg = make(chan *Message)

	opts := nats.Options{
		Url:            c.BrokerURL,
		Name:           "Co-co-co",
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
	}
	nc, err := opts.Connect()
	nc.Flush()
	if err != nil {
		log.Error(err)
		close(ps.msg)
		return nil
	}
	ps.js, err = nc.JetStream()
	if err != nil {
		log.Error(err)
		close(ps.msg)
		return nil
	}
	go func() {
		<-ctx.Done()
		log.Info("Broker done")
		nc.Close()
	}()
	return ps
}

func (ps *PubSub) GetMsgChan(ctx context.Context) chan *Message {

	if sub, err := ps.js.PullSubscribe(fmt.Sprintf("%s.*", ps.subj), ps.consumer); err == nil {
		go func() {
			defer close(ps.msg)
			for {
				message := &Message{}
				if bMsg, err := sub.Fetch(1); err == nil {
					log.Info(len(bMsg))
					if _, err := message.UnmarshalMsg(bMsg[0].Data); err == nil {
						ps.msg <- message
					} else {
						log.Error(err)
					}
				} else if err == nats.ErrTimeout {
					log.Info(err)
				} else {
					log.Error(err)
					return
				}
			}
		}()
	} else {
		log.Fatal(err)
	}
	return ps.msg
}

func (ps *PubSub) SendMsgChan(ctx context.Context) {
	tick := time.NewTicker(2 * time.Second)
	counter := 0
	for {
		select {
		case <-ctx.Done():
			tick.Stop()
			return
		case <-tick.C:
			msg := &Message{
				Time:    time.Now().String(),
				Counter: counter,
			}
			counter += 1
			b, err := msg.MarshalMsg(nil)
			if err != nil {
				log.Error(err)
			}
			n, err := ps.js.Publish(fmt.Sprintf("%s.%s", ps.subj, ps.consumerFilter), b)
			if err != nil {
				log.Error(err)

			}
			log.Info(n)
		}
	}
}
