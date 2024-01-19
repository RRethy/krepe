package cmd

import (
	"github.com/Shopify/krepe/krepe/pkg/run"
	"github.com/spf13/cobra"
)

var (
	runPipeline string
	runFunction string
	runPkg      string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a pipeline or function on a package",
	Run: func(cmd *cobra.Command, args []string) {
		err := run.Run(runPkg, runPipeline, runFunction)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&runPipeline, "pipeline", "p", "default", "name of the pipeline to run")
	runCmd.Flags().StringVarP(&runFunction, "function", "f", "", "name of the function to run")
	runCmd.Flags().StringVar(&runPkg, "pkg", ".", "path to the package to run")
}
