package cmd

import "github.com/spf13/cobra"

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Abre a interface interativa (TUI)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
