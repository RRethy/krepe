package cmd

import (
	"github.com/RRethy/krepe/krepe/pkg/update"
	"github.com/spf13/cobra"
)

var (
	updatePkgName string
	updateCmd     = &cobra.Command{
		Use:   "update",
		Short: "Update a package",
		Long: `Update a package.

Usage:
  update [packageRef]

Arguments:
  packageRef  The reference to a package in the form 'github.com/$OWNER/$REPO[$PATH]@$TAG'. This argument is required.

Example:
  update 'github.com/Owner/Repo/path/to/package@v1.0.0'`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := update.Update(pkgPath, args[0], installPkgName)
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updatePkgName, "name", "n", "", "name override of the package being updated")
}
