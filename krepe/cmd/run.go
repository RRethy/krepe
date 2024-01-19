package cmd

import (
	"github.com/RRethy/krepe/krepe/pkg/run"
	"github.com/spf13/cobra"
)

var (
	runPkg string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a pipeline or function on a package",
	Run: func(cmd *cobra.Command, args []string) {
		name := "default"
		if len(args) > 0 {
			name = args[0]
		}

		err := run.Run(runPkg, name)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&runPkg, "", "C", ".", "path to the package to run")
}
