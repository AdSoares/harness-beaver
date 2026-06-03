package i18n

// Portuguese catalog for cmd batch C (review, run, scan, shell, status, upgrade, version, ui).
func init() {
	register(map[string]string{
		// review.go
		"review [date]":                  "review [data]",
		"Analyse a day's learnings (AI)": "Analisa os aprendizados de um dia (IA)",
		"Analyse the commands recorded on a given day and generate a markdown\nlearnings file at ~/.harnessbeaver/learnings/<date>.md.\n\nWith no argument, uses the pending day (previous, not yet analysed) or today.\nDate in YYYY-MM-DD format. Use --dry-run to see the prompt without calling the AI.": "Analisa os comandos registrados em um dia e gera um markdown de\naprendizados em ~/.harnessbeaver/learnings/<data>.md.\n\nSem argumento, usa o dia pendente (anterior, sem análise) ou hoje.\nData no formato AAAA-MM-DD. Use --dry-run para ver o prompt sem chamar a IA.",
		"invalid date %q (use YYYY-MM-DD)":                        "data inválida %q (use AAAA-MM-DD)",
		"Existing analysis for %s (use --force to redo):\n\n%s\n": "Análise existente de %s (use --force para refazer):\n\n%s\n",
		"Analysis saved to: %s\n\n%s\n":                           "Análise salva em: %s\n\n%s\n",
		"print the prompt without calling the AI":                 "imprime o prompt sem chamar a IA",
		"redo analysis even if a cached result exists":            "refaz a análise mesmo se já houver cache",

		// run.go
		"run <command...>": "run <comando...>",
		"Execute a command (pwsh/cmd) and record it in the journal": "Executa um comando (pwsh/cmd) e registra no diário",
		"Execute a command on behalf of the user in PowerShell (default) or cmd and\nrecord the command + output in the journal (~/.harnessbeaver/logs).\n\nFlags must come before the command. Examples:\n  bvr run echo hello\n  bvr run --shell cmd \"echo hi & dir\"": "Executa um comando em nome do usuário no PowerShell (default) ou cmd e\nregistra comando + saída no diário de bordo (~/.harnessbeaver/logs).\n\nFlags devem vir antes do comando. Exemplos:\n  bvr run echo ola\n  bvr run --shell cmd \"echo oi & dir\"",
		"Hint: there is an unanalysed log for %s — run 'bvr review'.\n": "Dica: há log de %s não analisado — rode 'bvr review'.\n",
		"invalid shell %q (use: pwsh, cmd, bash, zsh)":                  "shell inválido %q (use: pwsh, cmd, bash, zsh)",
		"execution shell: pwsh|cmd|bash|zsh":                            "shell de execução: pwsh|cmd|bash|zsh",

		// scan.go
		"Scan a directory for code projects": "Escaneia um diretório em busca de projetos de código",
		"Scan the root directory (--root, or settings.scanRoot, or the current\ndirectory) and list subdirectories that look like code projects. With --import,\nrecords all candidates not yet registered.": "Escaneia o diretório raiz (--root, ou settings.scanRoot, ou o diretório\natual) e lista subdiretórios que parecem projetos de código. Com --import, grava\nno registro todos os candidatos ainda não cadastrados.",
		"No code projects found in %s\n":                       "Nenhum projeto de código encontrado em %s\n",
		" (already registered)":                                " (já cadastrado)",
		"\n%d project(s) imported.\n":                          "\n%d projeto(s) importado(s).\n",
		"\n%d candidate(s). Use --import to register.\n":       "\n%d candidato(s). Use --import para gravar.\n",
		"root directory to scan":                               "diretório raiz a escanear",
		"import new candidates":                                "importa os candidatos novos",
		"also offer parent directories that group ≥2 projects": "também oferece diretórios-pai que agrupam ≥2 projetos",
		"maximum depth":                                        "profundidade máxima",

		// shell.go
		"REPL that executes and records commands (pwsh/cmd)": "REPL que executa e registra comandos (pwsh/cmd)",
		"Open an interactive session: each line typed is executed in the chosen\nshell and recorded in the journal. Type 'exit' or 'quit' (or Ctrl+D) to leave.": "Abre uma sessão interativa: cada linha digitada é executada no shell\nescolhido e registrada no diário. Digite 'exit' ou 'quit' (ou Ctrl+D) para sair.",
		"bvr shell — %s. Commands are recorded in the journal. 'exit' to quit.\n":                                                                                "bvr shell — %s. Comandos são registrados no diário. 'exit' para sair.\n",
		"There is an unanalysed log for %s. Analyse now? [y/N] ":                                                                                                 "Há log de %s ainda não analisado. Analisar agora? [s/N] ",
		"analysis failed:":        "falha na análise:",
		"Analysis saved to: %s\n": "Análise salva em: %s\n",

		// status.go
		"Show git status of projects (branch, changes, ahead/behind)":                                                                                        "Mostra o status git dos projetos (branch, alterações, ahead/behind)",
		"With no arguments, shows the git status of all registered projects.\nWith project and/or package ids, restricts to those. Queries run in parallel.": "Sem argumentos, mostra o status git de todos os projetos registrados.\nCom ids de projeto e/ou pacote, restringe a esses. As consultas rodam em paralelo.",
		"No registered projects.":               "Nenhum projeto registrado.",
		"(not a git repository)":                "(não é repositório git)",
		"id not found (project or package): %s": "id não encontrado (projeto ou pacote): %s",

		// upgrade.go
		"Update the bvr binary from a source (URL or path)": "Atualiza o binário do bvr a partir de uma fonte (URL ou caminho)",
		"Download the new binary from --source (or settings.upgradeSource) and replace\nthe current executable. The source can be an http(s) URL or a local/shared file\npath. The previous binary is kept as a backup (.old).": "Baixa o binário novo de --source (ou settings.upgradeSource) e substitui o\nexecutável atual. A fonte pode ser uma URL http(s) ou um caminho de arquivo\nlocal/compartilhado. O binário anterior é mantido como backup (.old).",
		"no source: use --source <url|path> or 'bvr config set upgradeSource <...>'": "nenhuma fonte: use --source <url|caminho> ou 'bvr config set upgradeSource <...>'",
		"failed to fetch binary: %w":                                            "falha ao obter binário: %w",
		"could not move the current binary: %w":                                 "não foi possível mover o binário atual: %w",
		"could not install the new binary: %w":                                  "não foi possível instalar o novo binário: %w",
		"Updated: %s\n(backup of the previous binary at %s — safe to delete)\n": "Atualizado: %s\n(backup do anterior em %s — pode apagar)\n",
		"URL or path of the new binary":                                         "URL ou caminho do binário novo",

		// version.go
		"Show the bvr version": "Mostra a versão do bvr",

		// ui.go
		"Open the interactive interface (TUI)": "Abre a interface interativa (TUI)",
	})
}
