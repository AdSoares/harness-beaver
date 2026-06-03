package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
)

var listCmd = &cobra.Command{
	Use:       "list [projects|packages]",
	Short:     "Lista projetos e/ou pacotes registrados",
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
			fmt.Println("Projetos:")
			if len(cfg.Projects) == 0 {
				fmt.Println("  (nenhum)")
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
			fmt.Println("Pacotes:")
			if len(cfg.Packages) == 0 {
				fmt.Println("  (nenhum)")
			}
			for _, pk := range cfg.Packages {
				fmt.Printf("  %-20s %-28s mode=%s shell=%s projetos=[%s]\n",
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
