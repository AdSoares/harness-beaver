package i18n

// Portuguese catalog for the TUI (app, picker, styles).
func init() {
	register(map[string]string{
		// Run() — main menu title
		"HarnessBeaver — menu": "HarnessBeaver — menu",

		// Run() — review offer picker
		"There is an unanalyzed log from %s. Analyze learnings now?": "Há log de %s ainda não analisado. Analisar aprendizados agora?",
		"Yes, analyze": "Sim, analisar",
		"Not now":      "Agora não",

		// mainMenuItems()
		"Open projects":     "Abrir projetos",
		"Open package":      "Abrir pacote",
		"Open with layout":  "Abrir com layout",
		"Add project":       "Adicionar projeto",
		"Remove project":    "Remover projeto",
		"Scan & import":     "Escanear & importar",
		"Create package":    "Criar pacote",
		"Remove package":    "Remover pacote",
		"Execute command":   "Executar comando",
		"History":           "Histórico",
		"Aliases":           "Atalhos",
		"Analyze learnings": "Analisar aprendizados",
		"Settings":          "Configurações",
		"Quit":              "Sair",

		// packageItems()
		"%d project(s)": "%d projeto(s)",

		// modeItems()
		"(use global default)":                   "(usar default global)",
		"Tabs (one window, one tab per project)": "Abas (uma janela, uma aba por projeto)",
		"Separate windows":                       "Janelas separadas",

		// shellItems()
		"(use default)": "(usar default)",

		// onEnter() — scOpenPick
		"Select at least one project (space).": "Selecione ao menos um projeto (espaço).",

		// onEnter() — scOpenConfirm / scOpenMode
		"Open mode":          "Modo de abertura",
		"Shell for each tab": "Shell em cada aba",

		// onEnter() — scPkgOpenPick
		"Package not found.": "Pacote não encontrado.",

		// onEnter() — scProjRemove
		"Project removed: ": "Projeto removido: ",

		// onEnter() — scScan
		"%d project(s) imported.": "%d projeto(s) importado(s).",

		// onEnter() — scPkgAddPick / scPkgAddMode / scPkgAddShell
		"Select at least one project.": "Selecione ao menos um projeto.",
		"Package mode":                 "Modo do pacote",
		"Package shell":                "Shell do pacote",

		// onEnter() — scPkgAddShell
		"Package created: ": "Pacote criado: ",

		// onEnter() — scPkgRemove
		"Package removed: ": "Pacote removido: ",

		// onEnter() — scReviewPick
		"Day %s already has a saved analysis. What do you want?": "Dia %s já tem análise salva. O que deseja?",
		"Today (%s) already has an analysis. What do you want?":  "Hoje (%s) já tem uma análise. O que deseja?",
		"View existing analysis (cache)":                         "Ver análise existente (cache)",
		"Re-analyze now (run Claude again)":                      "Reanalisar agora (rodar o Claude de novo)",

		// onEnter() — scConfigList (Set picker title reuses "Set " key)
		"Set ": "Definir ",

		// onEnter() — scLayoutPick
		"No projects registered.":           "Nenhum projeto cadastrado.",
		"Open with layout — choose project": "Abrir com layout — escolha o projeto",

		// onEnter() — scLayoutProj
		"Project or layout not found.": "Projeto ou layout não encontrado.",
		"Failed to open: ":             "Falha ao abrir: ",
		"Opened %s with layout %s.":    "Aberto %s com layout %s.",

		// analyzeDoneMsg handler
		"Analysis failed — ": "Falha na análise — ",
		"Learnings from ":    "Aprendizados de ",
		"Saved to: ":         "Salvo em: ",

		// applyConfig()
		"Error saving: ": "Erro ao salvar: ",

		// dashTUI()
		"(empty)": "(vazio)",

		// beginAnalyze()
		"Analyzing learnings from ": "Analisando aprendizados de ",

		// showExistingLearning()
		"Analysis of ":                       "Análise de ",
		"Could not read existing analysis: ": "Não foi possível ler a análise existente: ",
		" (cache)":                           " (cache)",

		// dispatchMenu() — "open"
		"No projects registered. Use 'Add project' or 'Scan'.": "Nenhum projeto cadastrado. Use 'Adicionar projeto' ou 'Escanear'.",
		"Open projects (space to select, enter to confirm)":    "Abrir projetos (espaço marca, enter confirma)",

		// dispatchMenu() — "openpkg"
		"No packages registered. Use 'Create package'.": "Nenhum pacote cadastrado. Use 'Criar pacote'.",

		// dispatchMenu() — "addproj"
		"Project name (enter uses folder name)": "Nome do projeto (enter usa o nome da pasta)",

		// dispatchMenu() — "rmproj"
		"No projects to remove.": "Nenhum projeto para remover.",
		// "Remove project" already registered above (mainMenuItems shares same key)

		// dispatchMenu() — "addpkg"
		"Register projects before creating a package.": "Cadastre projetos antes de criar um pacote.",
		"Package name": "Nome do pacote",

		// dispatchMenu() — "rmpkg"
		"No packages to remove.": "Nenhum pacote para remover.",
		// "Remove package" already registered above (mainMenuItems shares same key)

		// dispatchMenu() — "exec"
		"Command to execute (e.g. git status)": "Comando a executar (ex: git status)",

		// dispatchMenu() — "review"
		"No logs yet. Use 'Execute command' to start recording.": "Sem logs ainda. Use 'Executar comando' para começar a registrar.",
		"no analysis":                        "sem análise",
		"today — can re-analyze":             "hoje — pode reanalisar",
		"already analyzed (cache available)": "já analisado (cache disponível)",
		"Analyze learnings — choose day":     "Analisar aprendizados — escolha o dia",

		// dispatchMenu() — "openlayout"
		"No layouts. Create one with 'bvr layout add <name> --preset dev'.": "Nenhum layout. Crie com 'bvr layout add <nome> --preset dev'.",
		"%d tab(s), %d pane(s)":            "%d aba(s), %d painel(éis)",
		"Open with layout — choose layout": "Abrir com layout — escolha o layout",

		// openHistory()
		"No history in the last 7 days.":     "Sem histórico nos últimos 7 dias.",
		"History — enter re-runs · / filter": "Histórico — enter re-executa · / filtra",

		// openAliases()
		"No aliases. Create one with 'bvr alias set <name> <command>'.": "Nenhum atalho. Crie com 'bvr alias set <nome> <comando>'.",
		"Aliases — enter runs · / filter":                               "Atalhos — enter executa · / filtra",

		// openConfigList()
		"Settings — choose a key": "Configurações — escolha uma chave",

		// startScan()
		"Error scanning: ":            "Erro ao escanear: ",
		"Nothing new to import in %s": "Nada novo para importar em %s",
		"Scan & import — %s":          "Escanear & importar — %s",

		// updateInput() — scProjAddName
		"Directory path (e.g. F:\\proj\\my-app)": "Caminho do diretório (ex: F:\\proj\\meu-app)",
		"Please enter a path.":                   "Informe um caminho.",
		"Project added: ":                        "Projeto adicionado: ",

		// updateInput() — scPkgAddName
		"Please enter a name.":               "Informe um nome.",
		"Package projects (space to select)": "Projetos do pacote (espaço marca)",

		// updateInput() — scExecInput
		"Please enter a command.": "Informe um comando.",

		// launch()
		"Opened %d project(s) in %s mode. (%d wt command)": "Aberto %d projeto(s) em modo %s. (%d comando wt)",

		// save()
		// "Error saving: " already registered above (applyConfig shares the same English key)

		// formatExecOutput()
		"(no output)": "(sem saída)",

		// lastLines()
		"…(%d previous lines omitted)": "…(%d linhas anteriores omitidas)",

		// openConfirmView()
		"Confirm projects to open": "Confirmar projetos a abrir",
		"  … and %d more":          "  … e mais %d",
		"%d project(s) selected":   "%d projeto(s) selecionado(s)",

		// View() — spinner
		" calling Claude, please wait…\n": " chamando o Claude, aguarde…\n",

		// inputTitle()
		"Add project — name":    "Adicionar projeto — nome",
		"Add project — path":    "Adicionar projeto — caminho",
		"Create package — name": "Criar pacote — nome",
		// "Execute command" already registered above (mainMenuItems shares same key)
		// "Set " already registered above

		// helpLine()
		"enter: confirm · esc: cancel · ctrl+c: quit":                         "enter: confirmar · esc: cancelar · ctrl+c: sair",
		"analyzing… (ctrl+c quits)":                                           "analisando… (ctrl+c sai)",
		"enter/esc: back to menu":                                             "enter/esc: voltar ao menu",
		"enter: open · esc: back to selection":                                "enter: abrir · esc: voltar à seleção",
		"filtering: type · ↑/↓ navigate · enter apply · esc clear":            "filtrando: digite · ↑/↓ navega · enter aplica · esc limpa",
		"↑/↓ navigate · / filter · enter select · q quit":                     "↑/↓ navegar · / filtrar · enter selecionar · q sair",
		"↑/↓ nav · ←/→ pages · / filter · space select · enter ok · esc back": "↑/↓ nav · ←/→ págs · / filtrar · espaço marcar · enter ok · esc voltar",
		"↑/↓ nav · ←/→ pages · / filter · enter select · esc back":            "↑/↓ nav · ←/→ págs · / filtrar · enter selecionar · esc voltar",

		// picker.go — view()
		"filter: ":                    "filtro: ",
		"(no items match the filter)": "(nenhum item corresponde ao filtro)",
		// "(empty)" already registered above
		"  ▲ %d above":        "  ▲ %d acima",
		"  ▼ %d below":        "  ▼ %d abaixo",
		"  %d–%d of %d":       "  %d–%d de %d",
		" (filtered from %d)": " (filtrados de %d)",
		" · %d selected":      " · %d marcado(s)",
	})
}
