package run

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/kaecloud/kaectl/api"
	"github.com/kaecloud/kaectl/internal/config"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/kaecloud/kaectl/pkg/iostreams"
	"github.com/kaecloud/kaectl/pkg/spec"
	"github.com/kaecloud/kaectl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

type RunOptions struct {
	Config      func() (*config.CmdConfig, error)
	AccessToken func() (string, error)
	IO          *iostreams.IOStreams

	Command  string
	SpecFile string
	Cluster string
}

func NewCmdRun(f *cmdutil.Factory, runF func(*RunOptions) error) *cobra.Command {
	opts := &RunOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		AccessToken: f.GetAccessToken,
	}

	cmd := &cobra.Command{
		Use:   "run <command>",
		Short: "Run a new job",
		Long:  `Run a new k8s job.`,
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
	 		# run echo command in k8s
	 		$ kaectl job run "echo hello world"
	   `),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(
				`A command should be supplied as an argument.
            `),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Command = args[0]
			if runF != nil {
				return runF(opts)
			}

			return runRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Cluster, "cluster", "", "the cluster used to run command")
	cmd.Flags().StringVar(&opts.SpecFile, "spec", "job.yaml", "the spec file")

	return cmd
}

func runRun(opts *RunOptions) error {
	// var obj map[string]interface{}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	if opts.Cluster == "" {
		opts.Cluster = cfg.JobDefaultCluster
	}
	tok, err := opts.AccessToken()
	if err != nil {
		return err
	}
	c := api.NewJobClient(cfg.JobServerUrl, tok)

	data, err := ioutil.ReadFile(opts.SpecFile)
	if err != nil {
		return err
	}
	sp, err := spec.FromYAML(data)
	if err != nil {
		return err
	}
	// generate a new job name
	sp.Name = fmt.Sprintf("%s-%s", sp.Name, utils.RandStringRunes(6))
	err = cmdutil.PrepareJob(sp, c)
	if err != nil {
		return err
	}
	if len(sp.Containers) != 1 {
		return errors.Errorf("only one Container is allowed in run command")
	}

	cmdList := []string{"sh", "-c", opts.Command}
	sp.Containers[0].Command = cmdList
	sp.Containers[0].Name = sp.Name

	yamlBytes, err := spec.ToYAML(sp)
	if err != nil {
		return err
	}
	obj := &spec.CreateJobArgs{
		Spec: string(yamlBytes),
		Cluster: opts.Cluster,
	}

	job, err := c.Create(obj)
	if err != nil {
		return err
	}
	if sp.Cron != nil {
		jobUrl := fmt.Sprintf("%s/#/jobs/%s/detail?cluster=%s", strings.TrimRight(cfg.JobServerUrl, "/"), sp.Name, opts.Cluster)
		fmt.Printf("this is a cron job, opening url %s in browser.", jobUrl)
		err = utils.OpenInBrowser(jobUrl)
		return err
	}
	outCh, err := c.Logs(job.Name, "", opts.Cluster, true)
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
			fmt.Printf("%s", v)
		}
	}
}
