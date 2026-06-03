package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
)

var addShell string

var addCmd = &cobra.Command{
	Use:   "add <nome> <path>",
	Short: "Registra um novo projeto (diretório)",
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
		fmt.Printf("Projeto adicionado: %s (%s)\n", p.ID, p.Path)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addShell, "shell", "", "shell padrão do projeto: claude|pwsh|cmd")
	rootCmd.AddCommand(addCmd)
}
