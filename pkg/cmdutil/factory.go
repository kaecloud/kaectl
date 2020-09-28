package cmdutil

import (
	"github.com/kaecloud/kaectl/auth"
	"github.com/kaecloud/kaectl/internal/config"
	"github.com/kaecloud/kaectl/pkg/iostreams"
	"net/http"
)

type Factory struct {
	IOStreams  *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (*config.CmdConfig, error)
	GetAccessToken   func() (string, error)
}

func NewFactory(appVersion string) *Factory {
	io := iostreams.System()

	var cachedConfig *config.CmdConfig
	var configError error
	var tok *auth.Token
	configFunc := func() (*config.CmdConfig, error) {
		if cachedConfig != nil || configError != nil {
			return cachedConfig, configError
		}
		cachedConfig, configError = config.LoadCmdConfig()
		// if errors.Is(configError, os.ErrNotExist) {
		// 	cachedConfig = config.NewBlankConfig()
		// 	configError = nil
		// }
		return cachedConfig, configError
	}

	accessTokenFunc := func() (string, error) {
		if tok == nil {
			cfg, err := configFunc()
			if err != nil {
				return "", err
			}
			tok, err = auth.GetAccessToken(cfg)
			if err != nil {
				return "", err
			}
		}
		return tok.AccessToken, nil
	}
	return &Factory{
		IOStreams:      io,
		Config:         configFunc,
		GetAccessToken: accessTokenFunc,
	}
}

