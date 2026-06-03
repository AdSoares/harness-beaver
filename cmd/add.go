package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
)

var addShell string

var addCmd = &cobra.Command{
	Use:   i18n.T("add <name> <path>"),
	Short: i18n.T("Register a new project (directory)"),
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell, err := validateShell(addShell)
		if err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		p, err := cfg.AddProject(args[0], args[1], shell)
		if err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Project added: %s (%s)\n"), p.ID, p.Path)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addShell, "shell", "", i18n.T("project default shell: claude|pwsh|cmd"))
	rootCmd.AddCommand(addCmd)
}
