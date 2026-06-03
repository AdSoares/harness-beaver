package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Gerencia atalhos de comandos (use 'bvr do <nome>' para executar)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Aliases) == 0 {
			fmt.Println("Nenhum atalho. Crie com: bvr alias set <nome> <comando...>")
			return nil
		}
		names := make([]string, 0, len(cfg.Aliases))
		for n := range cfg.Aliases {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			fmt.Printf("  %-16s %s\n", n, cfg.Aliases[n])
		}
		return nil
	},
}

var aliasSetCmd = &cobra.Command{
	Use:   "set <nome> <comando...>",
	Short: "Cria ou atualiza um atalho",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		name := args[0]
		command := strings.Join(args[1:], " ")
		cfg.SetAlias(name, command)
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("%s = %s\n", name, command)
		return nil
	},
}

var aliasRmCmd = &cobra.Command{
	Use:     "rm <nome>",
	Aliases: []string{"remove"},
	Short:   "Remove um atalho",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemoveAlias(args[0]) {
			return fmt.Errorf("atalho não encontrado: %s", args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Atalho removido: %s\n", args[0])
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasSetCmd, aliasRmCmd)
	rootCmd.AddCommand(aliasCmd)
}
