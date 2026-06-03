package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
)

var (
	pkgProjects string
	pkgMode     string
	pkgShell    string
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: i18n.T("Manage packages (groups of projects)"),
}

var pkgAddCmd = &cobra.Command{
	Use:   i18n.T("add <name>"),
	Short: i18n.T("Create a package with the given projects"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode, err := validateMode(pkgMode)
		if err != nil {
			return err
		}
		shell, err := validateShell(pkgShell)
		if err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		ids := splitCSV(pkgProjects)
		if len(ids) == 0 {
			return i18n.Errorf("provide --projects id1,id2,...")
		}
		for _, id := range ids {
			if _, ok := cfg.FindProject(id); !ok {
				return fmt.Errorf(i18n.T("project not found: %s"), id)
			}
		}
		pk, err := cfg.AddPackage(args[0], ids, mode, shell)
		if err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Package created: %s (%d projects)\n"), pk.ID, len(pk.ProjectIDs))
		return nil
	},
}

var pkgRemoveCmd = &cobra.Command{
	Use:     i18n.T("remove <pkgId>"),
	Aliases: []string{"rm"},
	Short:   i18n.T("Remove a package"),
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemovePackage(args[0]) {
			return fmt.Errorf(i18n.T("package not found: %s"), args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Package removed: %s\n"), args[0])
		return nil
	},
}

func splitCSV(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func init() {
	pkgAddCmd.Flags().StringVar(&pkgProjects, "projects", "", i18n.T("comma-separated project ids"))
	pkgAddCmd.Flags().StringVar(&pkgMode, "mode", "", i18n.T("package mode: tabs|windows"))
	pkgAddCmd.Flags().StringVar(&pkgShell, "shell", "", i18n.T("package shell: claude|pwsh|cmd"))
	pkgCmd.AddCommand(pkgAddCmd, pkgRemoveCmd)
	rootCmd.AddCommand(pkgCmd)
}
