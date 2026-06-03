package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
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
	Short: i18n.T("Scan a directory for code projects"),
	Long:  i18n.T("Scan the root directory (--root, or settings.scanRoot, or the current\ndirectory) and list subdirectories that look like code projects. With --import,\nrecords all candidates not yet registered."),
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
			fmt.Printf(i18n.T("No code projects found in %s\n"), root)
			return nil
		}

		imported := 0
		for _, c := range candidates {
			exists := projectExistsAtPath(cfg, c.Path)
			tag := ""
			if exists {
				tag = i18n.T(" (already registered)")
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
			fmt.Printf(i18n.T("\n%d project(s) imported.\n"), imported)
		} else {
			fmt.Printf(i18n.T("\n%d candidate(s). Use --import to register.\n"), len(candidates))
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
	scanCmd.Flags().StringVar(&scanRoot, "root", "", i18n.T("root directory to scan"))
	scanCmd.Flags().BoolVar(&scanImport, "import", false, i18n.T("import new candidates"))
	scanCmd.Flags().BoolVar(&scanGroups, "groups", false, i18n.T("also offer parent directories that group ≥2 projects"))
	scanCmd.Flags().IntVar(&scanDepth, "depth", scanner.DefaultDepth, i18n.T("maximum depth"))
	rootCmd.AddCommand(scanCmd)
}
