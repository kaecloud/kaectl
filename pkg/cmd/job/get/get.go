package get

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

type GetOptions struct {
	Config     func() (*config.CmdConfig, error)
	AccessToken func() (string, error)
	IO         *iostreams.IOStreams

	Name string
}

func NewCmdGet(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command{
	opts := &GetOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		AccessToken: f.GetAccessToken,
	}

	cmd := &cobra.Command{
		Use:   "get [<name>]",
		Short: "get a job",
		Long:  `get job by name.`,
		Example: heredoc.Doc(`
	 		# get job with specific name
	 		$ kaectl job get my-job
	   `),
	    Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}

			if runF != nil {
				return runF(opts)
			}

			// if opts.Template != "" && (opts.Homepage != "" || opts.Team != "" || !opts.EnableIssues || !opts.EnableWiki) {
			// 	return &cmdutil.FlagError{Err: errors.New(`The '--template' option is not supported with '--homepage, --team, --enable-issues or --enable-wiki'`)}
			// }

			return getRun(opts)
		},
	}

	// cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Description of repository")
	// cmd.Flags().StringVarP(&opts.Homepage, "homepage", "h", "", "Repository home page URL")

	return cmd
}

func getRun(opts *GetOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	tok, err := opts.AccessToken()
	if err != nil {
		return err
	}
	c := api.NewJobClient(cfg.JobServerUrl, tok)
	job, err := c.Get(opts.Name)
	if err != nil {
		return err
	}
	fmt.Printf("%s:\n  %s\n", utils.Bold("Name"), job.Name)
	fmt.Printf("%s:\n  %s\n", utils.Bold("Comment"), job.Comment)
	fmt.Printf("%s:\n  %s\n", utils.Bold("Spec"), job.SpecText)
	return nil
}
