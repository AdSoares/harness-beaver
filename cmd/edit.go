package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
)

var (
	editName  string
	editPath  string
	editShell string
)

var editCmd = &cobra.Command{
	Use:   i18n.T("edit <projId>"),
	Short: i18n.T("Change name, path or shell of a project"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell, err := validateShell(editShell)
		if err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		p, ok := cfg.FindProject(args[0])
		if !ok {
			return fmt.Errorf(i18n.T("project not found: %s"), args[0])
		}
		if cmd.Flags().Changed("name") {
			p.Name = editName
		}
		if cmd.Flags().Changed("path") {
			abs, err := filepath.Abs(editPath)
			if err != nil {
				return err
			}
			info, err := os.Stat(abs)
			if err != nil || !info.IsDir() {
				return fmt.Errorf(i18n.T("invalid path: %s"), editPath)
			}
			p.Path = abs
		}
		if cmd.Flags().Changed("shell") {
			p.DefaultShell = shell
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Project updated: %s\n"), p.ID)
		return nil
	},
}

func init() {
	editCmd.Flags().StringVar(&editName, "name", "", i18n.T("new name"))
	editCmd.Flags().StringVar(&editPath, "path", "", i18n.T("new path"))
	editCmd.Flags().StringVar(&editShell, "shell", "", i18n.T("new default shell: claude|pwsh|cmd"))
	rootCmd.AddCommand(editCmd)
}
