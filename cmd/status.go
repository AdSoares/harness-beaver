package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/gitstatus"
	"harnessbeaver/internal/i18n"
)

var statusCmd = &cobra.Command{
	Use:   "status [pkgId|projId ...]",
	Short: i18n.T("Show git status of projects (branch, changes, ahead/behind)"),
	Long:  i18n.T("With no arguments, shows the git status of all registered projects.\nWith project and/or package ids, restricts to those. Queries run in parallel."),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		projects, err := collectProjects(cfg, args)
		if err != nil {
			return err
		}
		if len(projects) == 0 {
			fmt.Println(i18n.T("No registered projects."))
			return nil
		}

		paths := make([]string, len(projects))
		for i, p := range projects {
			paths[i] = p.Path
		}
		res := gitstatus.QueryMany(paths)

		for _, p := range projects {
			label := res[p.Path].Label()
			if label == "" {
				label = i18n.T("(not a git repository)")
			}
			fmt.Printf("  %-20s %-30s %s\n", p.ID, label, p.Path)
		}
		return nil
	},
}

// collectProjects resolve args (ids de projeto/pacote) para uma lista de
// projetos sem duplicatas; sem args, devolve todos.
func collectProjects(cfg *config.Config, args []string) ([]config.Project, error) {
	if len(args) == 0 {
		return cfg.Projects, nil
	}
	var out []config.Project
	seen := map[string]bool{}
	add := func(p config.Project) {
		if !seen[p.ID] {
			seen[p.ID] = true
			out = append(out, p)
		}
	}
	for _, a := range args {
		if pk, ok := cfg.FindPackage(a); ok {
			for _, p := range cfg.ProjectsByIDs(pk.ProjectIDs) {
				add(p)
			}
			continue
		}
		if p, ok := cfg.FindProject(a); ok {
			add(*p)
			continue
		}
		return nil, fmt.Errorf(i18n.T("id not found (project or package): %s"), a)
	}
	return out, nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
