package mb

import (
	"github.com/helderfarias/go-api-kit/endpoint"
	"github.com/sirupsen/logrus"
)

type MessageBroker interface {
	Sub() Subscriber
	Pub() Publisher
}

type Subscriber interface {
	Delivery(receivers ...SubReceiver) error
}

type Publisher interface {
	Publish(stream, subj string, sender interface{}, o ...PubOpts) (interface{}, error)
}

type options struct {
	action endpoint.Endpoint
	args   map[string]interface{}
}

type SubReceiver func(o *options)

type PubOpts func(o *options)

type emptySub struct {
}

type emptyPub struct {
}

func (e *emptySub) Delivery(receivers ...SubReceiver) error {
	logrus.Warn("Subscriber not working...")
	return nil
}

func (e *emptyPub) Publish(stream, subj string, sender interface{}, o ...PubOpts) (interface{}, error) {
	logrus.Warn("Publisher not working...")
	return nil, nil
}
