package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
)

var removeCmd = &cobra.Command{
	Use:     "remove <projId>",
	Aliases: []string{"rm"},
	Short:   "Remove um projeto do registro",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemoveProject(args[0]) {
			return fmt.Errorf("projeto não encontrado: %s", args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Projeto removido: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
