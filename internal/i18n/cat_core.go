package i18n

// Portuguese catalog for the root command and cross-cutting strings.
// Keys are the exact English source strings passed to T().
func init() {
	register(map[string]string{
		"HarnessBeaver — multi-directory Claude Code launcher": "HarnessBeaver — launcher de Claude Code multi-diretório",
		"HarnessBeaver (bvr) opens Claude Code (or pwsh/cmd) across many\ndirectories at once, in Windows Terminal tabs or separate windows.\n\nWith no arguments it opens the interactive interface (TUI). With subcommands it\nworks as a scriptable CLI.": "HarnessBeaver (bvr) abre o Claude Code (ou pwsh/cmd) em vários\ndiretórios de uma vez, em abas do Windows Terminal ou janelas separadas.\n\nSem argumentos, abre a interface interativa (TUI). Com subcomandos, funciona\ncomo CLI scriptável.",
		"error:":                    "erro:",
		"interface language: en|pt": "idioma da interface: en|pt",
	})
}
