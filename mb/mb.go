package mb

import (
	"github.com/helderfarias/go-api-kit/endpoint"
	"github.com/sirupsen/logrus"
)

type Subscriber interface {
	Delivery(receivers ...SubRec) error

	Subscribe(stream, subj, cons string, cmd endpoint.Endpoint) SubRec

	Close() error
}

type Publisher interface {
	Publish(stream, subj string, sender interface{}, o ...PubOpts) (interface{}, error)
}

type options struct {
	action endpoint.Endpoint
	args   map[string]interface{}
}

type SubRec func(o *options)

type PubOpts func(o *options)

type emptySub struct {
}

type emptyPub struct {
}

func (e *emptySub) Delivery(receivers ...SubRec) error {
	logrus.Warn("Subscriber not working...")
	return nil
}

func (e *emptySub) Close() error {
	logrus.Warn("Subscriber not working...")
	return nil
}

func (e *emptySub) Subscribe(stream, subj, cons string, cmd endpoint.Endpoint) SubRec {
	logrus.Warn("Subscriber not working...")
	return nil
}

func (e *emptyPub) Publish(stream, subj string, sender interface{}, o ...PubOpts) (interface{}, error) {
	logrus.Warn("Publisher not working...")
	return nil, nil
}
