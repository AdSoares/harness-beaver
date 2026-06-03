package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/scanner"
)

var (
	scanRoot   string
	scanImport bool
	scanGroups bool
	scanDepth  int
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Escaneia um diretório em busca de projetos de código",
	Long: `Escaneia o diretório raiz (--root, ou settings.scanRoot, ou o diretório
atual) e lista subdiretórios que parecem projetos de código. Com --import, grava
no registro todos os candidatos ainda não cadastrados.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		root := scanRoot
		if root == "" {
			root = cfg.Settings.ScanRoot
		}
		if root == "" {
			if wd, err := os.Getwd(); err == nil {
				root = wd
			}
		}
		// Quando a raiz é informada explicitamente, memoriza para a próxima vez
		// (inclusive para a tela de scan da TUI).
		dirty := false
		if cmd.Flags().Changed("root") && !strings.EqualFold(cfg.Settings.ScanRoot, root) {
			cfg.Settings.ScanRoot = root
			dirty = true
		}
		candidates, err := scanner.Scan(root, scanDepth)
		if err != nil {
			return err
		}
		if scanGroups {
			candidates = append(candidates, scanner.Groups(candidates, root)...)
		}
		if len(candidates) == 0 {
			fmt.Printf("Nenhum projeto de código encontrado em %s\n", root)
			return nil
		}

		imported := 0
		for _, c := range candidates {
			exists := projectExistsAtPath(cfg, c.Path)
			tag := ""
			if exists {
				tag = " (já cadastrado)"
			}
			fmt.Printf("  %-28s %s  [%s]%s\n", c.Name, c.Path, strings.Join(c.Markers, ","), tag)
			if scanImport && !exists {
				if _, err := cfg.AddProject(c.Name, c.Path, ""); err == nil {
					imported++
				}
			}
		}
		if scanImport || dirty {
			if err := cfg.Save(); err != nil {
				return err
			}
		}
		if scanImport {
			fmt.Printf("\n%d projeto(s) importado(s).\n", imported)
		} else {
			fmt.Printf("\n%d candidato(s). Use --import para gravar.\n", len(candidates))
		}
		return nil
	},
}

func projectExistsAtPath(cfg *config.Config, path string) bool {
	for _, p := range cfg.Projects {
		if strings.EqualFold(p.Path, path) {
			return true
		}
	}
	return false
}

func init() {
	scanCmd.Flags().StringVar(&scanRoot, "root", "", "diretório raiz a escanear")
	scanCmd.Flags().BoolVar(&scanImport, "import", false, "importa os candidatos novos")
	scanCmd.Flags().BoolVar(&scanGroups, "groups", false, "também oferece diretórios-pai que agrupam ≥2 projetos")
	scanCmd.Flags().IntVar(&scanDepth, "depth", scanner.DefaultDepth, "profundidade máxima")
	rootCmd.AddCommand(scanCmd)
}
