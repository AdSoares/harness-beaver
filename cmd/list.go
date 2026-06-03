package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
)

var listCmd = &cobra.Command{
	Use:       "list [projects|packages]",
	Short:     i18n.T("List registered projects and/or packages"),
	ValidArgs: []string{"projects", "packages"},
	Args:      cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		what := ""
		if len(args) == 1 {
			what = args[0]
		}
		if what == "" || what == "projects" {
			fmt.Println(i18n.T("Projects:"))
			if len(cfg.Projects) == 0 {
				fmt.Println(i18n.T("  (none)"))
			}
			for _, p := range cfg.Projects {
				shell := p.DefaultShell
				if shell == "" {
					shell = "-"
				}
				fmt.Printf("  %-20s %-28s [%s] %s\n", p.ID, p.Name, shell, p.Path)
			}
		}
		if what == "" || what == "packages" {
			fmt.Println(i18n.T("Packages:"))
			if len(cfg.Packages) == 0 {
				fmt.Println(i18n.T("  (none)"))
			}
			for _, pk := range cfg.Packages {
				fmt.Printf("  %-20s %-28s mode=%s shell=%s projects=[%s]\n",
					pk.ID, pk.Name, dash(string(pk.Mode)), dash(string(pk.Shell)),
					strings.Join(pk.ProjectIDs, ","))
			}
		}
		return nil
	},
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func init() {
	rootCmd.AddCommand(listCmd)
}
