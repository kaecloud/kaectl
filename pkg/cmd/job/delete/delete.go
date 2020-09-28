package delete

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/kaecloud/kaectl/api"
	"github.com/kaecloud/kaectl/internal/config"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/kaecloud/kaectl/pkg/iostreams"
	"github.com/kaecloud/kaectl/utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	Config     func() (*config.CmdConfig, error)
	AccessToken func() (string, error)
	IO         *iostreams.IOStreams

	Name string
}

func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command{
	opts := &DeleteOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		AccessToken: f.GetAccessToken,
	}

	cmd := &cobra.Command{
		Use:   "delete [<name>]",
		Short: "delete a job",
		Long:  `delete job by name.`,
		Example: heredoc.Doc(`
	 		# delete job with specific name
	 		$ kaectl job delete my-job
	   `),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}

			if runF != nil {
				return runF(opts)
			}

			return deleteRun(opts)
		},
	}

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	tok, err := opts.AccessToken()
	if err != nil {
		return err
	}
	c := api.NewJobClient(cfg.JobServerUrl, tok)
	err = c.Delete(opts.Name)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", utils.Green(fmt.Sprintf("Delete job %s successfully", opts.Name)))
	return nil
}
