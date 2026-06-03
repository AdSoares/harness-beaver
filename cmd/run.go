package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/journal"
	"harnessbeaver/internal/runner"
)

var runShell string

var runCmd = &cobra.Command{
	Use:   "run <comando...>",
	Short: "Executa um comando (pwsh/cmd) e registra no diário",
	Long: `Executa um comando em nome do usuário no PowerShell (default) ou cmd e
registra comando + saída no diário de bordo (~/.harnessbeaver/logs).

Flags devem vir antes do comando. Exemplos:
  bvr run echo ola
  bvr run --shell cmd "echo oi & dir"`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		shell := cfg.Settings.DefaultRunShell
		if runShell != "" {
			s, err := validateRunShell(runShell)
			if err != nil {
				return err
			}
			shell = s
		}
		printPendingHint()

		command := strings.Join(args, " ")
		res, err := runner.Run(shell, command, "", true)
		if err != nil {
			return err
		}
		os.Exit(res.ExitCode)
		return nil
	},
}

// printPendingHint avisa (em stderr, sem bloquear) se há dia anterior sem análise.
func printPendingHint() {
	if d, _ := journal.PendingReviewDate(); d != "" {
		fmt.Fprintf(os.Stderr, "Dica: há log de %s não analisado — rode 'bvr review'.\n", d)
	}
}

// validateRunShell aceita os shells de execução (pwsh/cmd/bash/zsh); claude não
// executa comandos avulsos.
func validateRunShell(s string) (config.Shell, error) {
	for _, v := range config.ValidRunShells {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("shell inválido %q (use: pwsh, cmd, bash, zsh)", s)
}

func init() {
	runCmd.Flags().SetInterspersed(false)
	runCmd.Flags().StringVar(&runShell, "shell", "", "shell de execução: pwsh|cmd|bash|zsh")
	rootCmd.AddCommand(runCmd)
}
