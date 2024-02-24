package cmd

import (
	"github.com/RRethy/krepe/krepe/pkg/run"
	"github.com/spf13/cobra"
)

const (
	DEFAULT_PIPELINE_NAME = "default"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a pipeline on a package",
		Long: `Run a pipeline on a package.

Usage:
  krepe run [pipelineName]

Arguments:
  pipelineName  The name of the pipeline to run. This argument is optional. Default is 'default'.

Example:
  krepe run default`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := DEFAULT_PIPELINE_NAME
			if len(args) > 0 {
				name = args[0]
			}

			err := run.Run(name)
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}
