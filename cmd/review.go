package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/insights"
	"harnessbeaver/internal/journal"
)

var (
	reviewDryRun bool
	reviewForce  bool
)

var reviewCmd = &cobra.Command{
	Use:   "review [data]",
	Short: "Analisa os aprendizados de um dia (IA)",
	Long: `Analisa os comandos registrados em um dia e gera um markdown de
aprendizados em ~/.harnessbeaver/learnings/<data>.md.

Sem argumento, usa o dia pendente (anterior, sem análise) ou hoje.
Data no formato AAAA-MM-DD. Use --dry-run para ver o prompt sem chamar a IA.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		date := ""
		if len(args) == 1 {
			if _, err := time.Parse(journal.DateLayout, args[0]); err != nil {
				return fmt.Errorf("data inválida %q (use AAAA-MM-DD)", args[0])
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
			fmt.Printf("Análise existente de %s (use --force para refazer):\n\n%s\n", date, md)
			return nil
		}

		md, path, err := insights.Analyze(cfg, date, false)
		if err != nil {
			return err
		}
		markOffered(cfg)
		fmt.Printf("Análise salva em: %s\n\n%s\n", path, md)
		return nil
	},
}

// markOffered registra que a oferta/análise já ocorreu hoje (evita auto-nag).
func markOffered(cfg *config.Config) {
	cfg.Settings.LastReviewOffer = journal.Today()
	_ = cfg.Save()
}

func init() {
	reviewCmd.Flags().BoolVar(&reviewDryRun, "dry-run", false, "imprime o prompt sem chamar a IA")
	reviewCmd.Flags().BoolVar(&reviewForce, "force", false, "refaz a análise mesmo se já houver cache")
	rootCmd.AddCommand(reviewCmd)
}
