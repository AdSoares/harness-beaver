// Package i18n provides minimal runtime localization for bvr.
//
// English is the canonical in-code language (and the default at runtime).
// Portuguese is supplied by per-area catalogs that map the English source
// string to its translation. Wrap any user-facing string with T():
//
//	fmt.Println(i18n.T("Project added"))
//	fmt.Printf(i18n.T("Project added: %s (%s)\n"), id, path)
//
// For Printf-style strings, T() the format string and keep the verbs; the
// Portuguese catalog entry must carry the same verbs in the same order.
package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Lang is a supported language code.
type Lang string

const (
	En Lang = "en"
	Pt Lang = "pt"
)

var current = En

// ptCatalog maps English source strings to their Portuguese translations.
// It is populated by register() calls from the per-area catalog files.
var ptCatalog = map[string]string{}

// register merges a batch of English->Portuguese entries into the catalog.
// Called from init() in catalog files so areas stay conflict-free.
func register(m map[string]string) {
	for k, v := range m {
		ptCatalog[k] = v
	}
}

// init resolves the language as early as possible, before cobra command
// descriptions are built during package initialization, so help text is also
// localized. Precedence here: config file setting, then BVR_LANG env (env wins).
// A --lang flag, applied later at runtime, overrides both for command output.
func init() {
	if l := languageFromConfig(); l != "" {
		Set(l)
	}
	if v := os.Getenv("BVR_LANG"); v != "" {
		Set(v)
	}
}

// Set switches the active language. It accepts several spellings
// (e.g. "pt", "pt-BR", "portugues", "en", "english"); unknown values are ignored.
func Set(s string) {
	switch normalize(s) {
	case "pt", "ptbr", "br", "portugues", "portuguese":
		current = Pt
	case "en", "enus", "english", "ingles":
		current = En
	}
}

// Current returns the active language code.
func Current() Lang { return current }

// T returns the localized version of an English source string. For the default
// language (English) it returns the input unchanged; for Portuguese it returns
// the catalog entry, falling back to the English source when none exists.
func T(s string) string {
	if current == Pt {
		if v, ok := ptCatalog[s]; ok {
			return v
		}
	}
	return s
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, r := range []string{"-", "_", " "} {
		s = strings.ReplaceAll(s, r, "")
	}
	s = strings.ReplaceAll(s, "ê", "e")
	s = strings.ReplaceAll(s, "ç", "c")
	s = strings.ReplaceAll(s, "ã", "a")
	return s
}

// languageFromConfig reads only settings.language from the config file without
// importing the config package (which would create an import cycle).
func languageFromConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(home, ".harnessbeaver", "config.json"))
	if err != nil {
		return ""
	}
	var c struct {
		Settings struct {
			Language string `json:"language"`
		} `json:"settings"`
	}
	if json.Unmarshal(data, &c) != nil {
		return ""
	}
	return c.Settings.Language
}
