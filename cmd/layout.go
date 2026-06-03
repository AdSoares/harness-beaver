package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
	"harnessbeaver/internal/launcher"
)

var layoutCmd = &cobra.Command{
	Use:   "layout",
	Short: i18n.T("Manage layouts (tab/pane arrangements) for opening projects"),
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Layouts) == 0 {
			fmt.Println(i18n.T("No layouts. Create one with: bvr layout add <name> [--preset dev|triple]"))
			return nil
		}
		for _, l := range cfg.Layouts {
			panes := 0
			for _, t := range l.Tabs {
				panes += len(t.Panes)
			}
			fmt.Printf(i18n.T("  %-16s %-24s %d tab(s), %d pane(s)\n"), l.ID, l.Name, len(l.Tabs), panes)
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
	Use:   i18n.T("add <name>"),
	Short: i18n.T("Create a layout (preset: claude|dev|triple)"),
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
		fmt.Printf(i18n.T("Layout created: %s\n"), l.ID)
		return nil
	},
}

var layoutTabCmd = &cobra.Command{
	Use:   i18n.T("tab <layoutId>"),
	Short: i18n.T("Add a tab to the layout"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		l, ok := cfg.FindLayout(args[0])
		if !ok {
			return fmt.Errorf(i18n.T("layout not found: %s"), args[0])
		}
		l.Tabs = append(l.Tabs, config.Tab{Title: layoutTabTtl})
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Tab %d added to layout %s.\n"), len(l.Tabs), l.ID)
		return nil
	},
}

var layoutPaneCmd = &cobra.Command{
	Use:   i18n.T("pane <layoutId>"),
	Short: i18n.T("Add a pane to a tab in the layout"),
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
			return fmt.Errorf(i18n.T("layout not found: %s"), args[0])
		}
		if paneTab < 1 || paneTab > len(l.Tabs) {
			return fmt.Errorf(i18n.T("invalid tab %d (the layout has %d tab(s))"), paneTab, len(l.Tabs))
		}
		split := paneSplit
		if split != "" && split != "H" && split != "V" {
			return fmt.Errorf(i18n.T("invalid split %q (use H or V)"), split)
		}
		l.Tabs[paneTab-1].Panes = append(l.Tabs[paneTab-1].Panes, config.Pane{
			Shell: shell, Command: paneCommandF, Dir: paneDir, Split: split,
		})
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Pane added to tab %d of layout %s.\n"), paneTab, l.ID)
		return nil
	},
}

var layoutRmCmd = &cobra.Command{
	Use:     i18n.T("rm <layoutId>"),
	Aliases: []string{"remove"},
	Short:   i18n.T("Remove a layout"),
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemoveLayout(args[0]) {
			return fmt.Errorf(i18n.T("layout not found: %s"), args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Layout removed: %s\n"), args[0])
		return nil
	},
}

var layoutAssignCmd = &cobra.Command{
	Use:   i18n.T("assign <layoutId> <projId>"),
	Short: i18n.T("Set the default layout for a project"),
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if _, ok := cfg.FindLayout(args[0]); !ok {
			return fmt.Errorf(i18n.T("layout not found: %s"), args[0])
		}
		p, ok := cfg.FindProject(args[1])
		if !ok {
			return fmt.Errorf(i18n.T("project not found: %s"), args[1])
		}
		p.LayoutID = args[0]
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Project %s will now open with layout %s.\n"), p.ID, args[0])
		return nil
	},
}

var layoutShowCmd = &cobra.Command{
	Use:   i18n.T("show <layoutId>"),
	Short: i18n.T("Show the wt command that the layout would generate"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		l, ok := cfg.FindLayout(args[0])
		if !ok {
			return fmt.Errorf(i18n.T("layout not found: %s"), args[0])
		}
		project := config.Project{Name: "<project>", Path: "<project-dir>"}
		if showProjectID != "" {
			p, ok := cfg.FindProject(showProjectID)
			if !ok {
				return fmt.Errorf(i18n.T("project not found: %s"), showProjectID)
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
	return "", fmt.Errorf(i18n.T("invalid shell %q (use: claude, pwsh, cmd, bash, zsh)"), s)
}

func init() {
	layoutAddCmd.Flags().StringVar(&layoutPreset, "preset", "", i18n.T("initial preset: claude|dev|triple"))
	layoutTabCmd.Flags().StringVar(&layoutTabTtl, "title", "", i18n.T("tab title"))
	layoutPaneCmd.Flags().IntVar(&paneTab, "tab", 1, i18n.T("tab number (1-based)"))
	layoutPaneCmd.Flags().StringVar(&paneShell, "shell", "", i18n.T("pane shell: claude|pwsh|cmd|bash|zsh"))
	layoutPaneCmd.Flags().StringVar(&paneCommandF, "command", "", i18n.T("command to run in the pane"))
	layoutPaneCmd.Flags().StringVar(&paneSplit, "split", "V", i18n.T("split from previous pane: H|V"))
	layoutPaneCmd.Flags().StringVar(&paneDir, "dir", "", i18n.T("pane cwd (empty = project dir)"))
	layoutShowCmd.Flags().StringVar(&showProjectID, "project", "", i18n.T("project to fill in paths"))

	layoutCmd.AddCommand(layoutAddCmd, layoutTabCmd, layoutPaneCmd, layoutRmCmd, layoutAssignCmd, layoutShowCmd)
	rootCmd.AddCommand(layoutCmd)
}
