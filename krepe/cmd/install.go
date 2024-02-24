package cmd

import (
	"github.com/RRethy/krepe/krepe/pkg/install"
	"github.com/spf13/cobra"
)

var (
	installPkgName string
	installCmd     = &cobra.Command{
		Use:   "install",
		Short: "Import a package",
		Long: `Import a pacakge.

Usage:
  krepe install [path]

Arguments:
  path  The path to the package to install. Relative paths must be relative to the package being operated on. This argument is required.

Example:
  krepe install ../some_package`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := install.Install(args[0], installPkgName)
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installPkgName, "name", "n", "", "name override of the package being installed")
}
