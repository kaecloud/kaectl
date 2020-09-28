package main

import (
	"errors"
	"fmt"
	"github.com/kaecloud/kaectl/api"
	"github.com/kaecloud/kaectl/pkg/cmd/root"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/kaecloud/kaectl/version"
	"github.com/spf13/cobra"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	hasDebug := os.Getenv("DEBUG") != ""

	cmdFactory := cmdutil.NewFactory(version.Version)
	stderr := cmdFactory.IOStreams.ErrOut
	rootCmd := root.NewCmdRoot(cmdFactory, version.Version, version.BuildDate)

	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	// cmd, _, err := rootCmd.Traverse(expandedArgs)

	// authCheckEnabled := cmdutil.IsAuthCheckEnabled(cmd)

	// if authCheckEnabled {
	// 	hasAuth := false

	// 	cfg, err := cmdFactory.Config()
	// 	if err == nil {
	// 		hasAuth = cmdutil.CheckAuth(cfg)
	// 	}
	// }

	rootCmd.SetArgs(expandedArgs)

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		printError(stderr, err, cmd, hasDebug)

		var httpErr api.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == 401 {
			fmt.Println("hint: try authenticating with `gh auth login`")
		}

		os.Exit(1)
	}
	if root.HasFailed() {
		os.Exit(1)
	}

}

func printError(out io.Writer, err error, cmd *cobra.Command, debug bool) {
	if err == cmdutil.SilentError {
		return
	}

	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		fmt.Fprintf(out, "error connecting to %s\n", dnsError.Name)
		if debug {
			fmt.Fprintln(out, dnsError)
		}
		fmt.Fprintln(out, "check your internet connection or githubstatus.com")
		return
	}

	fmt.Fprintln(out, err)

	var flagError *cmdutil.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, cmd.UsageString())
	}
}

func isCompletionCommand() bool {
	return len(os.Args) > 1 && os.Args[1] == "completion"
}

