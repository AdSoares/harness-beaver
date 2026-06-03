package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
	"harnessbeaver/internal/insights"
	"harnessbeaver/internal/journal"
)

var (
	reviewDryRun bool
	reviewForce  bool
)

var reviewCmd = &cobra.Command{
	Use:   i18n.T("review [date]"),
	Short: i18n.T("Analyse a day's learnings (AI)"),
	Long:  i18n.T("Analyse the commands recorded on a given day and generate a markdown\nlearnings file at ~/.harnessbeaver/learnings/<date>.md.\n\nWith no argument, uses the pending day (previous, not yet analysed) or today.\nDate in YYYY-MM-DD format. Use --dry-run to see the prompt without calling the AI."),
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		date := ""
		if len(args) == 1 {
			if _, err := time.Parse(journal.DateLayout, args[0]); err != nil {
				return fmt.Errorf(i18n.T("invalid date %q (use YYYY-MM-DD)"), args[0])
			}
			date = args[0]
		} else if d, _ := journal.PendingReviewDate(); d != "" {
			date = d
		} else {
			date = journal.Today()
		}

		if reviewDryRun {
			prompt, _, err := insights.Analyze(cfg, date, true)
			if err != nil {
				return err
			}
			fmt.Println(prompt)
			return nil
		}

		// Cache: se já existe análise e não foi pedido --force, mostra a salva.
		if journal.HasLearning(date) && !reviewForce {
			md, err := journal.ReadLearning(date)
			if err != nil {
				return err
			}
			fmt.Printf(i18n.T("Existing analysis for %s (use --force to redo):\n\n%s\n"), date, md)
			return nil
		}

		md, path, err := insights.Analyze(cfg, date, false)
		if err != nil {
			return err
		}
		markOffered(cfg)
		fmt.Printf(i18n.T("Analysis saved to: %s\n\n%s\n"), path, md)
		return nil
	},
}

// markOffered registra que a oferta/análise já ocorreu hoje (evita auto-nag).
func markOffered(cfg *config.Config) {
	cfg.Settings.LastReviewOffer = journal.Today()
	_ = cfg.Save()
}

func init() {
	reviewCmd.Flags().BoolVar(&reviewDryRun, "dry-run", false, i18n.T("print the prompt without calling the AI"))
	reviewCmd.Flags().BoolVar(&reviewForce, "force", false, i18n.T("redo analysis even if a cached result exists"))
	rootCmd.AddCommand(reviewCmd)
}
