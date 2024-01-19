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
  install [packageRef]

Arguments:
  packageRef  The reference to a package in the form 'github.com/$OWNER/$REPO[$PATH]@$TAG'. This argument is required.

Example:
  install 'github.com/Owner/Repo/path/to/package@v1.0.0'`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := install.Install(pkgPath, args[0], installPkgName)
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installPkgName, "name", "n", "", "name override of the package being installed")
}
