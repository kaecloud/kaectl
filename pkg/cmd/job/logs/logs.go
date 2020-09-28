package logs

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/kaecloud/kaectl/api"
	"github.com/kaecloud/kaectl/internal/config"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/kaecloud/kaectl/pkg/iostreams"
	"github.com/spf13/cobra"
)

type LogsOptions struct {
	Config     func() (*config.CmdConfig, error)
	AccessToken func() (string, error)
	IO         *iostreams.IOStreams

	Name string
	Cluster string
	Follow bool
}

func NewCmdLogs(f *cmdutil.Factory, runF func(*LogsOptions) error) *cobra.Command{
	opts := &LogsOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		AccessToken: f.GetAccessToken,
	}

	cmd := &cobra.Command{
		Use:   "logs [<name>]",
		Short: "logs a job",
		Long:  `logs job by name.`,
		Example: heredoc.Doc(`
	 		# get job's log
	 		$ kaectl job logs my-job --follow
	   `),
	    Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}

			if runF != nil {
				return runF(opts)
			}

			return logsRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Cluster, "cluster", "c", "", "cluster")
	cmd.Flags().BoolVar(&opts.Follow, "follow", false, "follow the pod log")

	return cmd
}

func logsRun(opts *LogsOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	tok, err := opts.AccessToken()
	if err != nil {
		return err
	}
	c := api.NewJobClient(cfg.JobServerUrl, tok)
	outCh, err := c.Logs(opts.Name, "", opts.Cluster, opts.Follow)
	if err != nil {
		return err
	}
	for {
		item, ok := <- outCh
		if !ok {
			return nil
		}
		switch v := item.(type) {
		case error:
			return v
		case string:
			fmt.Printf("%s\n", v)
		}
	}
}
