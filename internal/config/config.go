package config

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/kaecloud/kaectl/utils"
	"io/ioutil"
	"strings"
)

type CmdConfig struct {
	SSOHost      string `json:"sso_host" yaml:"sso_host"`
	SSOUsername  string `json:"sso_username" yaml:"sso_username"`
	SSOPassword  string `json:"sso_password" yaml:"sso_password"`
	SSORealm     string `json:"sso_realm" yaml:"sso_realm"`
	SSOClientID  string `json:"sso_client_id" yaml:"sso_client_id"`
	JobServerUrl string `json:"job_server_url" yaml:"job_server_url"`
	AppServerUrl string `json:"app_server_url" yaml:"app_server_url"`
	JobDefaultCluster string `json:"job_default_cluster" yaml:"job_default_cluster"`
}

func LoadCmdConfig(opts ...string) (*CmdConfig, error) {
	var (
		cfg      CmdConfig
		filename  = "~/.kae/config.yaml"
	)
	if len(opts) == 1 {
		filename = opts[0]
	}
	filename = utils.ExpandUser(filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		err = yaml.Unmarshal(data, &cfg)
	} else {
		err = json.Unmarshal(data, &cfg)
	}
	return &cfg, err
}