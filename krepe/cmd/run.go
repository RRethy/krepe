package cmd

import (
	"github.com/RRethy/krepe/krepe/pkg/run"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a pipeline on a package",
		Long: `Run a pipeline on a package.

Usage:
  run [pipelineName]

Arguments:
  pipelineName  The name of the pipeline to run. This argument is optional. Default is 'default'.`,
		Run: func(cmd *cobra.Command, args []string) {
			name := "default"
			if len(args) > 0 {
				name = args[0]
			}

			err := run.Run(pkgPath, name)
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}
