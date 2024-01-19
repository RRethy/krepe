package cmd

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/run"
	"github.com/spf13/cobra"
)

var pipeline string
var function string
var pkg string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := run.Run(pkg, pipeline, function)
		if err != nil {
			panic(err)
		}
		fmt.Println("run succeeded")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	runCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "default", "TODO")
	runCmd.Flags().StringVarP(&function, "function", "f", "", "TODO")
	runCmd.Flags().StringVar(&pkg, "pkg", ".", "TODO")
}
