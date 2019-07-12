package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestLoadConfig(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8888").
		Get("/app/service/master").
		Reply(200).
		BodyString(`
			{
				"name":"accountservice-test",
				"profiles":["test"],
				"label":null,
				"version":null,
				"propertySources":[
					{"index":2,"name":"./service.yml","source":{"server_port":6767}},
					{"index":1,"name":"https://github.com/service","source":{"server_port":0,"server_name":"Accountservice"}},
					{"index":0,"name":"vault","source":{"server_port":10,"database":"postgres"}}
				]
			}`,
		)

	s := NewConfigServer("app", "service", "http://localhost:8888", "master", "token00")

	sources := map[string]interface{}{}

	err := s.Load(func(key string, value interface{}) {
		sources[key] = value
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, len(sources))
	assert.Equal(t, float64(6767), sources["server_port"])
	assert.Equal(t, "Accountservice", sources["server_name"])
	assert.Equal(t, "postgres", sources["database"])
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
