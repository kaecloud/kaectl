package job

import (
	"github.com/MakeNowJust/heredoc"
	jobCreateCmd "github.com/kaecloud/kaectl/pkg/cmd/job/create"
	jobGetCmd "github.com/kaecloud/kaectl/pkg/cmd/job/get"
	jobDeleteCmd "github.com/kaecloud/kaectl/pkg/cmd/job/delete"
	jobRunCmd "github.com/kaecloud/kaectl/pkg/cmd/job/run"
	jobLogsCmd "github.com/kaecloud/kaectl/pkg/cmd/job/logs"
	"github.com/kaecloud/kaectl/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdJob(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job <command>",
		Short: "Create, Get, List, Delete, and Run job",
		Long:  `Work with jobs`,
		Example: heredoc.Doc(`
			$ kaectl job create my-job --image ubuntu:16.04 --command "echo hello" --cluster mycluster
			$ kaectl job get my-job
			$ kaectl job delete my-job
		`),
		Annotations: map[string]string{
			"IsCore": "true",
			"help:arguments": heredoc.Doc(`
				A repository can be supplied as an argument in any of the following formats:
			"
			`),
		},
	}

	cmd.AddCommand(jobCreateCmd.NewCmdCreate(f, nil))
	cmd.AddCommand(jobGetCmd.NewCmdGet(f, nil))
	cmd.AddCommand(jobDeleteCmd.NewCmdDelete(f, nil))
	cmd.AddCommand(jobRunCmd.NewCmdRun(f, nil))
	cmd.AddCommand(jobLogsCmd.NewCmdLogs(f, nil))

	return cmd
}