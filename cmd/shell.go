package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/insights"
	"harnessbeaver/internal/journal"
	"harnessbeaver/internal/runner"
)

var shellShell string

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "REPL que executa e registra comandos (pwsh/cmd)",
	Long: `Abre uma sessão interativa: cada linha digitada é executada no shell
escolhido e registrada no diário. Digite 'exit' ou 'quit' (ou Ctrl+D) para sair.`,
	Args: cobra.NoArgs,
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

		fmt.Printf("bvr shell — %s. Comandos são registrados no diário. 'exit' para sair.\n", shell)
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
				fmt.Fprintln(os.Stderr, "erro:", err)
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
	fmt.Printf("Há log de %s ainda não analisado. Analisar agora? [s/N] ", d)
	ans, _ := reader.ReadString('\n')
	markOffered(cfg)
	if strings.EqualFold(strings.TrimSpace(ans), "s") {
		_, path, err := insights.Analyze(cfg, d, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, "falha na análise:", err)
			return
		}
		fmt.Printf("Análise salva em: %s\n", path)
	}
}

func init() {
	shellCmd.Flags().StringVar(&shellShell, "shell", "", "shell de execução: pwsh|cmd|bash|zsh")
	rootCmd.AddCommand(shellCmd)
}
