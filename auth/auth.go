package auth

import (
	"context"
	"github.com/Nerzal/gocloak/v7"
	"github.com/kaecloud/kaectl/internal/config"
	"strings"
)

type Token gocloak.JWT

func GetAccessToken(cfg *config.CmdConfig) (*Token, error){
	url := cfg.SSOHost
	if ! strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	client := gocloak.NewClient(url)
	ctx := context.Background()
	token, err := client.Login(ctx, cfg.SSOClientID, "", cfg.SSORealm, cfg.SSOUsername, cfg.SSOPassword)
	return (*Token)(token), err
}