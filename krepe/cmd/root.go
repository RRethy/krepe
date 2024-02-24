package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	pkgPath string
	rootCmd = &cobra.Command{
		Use:   "krepe",
		Short: "Kubernetes configuration management tooling.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			absPkgPath, err := filepath.Abs(pkgPath)
			if err != nil {
				panic(err)
			}

			err = os.Chdir(absPkgPath)
			if err != nil {
				panic(err)
			}
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&pkgPath, "pkgPath", "C", ".", "path to the package to run")
}
