package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/runner"
)

var doShell string

var doCmd = &cobra.Command{
	Use:   "do <nome> [args...]",
	Short: "Executa um atalho (alias) registrado, expandindo {1}..{N} e {*}",
	Long: `Executa o comando do atalho <nome>. Os argumentos extras substituem os
placeholders {1}, {2}, … e {*} (todos juntos); sem placeholders, são anexados.

Exemplos:
  bvr alias set deploy "npm run build && npm run deploy"
  bvr do deploy
  bvr alias set commit "git commit -m {1}"
  bvr do commit "fix: ajuste"`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		tmpl, ok := cfg.Alias(args[0])
		if !ok {
			return fmt.Errorf("atalho não encontrado: %s (veja 'bvr alias')", args[0])
		}
		shell := cfg.Settings.DefaultRunShell
		if doShell != "" {
			s, err := validateRunShell(doShell)
			if err != nil {
				return err
			}
			shell = s
		}
		command := config.ExpandAlias(tmpl, args[1:])
		fmt.Fprintf(os.Stderr, "» %s\n", command)
		res, err := runner.Run(shell, command, "", true)
		if err != nil {
			return err
		}
		os.Exit(res.ExitCode)
		return nil
	},
}

func init() {
	doCmd.Flags().SetInterspersed(false)
	doCmd.Flags().StringVar(&doShell, "shell", "", "shell de execução: pwsh|cmd|bash|zsh")
	rootCmd.AddCommand(doCmd)
}
