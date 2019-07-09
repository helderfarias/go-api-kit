package discovery

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type consulClient struct {
	client       *api.Client
	registration api.AgentServiceRegistration
}

type consulClientFallback struct {
}

type ConsulClientSettigs struct {
	ConsulAddress    string
	ConsulPort       string
	AdvertiseAddress string
	AdvertisePort    string
	ServiceName      string
	Tags             []string
	Endpoint         string
}

// NewConsulRegister method.
func NewConsulRegister(settings *ConsulClientSettigs) ServiceDiscoveryRegister {
	rand.Seed(time.Now().UTC().UnixNano())

	consulConfig := api.DefaultConfig()
	consulConfig.Address = settings.ConsulAddress + ":" + settings.ConsulPort
	cli, err := api.NewClient(consulConfig)
	if err != nil {
		return &consulClientFallback{}
	}

	check := api.AgentServiceCheck{
		HTTP:     "http://" + settings.AdvertiseAddress + ":" + settings.AdvertisePort + settings.Endpoint,
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Basic health checks",
	}

	port, _ := strconv.Atoi(settings.AdvertisePort)
	num := rand.Intn(100)
	asr := api.AgentServiceRegistration{
		ID:      settings.ServiceName + strconv.Itoa(num),
		Name:    settings.ServiceName,
		Address: settings.AdvertiseAddress,
		Port:    port,
		Tags:    settings.Tags,
		Check:   &check,
	}

	return &consulClient{client: cli, registration: asr}
}

func (c *consulClient) Register() {
	if err := c.client.Agent().ServiceRegister(&c.registration); err != nil {
		logrus.Error(err)
	}
}

func (c *consulClient) UnRegister() {
	if err := c.client.Agent().ServiceDeregister(c.registration.ID); err != nil {
		logrus.Error(err)
	}
}

func (*consulClientFallback) Register() {
	logrus.Warn("[Consult-Register] fallback")
}

func (*consulClientFallback) UnRegister() {
	logrus.Warn("[Consult-UnRegister] fallback")
}
