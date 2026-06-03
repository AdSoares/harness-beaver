// Package launcher monta e dispara comandos do Windows Terminal (wt) para
// abrir os projetos selecionados em abas ou janelas separadas.
package launcher

import (
	"os/exec"
	"strings"

	"harnessbeaver/internal/config"
)

// Item é um projeto resolvido com o shell que deve rodar nele.
type Item struct {
	Project config.Project
	Shell   config.Shell
}

// ResolveShell decide o shell efetivo: flag/pacote explícito > shell do
// projeto > default global > claude.
func ResolveShell(explicit config.Shell, p config.Project, def config.Shell) config.Shell {
	if explicit != "" {
		return explicit
	}
	if p.DefaultShell != "" {
		return p.DefaultShell
	}
	if def != "" {
		return def
	}
	return config.ShellClaude
}

// tabPalette são cores cíclicas para colorir abas/janelas quando habilitado.
var tabPalette = []string{
	"#5BC0EB", "#FDE74C", "#9BC53D", "#E55934", "#FA7921",
	"#C3A6FF", "#46B1C9", "#F26196", "#7AE582", "#FF8C42",
}

// colorFor devolve a cor (hex) do i-ésimo item, ciclando a paleta.
func colorFor(i int) string { return tabPalette[i%len(tabPalette)] }

// shellCommand mapeia o shell para os argumentos que iniciam a aba/janela.
// Para claude usamos um pwsh hospedeiro com -NoExit para a aba não fechar
// caso o claude encerre.
func shellCommand(s config.Shell) []string {
	switch s {
	case config.ShellPwsh:
		return []string{"pwsh", "-NoExit"}
	case config.ShellCmd:
		return []string{"cmd", "/k"}
	case config.ShellClaude:
		fallthrough
	default:
		return []string{"pwsh", "-NoExit", "-Command", "claude"}
	}
}

// tabArgs monta os argumentos de um único `new-tab`. `--suppressApplicationTitle`
// impede que o shell/app sobrescreva o título; color (hex) é opcional.
func tabArgs(it Item, color string) []string {
	args := []string{"new-tab", "--title", it.Project.Name, "--suppressApplicationTitle"}
	if color != "" {
		args = append(args, "--tabColor", color)
	}
	args = append(args, "-d", it.Project.Path)
	return append(args, shellCommand(it.Shell)...)
}

// buildTabsArgs monta os argumentos de uma única invocação `wt` com N abas,
// separadas pelo argumento literal ";" (o wt o trata como separador de
// subcomando; o os/exec do Windows não o coloca entre aspas).
func buildTabsArgs(items []Item, colorize bool) []string {
	var args []string
	for i, it := range items {
		if i > 0 {
			args = append(args, ";")
		}
		color := ""
		if colorize {
			color = colorFor(i)
		}
		args = append(args, tabArgs(it, color)...)
	}
	return args
}

// buildWindowArgs monta os argumentos de uma janela nova para um projeto.
func buildWindowArgs(it Item, color string) []string {
	return append([]string{"-w", "new"}, tabArgs(it, color)...)
}

// Plan devolve, sem executar, as linhas de comando `wt` que seriam disparadas.
func Plan(mode config.Mode, items []Item, colorize bool) []string {
	var lines []string
	if mode == config.ModeWindows {
		for i, it := range items {
			color := ""
			if colorize {
				color = colorFor(i)
			}
			lines = append(lines, "wt "+strings.Join(buildWindowArgs(it, color), " "))
		}
		return lines
	}
	return []string{"wt " + strings.Join(buildTabsArgs(items, colorize), " ")}
}

// Launch dispara o(s) comando(s) `wt`. Se dryRun for true, apenas devolve as
// linhas planejadas sem executar nada.
func Launch(mode config.Mode, items []Item, colorize, dryRun bool) ([]string, error) {
	lines := Plan(mode, items, colorize)
	if dryRun {
		return lines, nil
	}
	if mode == config.ModeWindows {
		for i, it := range items {
			color := ""
			if colorize {
				color = colorFor(i)
			}
			if err := exec.Command("wt", buildWindowArgs(it, color)...).Start(); err != nil {
				return lines, err
			}
		}
		return lines, nil
	}
	if err := exec.Command("wt", buildTabsArgs(items, colorize)...).Start(); err != nil {
		return lines, err
	}
	return lines, nil
}

// --- Layouts (Fase 3) ---

// paneCommand monta a invocação do shell de um painel, embutindo o comando
// opcional e mantendo o painel aberto após executá-lo.
func paneCommand(shell config.Shell, command string) []string {
	switch shell {
	case config.ShellPwsh:
		if command != "" {
			return []string{"pwsh", "-NoExit", "-Command", command}
		}
		return []string{"pwsh", "-NoExit"}
	case config.ShellCmd:
		if command != "" {
			return []string{"cmd", "/k", command}
		}
		return []string{"cmd", "/k"}
	case config.ShellBash:
		if command != "" {
			return []string{"bash", "-lc", command + "; exec bash"}
		}
		return []string{"bash"}
	case config.ShellZsh:
		if command != "" {
			return []string{"zsh", "-lc", command + "; exec zsh"}
		}
		return []string{"zsh"}
	case config.ShellClaude:
		fallthrough
	default:
		if command != "" {
			return []string{"pwsh", "-NoExit", "-Command", command + "; claude"}
		}
		return []string{"pwsh", "-NoExit", "-Command", "claude"}
	}
}

// buildLayoutArgs monta os argumentos `wt` de um layout para um projeto: o
// primeiro painel de cada aba vira `new-tab`; os demais, `split-pane -H|-V`.
func buildLayoutArgs(project config.Project, layout config.Layout, colorize bool) []string {
	var args []string
	first := true
	for ti, tab := range layout.Tabs {
		title := tab.Title
		if title == "" {
			title = project.Name
		}
		for pi, pane := range tab.Panes {
			if !first {
				args = append(args, ";")
			}
			first = false

			dir := pane.Dir
			if dir == "" {
				dir = project.Path
			}
			if pi == 0 {
				args = append(args, "new-tab")
			} else {
				args = append(args, "split-pane")
				if strings.EqualFold(pane.Split, "H") {
					args = append(args, "-H")
				} else {
					args = append(args, "-V")
				}
			}
			// Título fixo (nome do projeto/aba) em todos os painéis, sem deixar o app sobrescrever.
			args = append(args, "--title", title, "--suppressApplicationTitle")
			if colorize && pi == 0 {
				args = append(args, "--tabColor", colorFor(ti))
			}
			args = append(args, "-d", dir)
			args = append(args, paneCommand(pane.Shell, pane.Command)...)
		}
	}
	return args
}

// PlanLayout devolve a linha `wt` que seria disparada para um layout.
func PlanLayout(project config.Project, layout config.Layout, colorize bool) string {
	return "wt " + strings.Join(buildLayoutArgs(project, layout, colorize), " ")
}

// LaunchLayout abre um projeto usando um layout (uma janela com as abas/painéis).
func LaunchLayout(project config.Project, layout config.Layout, colorize, dryRun bool) (string, error) {
	line := PlanLayout(project, layout, colorize)
	if dryRun {
		return line, nil
	}
	if err := exec.Command("wt", buildLayoutArgs(project, layout, colorize)...).Start(); err != nil {
		return line, err
	}
	return line, nil
}
