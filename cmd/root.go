// Package cmd implementa a CLI clássica (cobra) do HarnessBeaver.
// Sem subcomando, abre a TUI rica.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/runner"
	"harnessbeaver/internal/tui"
)

// Version é injetada no build via -ldflags "-X harnessbeaver/cmd.Version=...".
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "bvr",
	Short: "HarnessBeaver — launcher de Claude Code multi-diretório",
	Version: Version,
	Long: `HarnessBeaver (bvr) abre o Claude Code (ou pwsh/cmd) em vários
diretórios de uma vez, em abas do Windows Terminal ou janelas separadas.

Sem argumentos, abre a interface interativa (TUI). Com subcomandos, funciona
como CLI scriptável.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Habilita a tag de projeto no diário (cwd -> projeto) para todos os comandos.
		if cfg, err := config.Load(); err == nil {
			runner.ResolveProjectID = cfg.ProjectIDForPath
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

func runTUI() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return tui.Run(cfg)
}

// Execute é o ponto de entrada chamado por main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}
