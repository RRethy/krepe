package cmd

import (
	"github.com/Shopify/krepe/krepe/pkg/run"
	"github.com/spf13/cobra"
)

var pipeline string
var function string
var pkg string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a pipeline or function on a package",
	Run: func(cmd *cobra.Command, args []string) {
		err := run.Run(pkg, pipeline, function)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "default", "name of the pipeline to run")
	runCmd.Flags().StringVarP(&function, "function", "f", "", "name of the function to run")
	runCmd.Flags().StringVar(&pkg, "pkg", ".", "path to the package to run")
}
