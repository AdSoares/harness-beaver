package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostra a versão do bvr",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("bvr %s (%s/%s, %s)\n", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
