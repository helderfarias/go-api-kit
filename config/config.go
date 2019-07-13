package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type ConfigServer interface {
	Load(env EnvSet) error
}

type EnvSet func(key string, value interface{})

func Nop(key string, value interface{}) {}

type configServer struct {
	application string
	profile     string
	server      string
	label       string
	token       string
}

type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Index  int                    `json:"index"`
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

func NewConfigServer(app, profile, server, label, token string) ConfigServer {
	return &configServer{
		application: app,
		profile:     profile,
		server:      server,
		label:       label,
		token:       token,
	}
}

// Load config from file to viper
func (s *configServer) Load(env EnvSet) error {
	url := fmt.Sprintf("%s/%s/%s/%s?apikey=%s", s.server, s.application, s.profile, s.label, s.token)
	logrus.Infof("Loading config from %s", url)

	body, err := s.fetch(url)
	if err != nil {
		return fmt.Errorf("Couldn't load configuration, cannot start. Terminating. Error:%s", err.Error())
	}

	cmd := Nop
	if env != nil {
		cmd = env
	}

	return s.parse(body, cmd)
}

func (s *configServer) fetch(url string) ([]byte, error) {
	client := resty.New()

	resp, err := client.
		SetDebug(false).
		SetDisableWarn(true).
		SetRetryCount(5).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("Couldn't load configuration, cannot start. Terminating. Error: %s", err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Couldn't load configuration, parse error: %s", resp.String())
	}

	return resp.Body(), nil
}

func (s *configServer) parse(body []byte, env EnvSet) error {
	var cloudConfig springCloudConfig

	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		return fmt.Errorf("Cannot parse configuration, message: %s", err.Error())
	}

	sort.SliceStable(cloudConfig.PropertySources, func(a, b int) bool {
		return cloudConfig.PropertySources[a].Index > cloudConfig.PropertySources[b].Index
	})

	for i := len(cloudConfig.PropertySources) - 1; i >= 0; i-- {
		props := cloudConfig.PropertySources[i]

		for key, value := range props.Source {
			env(key, value)
			logrus.Debugf("Loading config property %v => %v\n", key, value)
		}
	}

	return nil
}
