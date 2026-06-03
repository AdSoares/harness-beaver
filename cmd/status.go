package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/gitstatus"
)

var statusCmd = &cobra.Command{
	Use:   "status [pkgId|projId ...]",
	Short: "Mostra o status git dos projetos (branch, alterações, ahead/behind)",
	Long: `Sem argumentos, mostra o status git de todos os projetos registrados.
Com ids de projeto e/ou pacote, restringe a esses. As consultas rodam em paralelo.`,
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
			fmt.Println("Nenhum projeto registrado.")
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
				label = "(não é repositório git)"
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
		return nil, fmt.Errorf("id não encontrado (projeto ou pacote): %s", a)
	}
	return out, nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
