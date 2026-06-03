package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/launcher"
)

var layoutCmd = &cobra.Command{
	Use:   "layout",
	Short: "Gerencia layouts (arranjos de abas/painéis) para abrir projetos",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Layouts) == 0 {
			fmt.Println("Nenhum layout. Crie com: bvr layout add <nome> [--preset dev|triple]")
			return nil
		}
		for _, l := range cfg.Layouts {
			panes := 0
			for _, t := range l.Tabs {
				panes += len(t.Panes)
			}
			fmt.Printf("  %-16s %-24s %d aba(s), %d painel(éis)\n", l.ID, l.Name, len(l.Tabs), panes)
		}
		return nil
	},
}

var (
	layoutPreset  string
	layoutTabTtl  string
	paneTab       int
	paneShell     string
	paneCommandF  string
	paneSplit     string
	paneDir       string
	showProjectID string
)

var layoutAddCmd = &cobra.Command{
	Use:   "add <nome>",
	Short: "Cria um layout (preset: claude|dev|triple)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		l, err := cfg.AddLayout(args[0], layoutPreset)
		if err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Layout criado: %s\n", l.ID)
		return nil
	},
}

var layoutTabCmd = &cobra.Command{
	Use:   "tab <layoutId>",
	Short: "Adiciona uma aba ao layout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		l, ok := cfg.FindLayout(args[0])
		if !ok {
			return fmt.Errorf("layout não encontrado: %s", args[0])
		}
		l.Tabs = append(l.Tabs, config.Tab{Title: layoutTabTtl})
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Aba %d adicionada ao layout %s.\n", len(l.Tabs), l.ID)
		return nil
	},
}

var layoutPaneCmd = &cobra.Command{
	Use:   "pane <layoutId>",
	Short: "Adiciona um painel a uma aba do layout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell, err := validateLayoutShell(paneShell)
		if err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		l, ok := cfg.FindLayout(args[0])
		if !ok {
			return fmt.Errorf("layout não encontrado: %s", args[0])
		}
		if paneTab < 1 || paneTab > len(l.Tabs) {
			return fmt.Errorf("aba %d inválida (o layout tem %d aba(s))", paneTab, len(l.Tabs))
		}
		split := paneSplit
		if split != "" && split != "H" && split != "V" {
			return fmt.Errorf("split inválido %q (use H ou V)", split)
		}
		l.Tabs[paneTab-1].Panes = append(l.Tabs[paneTab-1].Panes, config.Pane{
			Shell: shell, Command: paneCommandF, Dir: paneDir, Split: split,
		})
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Painel adicionado à aba %d do layout %s.\n", paneTab, l.ID)
		return nil
	},
}

var layoutRmCmd = &cobra.Command{
	Use:     "rm <layoutId>",
	Aliases: []string{"remove"},
	Short:   "Remove um layout",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemoveLayout(args[0]) {
			return fmt.Errorf("layout não encontrado: %s", args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Layout removido: %s\n", args[0])
		return nil
	},
}

var layoutAssignCmd = &cobra.Command{
	Use:   "assign <layoutId> <projId>",
	Short: "Define o layout padrão de um projeto",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if _, ok := cfg.FindLayout(args[0]); !ok {
			return fmt.Errorf("layout não encontrado: %s", args[0])
		}
		p, ok := cfg.FindProject(args[1])
		if !ok {
			return fmt.Errorf("projeto não encontrado: %s", args[1])
		}
		p.LayoutID = args[0]
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("Projeto %s agora abre com o layout %s.\n", p.ID, args[0])
		return nil
	},
}

var layoutShowCmd = &cobra.Command{
	Use:   "show <layoutId>",
	Short: "Mostra o comando wt que o layout geraria",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		l, ok := cfg.FindLayout(args[0])
		if !ok {
			return fmt.Errorf("layout não encontrado: %s", args[0])
		}
		project := config.Project{Name: "<projeto>", Path: "<dir-do-projeto>"}
		if showProjectID != "" {
			p, ok := cfg.FindProject(showProjectID)
			if !ok {
				return fmt.Errorf("projeto não encontrado: %s", showProjectID)
			}
			project = *p
		}
		fmt.Println(launcher.PlanLayout(project, *l, cfg.Settings.TabColors))
		return nil
	},
}

// validateLayoutShell aceita os shells válidos num painel (inclui claude).
func validateLayoutShell(s string) (config.Shell, error) {
	if s == "" {
		return config.ShellClaude, nil
	}
	switch config.Shell(s) {
	case config.ShellClaude, config.ShellPwsh, config.ShellCmd, config.ShellBash, config.ShellZsh:
		return config.Shell(s), nil
	}
	return "", fmt.Errorf("shell inválido %q (use: claude, pwsh, cmd, bash, zsh)", s)
}

func init() {
	layoutAddCmd.Flags().StringVar(&layoutPreset, "preset", "", "preset inicial: claude|dev|triple")
	layoutTabCmd.Flags().StringVar(&layoutTabTtl, "title", "", "título da aba")
	layoutPaneCmd.Flags().IntVar(&paneTab, "tab", 1, "número da aba (1-based)")
	layoutPaneCmd.Flags().StringVar(&paneShell, "shell", "", "shell do painel: claude|pwsh|cmd|bash|zsh")
	layoutPaneCmd.Flags().StringVar(&paneCommandF, "command", "", "comando a rodar no painel")
	layoutPaneCmd.Flags().StringVar(&paneSplit, "split", "V", "divisão a partir do painel anterior: H|V")
	layoutPaneCmd.Flags().StringVar(&paneDir, "dir", "", "cwd do painel (vazio = dir do projeto)")
	layoutShowCmd.Flags().StringVar(&showProjectID, "project", "", "projeto p/ preencher os caminhos")

	layoutCmd.AddCommand(layoutAddCmd, layoutTabCmd, layoutPaneCmd, layoutRmCmd, layoutAssignCmd, layoutShowCmd)
	rootCmd.AddCommand(layoutCmd)
}
