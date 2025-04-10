package cmd

import (
	"errors"

	"github.com/RRethy/krepe/krepe/pkg/update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a package",
	Long: `Update a package.

Usage:
  krepe update packageImport

Arguments:
  packageImport  The name of the imported package to update. This argument is required.

Example:
  krepe update some_imported_package`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("packageName is required")
		}

		err := update.Update(args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	// rootCmd.AddCommand(updateCmd)
}
