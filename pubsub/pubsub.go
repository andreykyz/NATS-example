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
	topic string
	group string
	js    nats.JetStreamContext
	msg   chan *Message
}

func NewPubSub(ctx context.Context, c *config.Config) *PubSub {
	ps := &PubSub{topic: c.Topic, group: c.Group}
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

	if sub, err := ps.js.QueueSubscribeSync(ps.topic, ps.group); err == nil {
		go func() {
			defer close(ps.msg)
			for {
				message := &Message{}
				if bMsg, err := sub.NextMsgWithContext(ctx); err == nil {
					if err, _ := message.UnmarshalMsg(bMsg.Data); err == nil {
						ps.msg <- message
					} else {
						log.Error(err)
					}
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
	tick := time.NewTicker(5 * time.Second)
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
			b, err := msg.MarshalMsg(nil)
			if err != nil {
				log.Error(err)
			}
			ps.js.Publish(ps.topic, b)
		}
	}
}
