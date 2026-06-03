package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/launcher"
)

var (
	openMode     string
	openShell    string
	openDryRun   bool
	openLayout   string
	openNoLayout bool
)

var openCmd = &cobra.Command{
	Use:   "open <pkgId|projId> [outros...]",
	Short: "Abre projetos ou um pacote no Windows Terminal",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// Abertura com layout: aplica quando há um único projeto e um layout
		// (flag --layout, ou o layout padrão do projeto).
		if layout, project, ok, err := resolveLayoutOpen(cfg, args); err != nil {
			return err
		} else if ok {
			line, err := launcher.LaunchLayout(*project, *layout, cfg.Settings.TabColors, openDryRun)
			if err != nil {
				return err
			}
			if openDryRun {
				fmt.Println(line)
			} else {
				fmt.Printf("Aberto %s com o layout %s.\n", project.Name, layout.ID)
			}
			return nil
		}

		mode, items, err := resolveOpen(cfg, args, openMode, openShell)
		if err != nil {
			return err
		}
		lines, err := launcher.Launch(mode, items, cfg.Settings.TabColors, openDryRun)
		if err != nil {
			return err
		}
		if openDryRun {
			for _, l := range lines {
				fmt.Println(l)
			}
			return nil
		}
		fmt.Printf("Aberto %d projeto(s) em modo %s.\n", len(items), mode)
		return nil
	},
}

// resolveLayoutOpen decide se a abertura deve usar um layout. Retorna ok=true
// quando há exatamente um projeto alvo e um layout aplicável (flag --layout ou
// o LayoutID do projeto, salvo --no-layout).
func resolveLayoutOpen(cfg *config.Config, args []string) (*config.Layout, *config.Project, bool, error) {
	if openNoLayout {
		return nil, nil, false, nil
	}
	// Resolve um único projeto a partir dos args (não expande pacotes aqui).
	if len(args) != 1 {
		if openLayout != "" {
			return nil, nil, false, fmt.Errorf("--layout requer exatamente um projeto")
		}
		return nil, nil, false, nil
	}
	project, ok := cfg.FindProject(args[0])
	if !ok {
		if openLayout != "" {
			return nil, nil, false, fmt.Errorf("--layout requer um projeto (não um pacote): %s", args[0])
		}
		return nil, nil, false, nil
	}

	layoutID := openLayout
	if layoutID == "" {
		layoutID = project.LayoutID
	}
	if layoutID == "" {
		return nil, nil, false, nil // sem layout: segue abertura normal
	}
	layout, ok := cfg.FindLayout(layoutID)
	if !ok {
		return nil, nil, false, fmt.Errorf("layout não encontrado: %s", layoutID)
	}
	return layout, project, true, nil
}

func init() {
	openCmd.Flags().StringVar(&openMode, "mode", "", "modo de abertura: tabs|windows")
	openCmd.Flags().StringVar(&openShell, "shell", "", "shell por aba: claude|pwsh|cmd")
	openCmd.Flags().BoolVar(&openDryRun, "dry-run", false, "imprime os comandos wt sem executar")
	openCmd.Flags().StringVar(&openLayout, "layout", "", "abre o projeto com este layout")
	openCmd.Flags().BoolVar(&openNoLayout, "no-layout", false, "ignora o layout padrão do projeto")
	rootCmd.AddCommand(openCmd)
}
