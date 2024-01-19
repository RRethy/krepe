package cmd

import (
	"github.com/RRethy/krepe/krepe/pkg/install"
	"github.com/spf13/cobra"
)

var (
	installPkg  string
	installName string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		err := install.Install(installPkg, url, installName)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installName, "name", "n", "", "TODO")
	installCmd.Flags().StringVarP(&installPkg, "", "C", ".", "TODO")
}
