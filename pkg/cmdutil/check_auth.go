package cmdutil

import (
	"github.com/spf13/cobra"
)

func DisableAuthCheck(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}

	cmd.Annotations["skipAuthCheck"] = "true"
}

// func CheckAuth(cfg config.Config) bool {
// 	hosts, err := cfg.Hosts()
// 	if err != nil {
// 		return false
// 	}
//
// 	for _, hostname := range hosts {
// 		token, _ := cfg.Get(hostname, "oauth_token")
// 		if token != "" {
// 			return true
// 		}
// 	}
//
// 	return false
// }

func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	if !cmd.Runnable() {
		return false
	}
	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["skipAuthCheck"] == "true" {
			return false
		}
	}

	return true
}