package cmd

import (
	"github.com/Shopify/krepe/krepe/pkg/install"
	"github.com/spf13/cobra"
)

var (
	installUrl  string
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
	Run: func(cmd *cobra.Command, args []string) {
		err := install.Install(installPkg, installUrl, installName)
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
