package i18n

// Portuguese catalog for cmd batch A (add, alias, config, do, each, edit, helpers).
func init() {
	register(map[string]string{
		// add.go
		"add <name> <path>":                      "add <nome> <path>",
		"Register a new project (directory)":     "Registra um novo projeto (diretório)",
		"Project added: %s (%s)\n":               "Projeto adicionado: %s (%s)\n",
		"project default shell: claude|pwsh|cmd": "shell padrão do projeto: claude|pwsh|cmd",

		// alias.go
		"Manage command shortcuts (use 'bvr do <name>' to run)":            "Gerencia atalhos de comandos (use 'bvr do <nome>' para executar)",
		"No shortcuts. Create one with: bvr alias set <name> <command...>": "Nenhum atalho. Crie com: bvr alias set <nome> <comando...>",
		"set <name> <command...>":                                          "set <nome> <comando...>",
		"Create or update a shortcut":                                      "Cria ou atualiza um atalho",
		"rm <name>":                                                        "rm <nome>",
		"Remove a shortcut":                                                "Remove um atalho",
		"shortcut not found: %s":                                           "atalho não encontrado: %s",
		"Shortcut removed: %s\n":                                           "Atalho removido: %s\n",

		// config.go
		"config [get <key> | set <key> <value>]":                    "config [get <chave> | set <chave> <valor>]",
		"Show or change bvr settings":                               "Mostra ou altera as configurações do bvr",
		"usage: bvr config get <key>":                               "uso: bvr config get <chave>",
		"unknown key %q (valid: %s)":                                "chave desconhecida %q (válidas: %s)",
		"usage: bvr config set <key> <value>":                       "uso: bvr config set <chave> <valor>",
		"invalid subcommand %q (use: get, set, or nothing to list)": "subcomando inválido %q (use: get, set, ou nada para listar)",

		// do.go
		"do <name> [args...]":                                "do <nome> [args...]",
		"Run a registered alias, expanding {1}..{N} and {*}": "Executa um atalho (alias) registrado, expandindo {1}..{N} e {*}",
		"alias not found: %s (see 'bvr alias')":              "atalho não encontrado: %s (veja 'bvr alias')",
		"execution shell: pwsh|cmd|bash|zsh":                 "shell de execução: pwsh|cmd|bash|zsh",

		// each.go
		"each <pkgId|projId...> -- <command...>":                                       "each <pkgId|projId...> -- <comando...>",
		"Run a command in each project (from a package or ids)":                        "Roda um comando em cada projeto (de um pacote ou ids)",
		"use '--' to separate targets from command: bvr each <ids...> -- <command...>": "use '--' para separar alvos do comando: bvr each <ids...> -- <comando...>",
		"provide at least one project/package before '--'":                             "informe ao menos um projeto/pacote antes de '--'",
		"provide the command after '--'":                                               "informe o comando após '--'",
		"No target projects.":                                                          "Nenhum projeto alvo.",
		"\nSummary: %d project(s), %d with error (exit != 0).\n":                       "\nResumo: %d projeto(s), %d com erro (exit != 0).\n",
		"\n=== %s · %s ===\n":                                                          "\n=== %s · %s ===\n",
		"(exit %d, %dms)\n":                                                            "(exit %d, %dms)\n",
		"\n=== %s · %s === (exit %d, %dms)\n":                                          "\n=== %s · %s === (exit %d, %dms)\n",
		"run in parallel (output captured per project)":                                "roda em paralelo (saída capturada por projeto)",

		// edit.go
		"edit <projId>": "edit <projId>",
		"Change name, path or shell of a project": "Altera nome, path ou shell de um projeto",
		"project not found: %s":                   "projeto não encontrado: %s",
		"invalid path: %s":                        "path inválido: %s",
		"Project updated: %s\n":                   "Projeto atualizado: %s\n",
		"new name":                                "novo nome",
		"new path":                                "novo caminho",
		"new default shell: claude|pwsh|cmd":      "novo shell padrão: claude|pwsh|cmd",

		// helpers.go
		"invalid shell %q (use: claude, pwsh, cmd)": "shell inválido %q (use: claude, pwsh, cmd)",
		"invalid mode %q (use: tabs, windows)":      "modo inválido %q (use: tabs, windows)",
		"id not found (project or package): %s":     "id não encontrado (projeto ou pacote): %s",
		"no projects to open":                       "nenhum projeto para abrir",
	})

	// Long descriptions use raw string literals (real newlines), registered separately.
	register(map[string]string{
		`Run the alias command <name>. Extra arguments replace the
placeholders {1}, {2}, … and {*} (all together); without placeholders, they are appended.

Examples:
  bvr alias set deploy "npm run build && npm run deploy"
  bvr do deploy
  bvr alias set commit "git commit -m {1}"
  bvr do commit "fix: adjustment"`: `Executa o comando do atalho <nome>. Os argumentos extras substituem os
placeholders {1}, {2}, … e {*} (todos juntos); sem placeholders, são anexados.

Exemplos:
  bvr alias set deploy "npm run build && npm run deploy"
  bvr do deploy
  bvr alias set commit "git commit -m {1}"
  bvr do commit "fix: ajuste"`,

		`Run the same command in each target project's directory.
Targets (project and/or package ids) come before '--' and the command comes after.
Each run is recorded in the journal (with the project's cwd).

Examples:
  bvr each smb-ativo -- git pull
  bvr each repair beauty --shell bash -- git status -s
  bvr each smb-ativo --parallel -- git fetch`: `Executa o mesmo comando no diretório de cada projeto alvo.
Os alvos (ids de projeto e/ou pacote) vêm antes de '--' e o comando vem depois.
Cada execução é registrada no diário (com o cwd do projeto).

Exemplos:
  bvr each smb-ativo -- git pull
  bvr each repair beauty --shell bash -- git status -s
  bvr each smb-ativo --parallel -- git fetch`,
	})
}
