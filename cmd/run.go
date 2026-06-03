package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
	"harnessbeaver/internal/journal"
	"harnessbeaver/internal/runner"
)

var runShell string

var runCmd = &cobra.Command{
	Use:                i18n.T("run <command...>"),
	Short:              i18n.T("Execute a command (pwsh/cmd) and record it in the journal"),
	Long:               i18n.T("Execute a command on behalf of the user in PowerShell (default) or cmd and\nrecord the command + output in the journal (~/.harnessbeaver/logs).\n\nFlags must come before the command. Examples:\n  bvr run echo hello\n  bvr run --shell cmd \"echo hi & dir\""),
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
		fmt.Fprintf(os.Stderr, i18n.T("Hint: there is an unanalysed log for %s — run 'bvr review'.\n"), d)
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
	return "", fmt.Errorf(i18n.T("invalid shell %q (use: pwsh, cmd, bash, zsh)"), s)
}

func init() {
	runCmd.Flags().SetInterspersed(false)
	runCmd.Flags().StringVar(&runShell, "shell", "", i18n.T("execution shell: pwsh|cmd|bash|zsh"))
	rootCmd.AddCommand(runCmd)
}
