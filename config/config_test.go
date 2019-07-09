package config

import (
	"testing"

	"gopkg.in/h2non/gock.v1"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8888").
		Get("/app/service/master").
		Reply(200).
		BodyString(`{"name":"accountservice-test","profiles":["test"],"label":null,"version":null,"propertySources":[{"name":"file:/config-repo/accountservice-test.yml","source":{"server_port":6767,"server_name":"Accountservice"}}]}`)

	s := NewConfigServer("app", "service", "http://localhost:8888", "master", "token00")

	sources := map[string]interface{}{}

	err := s.Load(func(key string, value interface{}) {
		sources[key] = value
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(sources))
	assert.Equal(t, float64(6767), sources["server_port"])
	assert.Equal(t, "Accountservice", sources["server_name"])
}

func TestLoadConfigWithoutEnvSet(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8888").
		Get("/app/service/master").
		Reply(200).
		BodyString(`{"name":"accountservice-test","profiles":["test"],"label":null,"version":null,"propertySources":[{"name":"file:/config-repo/accountservice-test.yml","source":{"server_port":6767,"server_name":"Accountservice"}}]}`)

	s := NewConfigServer("app", "service", "http://localhost:8888", "master", "token00")

	err := s.Load(nil)

	assert.NoError(t, err)
}
