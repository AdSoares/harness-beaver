package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/i18n"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: i18n.T("Show the bvr version"),
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("bvr %s (%s/%s, %s)\n", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
