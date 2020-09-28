package create

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/kaecloud/kaectl/api"
	"github.com/kaecloud/kaectl/internal/config"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/kaecloud/kaectl/pkg/iostreams"
	"github.com/kaecloud/kaectl/pkg/spec"
	"github.com/kaecloud/kaectl/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type CreateOptions struct {
	Config     func() (*config.CmdConfig, error)
	AccessToken func() (string, error)
	IO         *iostreams.IOStreams

	Name string
	Image string
	Command string
	Shell bool
	SpecFile string
	Cluster string
}

func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command{
	opts := &CreateOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		AccessToken: f.GetAccessToken,
	}

	cmd := &cobra.Command{
	 	Use:   "create [<name>]",
	 	Short: "Create a new job",
	 	Long:  `Create a new k8s job.`,
	 	Args:  cobra.MaximumNArgs(1),
	 	Example: heredoc.Doc(`
	 		# create a job with a specific name
	 		$ kaectl job create my-job --image ubuntu:18.04 --command "echo hello world" --cluster mycluster
	   `),
	 	Annotations: map[string]string{
	 		"help:arguments": heredoc.Doc(
	 			`A name should be supplied as an argument
            "`),
	 	},
	 	RunE: func(cmd *cobra.Command, args []string) error {
	 		if len(args) > 0 {
	 			opts.Name = args[0]
	 		}

	 		if runF != nil {
	 			return runF(opts)
	 		}

	 		return createRun(opts)
	 	},
	}

	// cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Job name")
	cmd.Flags().StringVarP(&opts.Image, "image", "i", "", "Job's docker image")
	cmd.Flags().StringVarP(&opts.Command, "command", "c", "", "the command needs to run")
	cmd.Flags().BoolVar(&opts.Shell, "shell", true, "Use shell to run the command")
	cmd.Flags().StringVar(&opts.SpecFile, "spec", "job.yaml", "the spec file")
	cmd.Flags().StringVar(&opts.Cluster, "cluster", "", "cluster name")

	return cmd
}

func createRun(opts *CreateOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	tok, err := opts.AccessToken()
	if err != nil {
		return err
	}
	c := api.NewJobClient(cfg.JobServerUrl, tok)

	var obj *spec.CreateJobArgs

	if utils.FileExists(opts.SpecFile) {
		data, err := ioutil.ReadFile(opts.SpecFile)
		if err != nil {
			return err
		}
		sp, err := spec.FromYAML(data)
		if err != nil {
			return err
		}
		err = cmdutil.PrepareJob(sp, c)
		if err != nil {
			return err
		}
		obj = &spec.CreateJobArgs{
			Spec: string(data),
			Cluster: opts.Cluster,
		}
	} else {
		obj = &spec.CreateJobArgs{
			Name: opts.Name,
			Image: opts.Image,
			Shell: opts.Shell,
			Command: opts.Command,
			Cluster: opts.Cluster,
		}
	}

	job, err := c.Create(obj)
	if err != nil {
		return err
	}

	fmt.Printf("Create job %s successfully\n", job.Name)
	return nil
}
