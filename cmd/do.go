package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
	"harnessbeaver/internal/runner"
)

var doShell string

var doCmd = &cobra.Command{
	Use:   i18n.T("do <name> [args...]"),
	Short: i18n.T("Run a registered alias, expanding {1}..{N} and {*}"),
	Long: i18n.T(`Run the alias command <name>. Extra arguments replace the
placeholders {1}, {2}, … and {*} (all together); without placeholders, they are appended.

Examples:
  bvr alias set deploy "npm run build && npm run deploy"
  bvr do deploy
  bvr alias set commit "git commit -m {1}"
  bvr do commit "fix: adjustment"`),
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		tmpl, ok := cfg.Alias(args[0])
		if !ok {
			return fmt.Errorf(i18n.T("alias not found: %s (see 'bvr alias')"), args[0])
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
	doCmd.Flags().StringVar(&doShell, "shell", "", i18n.T("execution shell: pwsh|cmd|bash|zsh"))
	rootCmd.AddCommand(doCmd)
}
