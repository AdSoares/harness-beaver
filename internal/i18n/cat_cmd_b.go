package i18n

// Portuguese catalog for cmd batch B (history, layout, list, open, pkg, remove).
func init() {
	register(map[string]string{
		// history.go
		"List command history from the journal (with filters)": "Lista o histórico de comandos do diário (com filtros)",
		`Shows commands recorded in the last --days days, from most recent
to oldest. Filter by project, text, or only those that failed.`: `Mostra os comandos registrados nos últimos --days dias, do mais recente
para o mais antigo. Filtre por projeto, texto ou apenas os que falharam.`,
		"No history entries for the given filters.": "Nenhuma entrada no histórico para os filtros dados.",
		"filter by project id":                      "filtra por id de projeto",
		"filter by text in command":                 "filtra por texto no comando",
		"only commands with exit != 0":              "apenas comandos com exit != 0",
		"how many days back":                        "quantos dias para trás",
		"maximum number of entries shown":           "máximo de entradas exibidas",

		// layout.go
		"Manage layouts (tab/pane arrangements) for opening projects":              "Gerencia layouts (arranjos de abas/painéis) para abrir projetos",
		"No layouts. Create one with: bvr layout add <name> [--preset dev|triple]": "Nenhum layout. Crie com: bvr layout add <nome> [--preset dev|triple]",
		"  %-16s %-24s %d tab(s), %d pane(s)\n":                                    "  %-16s %-24s %d aba(s), %d painel(éis)\n",
		"Create a layout (preset: claude|dev|triple)":                              "Cria um layout (preset: claude|dev|triple)",
		"Layout created: %s\n":                                                     "Layout criado: %s\n",
		"tab <layoutId>":                                                           "tab <layoutId>",
		"Add a tab to the layout":                                                  "Adiciona uma aba ao layout",
		"layout not found: %s":                                                     "layout não encontrado: %s",
		"Tab %d added to layout %s.\n":                                             "Aba %d adicionada ao layout %s.\n",
		"pane <layoutId>":                                                          "pane <layoutId>",
		"Add a pane to a tab in the layout":                                        "Adiciona um painel a uma aba do layout",
		"invalid tab %d (the layout has %d tab(s))":                                "aba %d inválida (o layout tem %d aba(s))",
		"invalid split %q (use H or V)":                                            "split inválido %q (use H ou V)",
		"Pane added to tab %d of layout %s.\n":                                     "Painel adicionado à aba %d do layout %s.\n",
		"rm <layoutId>":                                                            "rm <layoutId>",
		"Remove a layout":                                                          "Remove um layout",
		"Layout removed: %s\n":                                                     "Layout removido: %s\n",
		"assign <layoutId> <projId>":                                               "assign <layoutId> <projId>",
		"Set the default layout for a project":                                     "Define o layout padrão de um projeto",
		"project not found: %s":                                                    "projeto não encontrado: %s",
		"Project %s will now open with layout %s.\n":                               "Projeto %s agora abre com o layout %s.\n",
		"show <layoutId>":                                                          "show <layoutId>",
		"Show the wt command that the layout would generate":                       "Mostra o comando wt que o layout geraria",
		"invalid shell %q (use: claude, pwsh, cmd, bash, zsh)":                     "shell inválido %q (use: claude, pwsh, cmd, bash, zsh)",
		"initial preset: claude|dev|triple":                                        "preset inicial: claude|dev|triple",
		"tab title":                                                                "título da aba",
		"tab number (1-based)":                                                     "número da aba (1-based)",
		"pane shell: claude|pwsh|cmd|bash|zsh":                                     "shell do painel: claude|pwsh|cmd|bash|zsh",
		"command to run in the pane":                                               "comando a rodar no painel",
		"split from previous pane: H|V":                                            "divisão a partir do painel anterior: H|V",
		"pane cwd (empty = project dir)":                                           "cwd do painel (vazio = dir do projeto)",
		"project to fill in paths":                                                 "projeto p/ preencher os caminhos",

		// list.go
		"List registered projects and/or packages": "Lista projetos e/ou pacotes registrados",
		"Projects:": "Projetos:",
		"  (none)":  "  (nenhum)",
		"Packages:": "Pacotes:",
		"  %-20s %-28s mode=%s shell=%s projects=[%s]\n": "  %-20s %-28s mode=%s shell=%s projetos=[%s]\n",

		// open.go
		"open <pkgId|projId> [others...]":                 "open <pkgId|projId> [outros...]",
		"Open projects or a package in Windows Terminal":  "Abre projetos ou um pacote no Windows Terminal",
		"Opened %s with layout %s.\n":                     "Aberto %s com o layout %s.\n",
		"Opened %d project(s) in mode %s.\n":              "Aberto %d projeto(s) em modo %s.\n",
		"--layout requires exactly one project":           "--layout requer exatamente um projeto",
		"--layout requires a project (not a package): %s": "--layout requer um projeto (não um pacote): %s",
		"open mode: tabs|windows":                         "modo de abertura: tabs|windows",
		"shell per tab: claude|pwsh|cmd":                  "shell por aba: claude|pwsh|cmd",
		"print wt commands without executing":             "imprime os comandos wt sem executar",
		"open the project with this layout":               "abre o projeto com este layout",
		"ignore the project's default layout":             "ignora o layout padrão do projeto",

		// pkg.go — "add <name>" and "project not found: %s" are shared with layout.go above
		"Manage packages (groups of projects)":     "Gerencia pacotes (grupos de projetos)",
		"add <name>":                               "add <nome>",
		"Create a package with the given projects": "Cria um pacote com os projetos indicados",
		"provide --projects id1,id2,...":           "informe --projects id1,id2,...",
		"Package created: %s (%d projects)\n":      "Pacote criado: %s (%d projetos)\n",
		"remove <pkgId>":                           "remove <pkgId>",
		"Remove a package":                         "Remove um pacote",
		"package not found: %s":                    "pacote não encontrado: %s",
		"Package removed: %s\n":                    "Pacote removido: %s\n",
		"comma-separated project ids":              "ids de projetos separados por vírgula",
		"package mode: tabs|windows":               "modo do pacote: tabs|windows",
		"package shell: claude|pwsh|cmd":           "shell do pacote: claude|pwsh|cmd",

		// remove.go — "project not found: %s" is shared with layout.go above
		"remove <projId>":                    "remove <projId>",
		"Remove a project from the registry": "Remove um projeto do registro",
		"Project removed: %s\n":              "Projeto removido: %s\n",
	})
}
