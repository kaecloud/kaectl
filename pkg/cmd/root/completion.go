package root

import (
	"errors"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/kaecloud/kaectl/pkg/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdCompletion(io *iostreams.IOStreams) *cobra.Command {
	var shellType string

	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: heredoc.Doc(`
			Generate shell completion scripts for kaectl commands.
			The output of this command will be computer code and is meant to be saved to a
			file or immediately evaluated by an interactive shell.
			For example, for bash you could add this to your '~/.bash_profile':
				eval "$(kaectl completion -s bash)"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if shellType == "" {
				if io.IsStdoutTTY() {
					return &cmdutil.FlagError{Err: errors.New("error: the value for `--shell` is required")}
				}
				shellType = "bash"
			}

			w := io.Out
			rootCmd := cmd.Parent()

			switch shellType {
			case "bash":
				return rootCmd.GenBashCompletion(w)
			case "zsh":
				return rootCmd.GenZshCompletion(w)
			case "powershell":
				return rootCmd.GenPowerShellCompletion(w)
			case "fish":
				return rootCmd.GenFishCompletion(w, true)
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.Flags().StringVarP(&shellType, "shell", "s", "", "Shell type: {bash|zsh|fish|powershell}")

	return cmd
}