package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
	"harnessbeaver/internal/insights"
	"harnessbeaver/internal/journal"
	"harnessbeaver/internal/runner"
)

var shellShell string

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: i18n.T("REPL that executes and records commands (pwsh/cmd)"),
	Long:  i18n.T("Open an interactive session: each line typed is executed in the chosen\nshell and recorded in the journal. Type 'exit' or 'quit' (or Ctrl+D) to leave."),
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		shell := cfg.Settings.DefaultRunShell
		if shellShell != "" {
			s, err := validateRunShell(shellShell)
			if err != nil {
				return err
			}
			shell = s
		}

		reader := bufio.NewReader(os.Stdin)
		offerPendingReview(cfg, reader)

		fmt.Printf(i18n.T("bvr shell — %s. Commands are recorded in the journal. 'exit' to quit.\n"), shell)
		for {
			fmt.Printf("bvr[%s]> ", shell)
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Println()
				break
			}
			if err != nil {
				return err
			}
			command := strings.TrimSpace(line)
			if command == "" {
				continue
			}
			if command == "exit" || command == "quit" {
				break
			}
			if _, err := runner.Run(shell, command, "", true); err != nil {
				fmt.Fprintln(os.Stderr, i18n.T("error:"), err)
			}
		}
		return nil
	},
}

// offerPendingReview pergunta (uma vez por dia) se o usuário quer analisar o
// dia anterior pendente. Em caso afirmativo, gera e salva a análise.
func offerPendingReview(cfg *config.Config, reader *bufio.Reader) {
	if cfg.Settings.LastReviewOffer == journal.Today() {
		return
	}
	d, _ := journal.PendingReviewDate()
	if d == "" {
		return
	}
	fmt.Printf(i18n.T("There is an unanalysed log for %s. Analyse now? [y/N] "), d)
	ans, _ := reader.ReadString('\n')
	markOffered(cfg)
	if strings.EqualFold(strings.TrimSpace(ans), "y") {
		_, path, err := insights.Analyze(cfg, d, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, i18n.T("analysis failed:"), err)
			return
		}
		fmt.Printf(i18n.T("Analysis saved to: %s\n"), path)
	}
}

func init() {
	shellCmd.Flags().StringVar(&shellShell, "shell", "", i18n.T("execution shell: pwsh|cmd|bash|zsh"))
	rootCmd.AddCommand(shellCmd)
}
