package root

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
	jobCmd "github.com/kaecloud/kaectl/pkg/cmd/job"
)

func NewCmdRoot(f *cmdutil.Factory, version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kaectl <command> <subcommand> [flags]",
		Short: "KAE CLI",
		Long:  `Work seamlessly with KAE from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ kaectl job create
			$ kaectl job get
		`),
		Annotations: map[string]string{
			"help:feedback": heredoc.Doc(`
				Open an issue using “gh issue create -R cli/cli”
			`),
			"help:environment": heredoc.Doc(`
				GITHUB_TOKEN: an authentication token for API requests. Setting this avoids being
				prompted to authenticate and overrides any previously stored credentials.
	
				GH_REPO: specify the GitHub repository in the "[HOST/]OWNER/REPO" format for commands
				that otherwise operate on a local repository.
				GH_HOST: specify the GitHub hostname for commands that would otherwise assume
				the "github.com" host when not in a context of an existing repository.
	
				BROWSER: the web browser to use for opening links.
	
				DEBUG: set to any value to enable verbose output to standard error. Include values "api"
				or "oauth" to print detailed information about HTTP requests or authentication flow.
	
				GLAMOUR_STYLE: the style to use for rendering Markdown. See
				https://github.com/charmbracelet/glamour#styles
	
				NO_COLOR: avoid printing ANSI escape sequences for color output.
			`),
		},
	}

	version = strings.TrimPrefix(version, "v")
	if buildDate == "" {
		cmd.Version = version
	} else {
		cmd.Version = fmt.Sprintf("%s (%s)", version, buildDate)
	}
	versionOutput := fmt.Sprintf("kaectl version %s\n", cmd.Version)
	cmd.AddCommand(&cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(versionOutput)
		},
	})
	cmd.SetVersionTemplate(versionOutput)
	cmd.Flags().Bool("version", false, "Show kaectl version")

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetHelpFunc(rootHelpFunc)
	cmd.SetUsageFunc(rootUsageFunc)

	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if err == pflag.ErrHelp {
			return err
		}
		return &cmdutil.FlagError{Err: err}
	})

	cmdutil.DisableAuthCheck(cmd)

	// CHILD COMMANDS

	cmd.AddCommand(jobCmd.NewCmdJob(f))
	cmd.AddCommand(NewCmdCompletion(f.IOStreams))

	return cmd
}

