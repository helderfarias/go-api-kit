package mb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/helderfarias/go-api-kit/endpoint"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type natsServer struct {
	nc *nats.Conn
}

type natsSubscriber struct {
	nc   *nats.Conn
	js   nats.JetStreamContext
	done chan error
}

type natsPublisher struct {
	nc *nats.Conn
}

func NewNatsServer() *natsServer {
	opts := []nats.Option{}

	if viper.GetString("nats_client_name") != "" {
		opts = append(opts, nats.Name(viper.GetString("nats_client_name")))
	}

	if viper.GetString("nats_auth_token") != "" {
		opts = append(opts, nats.Token(viper.GetString("nats_auth_token")))
	}

	nc, err := nats.Connect(viper.GetString("nats_servers"), opts...)
	if err != nil {
		logrus.Warnf("nats -> %v", err)
	}

	return &natsServer{nc: nc}
}

func NewNatsPublisher(nc *nats.Conn) *natsPublisher {
	return &natsPublisher{
		nc: nc,
	}
}

func NewNatsSubscriber(nc *nats.Conn) *natsSubscriber {
	js, err := nc.JetStream()
	if err != nil {
		logrus.Warn(err)
	} else {
		logrus.Info("getting jetstream context")
	}

	return &natsSubscriber{
		nc:   nc,
		js:   js,
		done: make(chan error),
	}
}

func (p *natsServer) Pub() Publisher {
	if p.nc == nil {
		return &emptyPub{}
	}

	return NewNatsPublisher(p.nc)
}

func (p *natsServer) Sub() Subscriber {
	if p.nc == nil {
		return &emptySub{}
	}

	return NewNatsSubscriber(p.nc)
}

func (p *natsPublisher) Publish(stream, subj string, sender interface{}, opts ...PubOpts) (interface{}, error) {
	js, err := p.nc.JetStream()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	str, err := js.StreamInfo(stream)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if str == nil {
		return nil, fmt.Errorf("stream not found: %v", stream)
	}

	var data []byte

	if b, ok := sender.([]byte); ok {
		data = b
	} else if s, ok := sender.(string); ok {
		data = []byte(s)
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		data = b
	}

	args := []nats.PubOpt{}

	for _, fill := range opts {
		o := options{args: map[string]interface{}{}}
		fill(&o)

		if values, ok := o.args["pub_args"].([]nats.PubOpt); ok {
			args = append(args, values...)
		}
	}

	pub, err := js.Publish(subj, data, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return pub, nil
}

func (c *natsSubscriber) Delivery(receivers ...SubRec) error {
	subs := []*nats.Subscription{}

	for _, fillArgs := range receivers {
		opts := options{
			args: map[string]interface{}{},
		}

		fillArgs(&opts)

		sub, err := c.js.Subscribe(
			opts.args["subject"].(string),
			c.withServeMsg(opts.action),
			nats.Durable(opts.args["consumer"].(string)),
		)
		if err != nil {
			logrus.Warn(err)
		}

		subs = append(subs, sub)
	}

	logrus.Infof("subscribers for messages")

	<-c.done

	logrus.Infof("drain subscribers")
	for _, s := range subs {
		if err := s.Drain(); err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

func (c *natsSubscriber) Shutdown() {
	logrus.Infof("drain nats connection")

	if err := c.nc.Drain(); err != nil {
		logrus.Error(err)
	}

	c.done <- nil
}

func (c *natsSubscriber) Subscribe(stream, subj, cons string, cmd endpoint.Endpoint) SubRec {
	return func(o *options) {
		o.action = cmd
		o.args["subject"] = subj
		o.args["consumer"] = cons

		str, err := c.js.StreamInfo(stream)
		if err != nil {
			logrus.Warn(err)
		}

		if str == nil {
			logrus.Infof("creating stream %q and subject %q", stream, subj)

			info, err := c.js.AddStream(&nats.StreamConfig{
				Name:     stream,
				Subjects: []string{subj},
			})

			if err != nil {
				logrus.Warnf("js stream -> %v, %v", err, info)
			}
		}
	}
}

func (c *natsSubscriber) withServeMsg(fireAndForget endpoint.Endpoint) nats.MsgHandler {
	return func(msg *nats.Msg) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		fireAndForget(ctx, msg)
	}
}
