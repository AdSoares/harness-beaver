// Package cmd implementa a CLI clássica (cobra) do HarnessBeaver.
// Sem subcomando, abre a TUI rica.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
	"harnessbeaver/internal/runner"
	"harnessbeaver/internal/tui"
)

// Version é injetada no build via -ldflags "-X harnessbeaver/cmd.Version=...".
var Version = "dev"

// langFlag is the global --lang override (en|pt). It wins over the config
// setting and BVR_LANG for the current invocation's runtime output.
var langFlag string

var rootCmd = &cobra.Command{
	Use:     "bvr",
	Short:   i18n.T("HarnessBeaver — multi-directory Claude Code launcher"),
	Version: Version,
	Long: i18n.T(`HarnessBeaver (bvr) opens Claude Code (or pwsh/cmd) across many
directories at once, in Windows Terminal tabs or separate windows.

With no arguments it opens the interactive interface (TUI). With subcommands it
works as a scriptable CLI.`),
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if langFlag != "" {
			i18n.Set(langFlag)
		}
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
		fmt.Fprintln(os.Stderr, i18n.T("error:"), err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", i18n.T("interface language: en|pt"))
}
