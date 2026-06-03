package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
)

var removeCmd = &cobra.Command{
	Use:     i18n.T("remove <projId>"),
	Aliases: []string{"rm"},
	Short:   i18n.T("Remove a project from the registry"),
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemoveProject(args[0]) {
			return fmt.Errorf(i18n.T("project not found: %s"), args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Project removed: %s\n"), args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
