package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	pkgPath string
	rootCmd = &cobra.Command{
		Use:   "krepe",
		Short: "Kubernetes configuration management tooling",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&pkgPath, "pkg", "C", ".", "path to the package to run")
}
