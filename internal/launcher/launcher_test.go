package launcher

import (
	"reflect"
	"strings"
	"testing"

	"harnessbeaver/internal/config"
)

func TestPaneCommand(t *testing.T) {
	cases := []struct {
		shell   config.Shell
		command string
		want    []string
	}{
		{config.ShellClaude, "", []string{"pwsh", "-NoExit", "-Command", "claude"}},
		{config.ShellPwsh, "", []string{"pwsh", "-NoExit"}},
		{config.ShellPwsh, "npm run dev", []string{"pwsh", "-NoExit", "-Command", "npm run dev"}},
		{config.ShellCmd, "dir", []string{"cmd", "/k", "dir"}},
		{config.ShellBash, "ls", []string{"bash", "-lc", "ls; exec bash"}},
	}
	for _, c := range cases {
		if got := paneCommand(c.shell, c.command); !reflect.DeepEqual(got, c.want) {
			t.Errorf("paneCommand(%q,%q) = %v, esperava %v", c.shell, c.command, got, c.want)
		}
	}
}

func TestBuildLayoutArgs(t *testing.T) {
	project := config.Project{Name: "repair", Path: `C:\p\repair`}
	layout := config.Layout{
		Name: "dev",
		Tabs: []config.Tab{
			{Panes: []config.Pane{
				{Shell: config.ShellClaude},
				{Shell: config.ShellPwsh, Split: "V", Command: "npm run dev"},
			}},
			{Title: "git", Panes: []config.Pane{
				{Shell: config.ShellPwsh, Command: "git status"},
			}},
		},
	}
	got := buildLayoutArgs(project, layout, false)
	want := []string{
		"new-tab", "--title", "repair", "--suppressApplicationTitle", "-d", `C:\p\repair`, "pwsh", "-NoExit", "-Command", "claude",
		";",
		"split-pane", "-V", "--title", "repair", "--suppressApplicationTitle", "-d", `C:\p\repair`, "pwsh", "-NoExit", "-Command", "npm run dev",
		";",
		"new-tab", "--title", "git", "--suppressApplicationTitle", "-d", `C:\p\repair`, "pwsh", "-NoExit", "-Command", "git status",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildLayoutArgs =\n%v\nesperava\n%v", got, want)
	}
}

func TestBuildLayoutArgsUsesPaneDir(t *testing.T) {
	project := config.Project{Name: "p", Path: `C:\proj`}
	layout := config.Layout{Tabs: []config.Tab{{Panes: []config.Pane{
		{Shell: config.ShellPwsh, Dir: `C:\outro`},
	}}}}
	line := PlanLayout(project, layout, false)
	if !strings.Contains(line, `-d C:\outro`) {
		t.Errorf("esperava cwd do painel no comando: %s", line)
	}
}

func TestColorize(t *testing.T) {
	items := []Item{
		{Project: config.Project{Name: "a", Path: `C:\a`}, Shell: config.ShellPwsh},
		{Project: config.Project{Name: "b", Path: `C:\b`}, Shell: config.ShellPwsh},
	}
	// Sem cor: nenhum --tabColor.
	if got := strings.Join(buildTabsArgs(items, false), " "); strings.Contains(got, "--tabColor") {
		t.Errorf("colorize=false não deveria conter --tabColor: %s", got)
	}
	// Com cor: cada aba recebe uma cor distinta da paleta.
	got := buildTabsArgs(items, true)
	joined := strings.Join(got, " ")
	if !strings.Contains(joined, "--tabColor "+colorFor(0)) || !strings.Contains(joined, "--tabColor "+colorFor(1)) {
		t.Errorf("esperava cores distintas por aba: %s", joined)
	}
	if colorFor(0) == colorFor(1) {
		t.Error("as duas primeiras cores da paleta deveriam diferir")
	}
}
