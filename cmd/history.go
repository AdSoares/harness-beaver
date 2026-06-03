package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/i18n"
	"harnessbeaver/internal/journal"
)

var (
	histProject string
	histGrep    string
	histFailed  bool
	histDays    int
	histLimit   int
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: i18n.T("List command history from the journal (with filters)"),
	Long: i18n.T(`Shows commands recorded in the last --days days, from most recent
to oldest. Filter by project, text, or only those that failed.`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		to := journal.Today()
		from := time.Now().AddDate(0, 0, -(histDays - 1)).Format(journal.DateLayout)

		entries, err := journal.ReadRange(from, to)
		if err != nil {
			return err
		}
		entries = journal.Filter(entries, journal.FilterOpts{
			ProjectID:  histProject,
			Grep:       histGrep,
			FailedOnly: histFailed,
		})
		if len(entries) == 0 {
			fmt.Println(i18n.T("No history entries for the given filters."))
			return nil
		}

		// Mais recentes primeiro, limitado.
		shown := 0
		for i := len(entries) - 1; i >= 0 && shown < histLimit; i-- {
			e := entries[i]
			proj := e.ProjectID
			if proj == "" {
				proj = "-"
			}
			fmt.Printf("  %s [%s] (%s) exit=%d  %s\n",
				e.Time.Format("2006-01-02 15:04:05"), e.Shell, proj, e.ExitCode, e.Command)
			shown++
		}
		return nil
	},
}

func init() {
	historyCmd.Flags().StringVar(&histProject, "project", "", i18n.T("filter by project id"))
	historyCmd.Flags().StringVar(&histGrep, "grep", "", i18n.T("filter by text in command"))
	historyCmd.Flags().BoolVar(&histFailed, "failed", false, i18n.T("only commands with exit != 0"))
	historyCmd.Flags().IntVar(&histDays, "days", 7, i18n.T("how many days back"))
	historyCmd.Flags().IntVar(&histLimit, "limit", 50, i18n.T("maximum number of entries shown"))
	rootCmd.AddCommand(historyCmd)
}
