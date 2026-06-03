package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
)

var (
	pkgProjects string
	pkgMode     string
	pkgShell    string
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Gerencia pacotes (grupos de projetos)",
}

var pkgAddCmd = &cobra.Command{
	Use:   "add <nome>",
	Short: "Cria um pacote com os projetos indicados",
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
			return fmt.Errorf("informe --projects id1,id2,...")
		}
		for _, id := range ids {
			if _, ok := cfg.FindProject(id); !ok {
				return fmt.Errorf("projeto não encontrado: %s", id)
			}
		}
		pk, err := cfg.AddPackage(args[0], ids, mode, shell)
		if err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Pacote criado: %s (%d projetos)\n", pk.ID, len(pk.ProjectIDs))
		return nil
	},
}

var pkgRemoveCmd = &cobra.Command{
	Use:     "remove <pkgId>",
	Aliases: []string{"rm"},
	Short:   "Remove um pacote",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemovePackage(args[0]) {
			return fmt.Errorf("pacote não encontrado: %s", args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Pacote removido: %s\n", args[0])
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
	pkgAddCmd.Flags().StringVar(&pkgProjects, "projects", "", "ids de projetos separados por vírgula")
	pkgAddCmd.Flags().StringVar(&pkgMode, "mode", "", "modo do pacote: tabs|windows")
	pkgAddCmd.Flags().StringVar(&pkgShell, "shell", "", "shell do pacote: claude|pwsh|cmd")
	pkgCmd.AddCommand(pkgAddCmd, pkgRemoveCmd)
	rootCmd.AddCommand(pkgCmd)
}
