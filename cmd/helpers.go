package cmd

import (
	"fmt"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/launcher"
)

// validateShell garante que a string corresponde a um shell suportado.
func validateShell(s string) (config.Shell, error) {
	if s == "" {
		return "", nil
	}
	for _, v := range config.ValidShells {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("shell inválido %q (use: claude, pwsh, cmd)", s)
}

// validateMode garante que a string corresponde a um modo suportado.
func validateMode(s string) (config.Mode, error) {
	if s == "" {
		return "", nil
	}
	for _, v := range config.ValidModes {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("modo inválido %q (use: tabs, windows)", s)
}

// resolveOpen interpreta os argumentos (ids de pacote e/ou projeto) e as flags
// de modo/shell, devolvendo o modo efetivo e os itens prontos para o launcher.
func resolveOpen(cfg *config.Config, args []string, modeFlag, shellFlag string) (config.Mode, []launcher.Item, error) {
	mode, err := validateMode(modeFlag)
	if err != nil {
		return "", nil, err
	}
	shell, err := validateShell(shellFlag)
	if err != nil {
		return "", nil, err
	}

	var projects []config.Project
	var pkgMode config.Mode
	var pkgShell config.Shell

	for _, a := range args {
		if pk, ok := cfg.FindPackage(a); ok {
			pkgMode = pk.Mode
			pkgShell = pk.Shell
			projects = append(projects, cfg.ProjectsByIDs(pk.ProjectIDs)...)
			continue
		}
		if pr, ok := cfg.FindProject(a); ok {
			projects = append(projects, *pr)
			continue
		}
		return "", nil, fmt.Errorf("id não encontrado (projeto ou pacote): %s", a)
	}
	if len(projects) == 0 {
		return "", nil, fmt.Errorf("nenhum projeto para abrir")
	}

	// Modo: default global < pacote < flag.
	effMode := config.ModeTabs
	if cfg.Settings.DefaultMode != "" {
		effMode = cfg.Settings.DefaultMode
	}
	if pkgMode != "" {
		effMode = pkgMode
	}
	if mode != "" {
		effMode = mode
	}

	// Shell explícito: pacote < flag.
	explicitShell := pkgShell
	if shell != "" {
		explicitShell = shell
	}

	items := make([]launcher.Item, 0, len(projects))
	for _, p := range projects {
		items = append(items, launcher.Item{
			Project: p,
			Shell:   launcher.ResolveShell(explicitShell, p, cfg.Settings.DefaultShell),
		})
	}
	return effMode, items, nil
}
