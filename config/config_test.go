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
					{"index":2,"name":"vault","source":{"server_port":10,"database":"postgres"}},				
					{"index":1,"name":"./service.yml","source":{"server_port":6767}},
					{"index":0,"name":"https://github.com/service","source":{"server_port":0,"server_name":"Accountservice"}}
				]
			}`,
		)

	s := NewConfigServer(
		App("app"),
		Profile("service"),
		Server("http://localhost:8888"),
		Label("master"),
		Token("token00"),
	)

	sources := map[string]interface{}{}

	err := s.Load(func(key string, value interface{}) {
		sources[key] = value
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, len(sources))
	assert.Equal(t, float64(10), sources["server_port"])
	assert.Equal(t, "Accountservice", sources["server_name"])
	assert.Equal(t, "postgres", sources["database"])
}

func TestLoadConfigWithoutEnvSet(t *testing.T) {
	defer gock.Off()
	gock.New("http://localhost:8888").
		Get("/app/service/master").
		Reply(200).
		BodyString(`{"name":"accountservice-test","profiles":["test"],"label":null,"version":null,"propertySources":[{"name":"file:/config-repo/accountservice-test.yml","source":{"server_port":6767,"server_name":"Accountservice"}}]}`)

	s := NewConfigServer(
		App("app"),
		Profile("service"),
		Server("http://localhost:8888"),
		Label("master"),
		Token("token00"),
		VaultToken("vault_token_00"))

	err := s.Load(nil)

	assert.NoError(t, err)
}

func TestLoadConfigFromLocalYamlFile(t *testing.T) {
	s := NewConfigServer(LocalYamlFile("testdata/file.yml"))

	sources := map[string]interface{}{}

	err := s.Load(func(key string, value interface{}) {
		sources[key] = value
	})

	assert.Nil(t, err)
	assert.Equal(t, 5, len(sources))
	assert.Equal(t, sources["rest_server_host"], "http://localhost")
	assert.Equal(t, sources["grpc_server_host"], "grpc://localhost")
	assert.Equal(t, sources["rest_server_port"], 4002)
	assert.Equal(t, sources["grpc_server_port"], 50051)
	assert.Equal(t, sources["gateway_url"], "http://localhost:3010")
}

func TestShouldErrorFileIfNotExistsWhenLoadConfigFromLocalYamlFile(t *testing.T) {
	s := NewConfigServer(LocalYamlFile("testdata/file_is_not_exists.yml"))

	sources := map[string]interface{}{}

	err := s.Load(func(key string, value interface{}) {
		sources[key] = value
	})

	assert.EqualError(t, err, "open testdata/file_is_not_exists.yml: no such file or directory")
}

func TestLoadConfigExternalFile(t *testing.T) {
	s := NewConfigServer(LocalYamlFile("testdata/file_embbed.yml"))

	sources := map[string]interface{}{}

	err := s.Load(func(key string, value interface{}) {
		sources[key] = value
	})

	assert.Nil(t, err)
	assert.Equal(t, 5, len(sources))
	assert.Equal(t, sources["rest_server_host"], "http://localhost")
	assert.Equal(t, sources["grpc_server_host"], "grpc://localhost")
	assert.Equal(t, sources["rest_server_port"], 4002)
	assert.Equal(t, sources["grpc_server_port"], 50051)
	assert.Equal(t, sources["gateway_url"], "http://localhost:3010")
}
