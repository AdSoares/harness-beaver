package cmd

import (
	"github.com/spf13/cobra"

	"harnessbeaver/internal/i18n"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: i18n.T("Open the interactive interface (TUI)"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
