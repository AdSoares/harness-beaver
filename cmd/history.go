package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

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
	Short: "Lista o histórico de comandos do diário (com filtros)",
	Long: `Mostra os comandos registrados nos últimos --days dias, do mais recente
para o mais antigo. Filtre por projeto, texto ou apenas os que falharam.`,
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
			fmt.Println("Nenhuma entrada no histórico para os filtros dados.")
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
	historyCmd.Flags().StringVar(&histProject, "project", "", "filtra por id de projeto")
	historyCmd.Flags().StringVar(&histGrep, "grep", "", "filtra por texto no comando")
	historyCmd.Flags().BoolVar(&histFailed, "failed", false, "apenas comandos com exit != 0")
	historyCmd.Flags().IntVar(&histDays, "days", 7, "quantos dias para trás")
	historyCmd.Flags().IntVar(&histLimit, "limit", 50, "máximo de entradas exibidas")
	rootCmd.AddCommand(historyCmd)
}
