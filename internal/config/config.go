// Package config define o modelo de dados do HarnessBeaver e cuida da
// persistência em %USERPROFILE%\.harnessbeaver\config.json.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"harnessbeaver/internal/i18n"
)

// Shell é o programa que roda dentro de cada aba/janela aberta.
type Shell string

// Mode define como os projetos selecionados são abertos.
type Mode string

const (
	ShellClaude Shell = "claude"
	ShellPwsh   Shell = "pwsh"
	ShellCmd    Shell = "cmd"
	ShellBash   Shell = "bash"
	ShellZsh    Shell = "zsh"

	ModeTabs    Mode = "tabs"    // uma janela do Windows Terminal, uma aba por projeto
	ModeWindows Mode = "windows" // uma janela separada por projeto
)

// ValidShells (launcher: claude/pwsh/cmd via wt), ValidRunShells (runner:
// pwsh/cmd/bash/zsh) e ValidModes servem para validação de flags na CLI.
var (
	ValidShells    = []Shell{ShellClaude, ShellPwsh, ShellCmd}
	ValidRunShells = []Shell{ShellPwsh, ShellCmd, ShellBash, ShellZsh}
	ValidModes     = []Mode{ModeTabs, ModeWindows}
)

// DefaultRunShellForOS devolve o shell de execução padrão do sistema operacional.
func DefaultRunShellForOS() Shell {
	if runtime.GOOS == "windows" {
		return ShellPwsh
	}
	return ShellBash
}

// Project é um diretório registrado que pode ser aberto.
type Project struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	DefaultShell Shell    `json:"defaultShell,omitempty"`
	LayoutID     string   `json:"layoutId,omitempty"` // layout padrão ao abrir (Fase 3)
	Tags         []string `json:"tags,omitempty"`
}

// Pane é um painel dentro de uma aba: um shell, com comando e cwd opcionais.
// Split indica como o painel é criado a partir do anterior: "" = primeira aba
// (new-tab), "H" = divisão horizontal, "V" = divisão vertical.
type Pane struct {
	Shell   Shell  `json:"shell,omitempty"`
	Command string `json:"command,omitempty"`
	Dir     string `json:"dir,omitempty"` // vazio = diretório do projeto
	Split   string `json:"split,omitempty"`
}

// Tab é uma aba do Windows Terminal com um ou mais painéis.
type Tab struct {
	Title string `json:"title,omitempty"`
	Panes []Pane `json:"panes"`
}

// Layout é um arranjo nomeado de abas/painéis para abrir um projeto.
type Layout struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tabs []Tab  `json:"tabs"`
}

// Package é um grupo nomeado de projetos com modo/shell opcionais.
type Package struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	ProjectIDs []string `json:"projectIds"`
	Mode       Mode     `json:"mode,omitempty"`
	Shell      Shell    `json:"shell,omitempty"`
}

// Settings guarda os defaults globais.
type Settings struct {
	ScanRoot     string `json:"scanRoot"`
	DefaultMode  Mode   `json:"defaultMode"`
	DefaultShell Shell  `json:"defaultShell"`

	// Runner / aprendizados (Feature 2)
	DefaultRunShell Shell  `json:"defaultRunShell"` // shell do run/shell (default OS-aware)
	InsightsEngine  string `json:"insightsEngine"`  // "auto" | "claude" | "api"
	InsightsModel   string `json:"insightsModel"`   // modelo p/ o caminho da API
	LastReviewOffer string `json:"lastReviewOffer"` // data da última auto-oferta (AAAA-MM-DD)

	// Fundações (Fase 0)
	LearningsExtraDir string `json:"learningsExtraDir"` // ""=desligado; cópia extra das análises
	UpgradeSource     string `json:"upgradeSource"`     // URL ou caminho de onde baixar o binário novo

	TabColors bool `json:"tabColors"` // cor diferente por aba/janela aberta

	Language string `json:"language,omitempty"` // idioma da interface: "en" (default) | "pt"
}

// Config é a raiz serializada para o config.json.
type Config struct {
	Version  int               `json:"version"`
	Settings Settings          `json:"settings"`
	Projects []Project         `json:"projects"`
	Packages []Package         `json:"packages"`
	Aliases  map[string]string `json:"aliases,omitempty"` // nome -> comando (Fase 1)
	Layouts  []Layout          `json:"layouts,omitempty"` // arranjos de abas/painéis (Fase 3)
}

const currentVersion = 1

// Dir retorna a pasta de configuração (~/.harnessbeaver).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".harnessbeaver"), nil
}

// Path retorna o caminho completo do config.json.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// LogsDir retorna a pasta dos logs diários (~/.harnessbeaver/logs).
func LogsDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "logs"), nil
}

// LearningsDir retorna a pasta das análises (~/.harnessbeaver/learnings).
func LearningsDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "learnings"), nil
}

// Default devolve uma configuração inicial vazia.
func Default() *Config {
	return &Config{
		Version: currentVersion,
		Settings: Settings{
			ScanRoot:        "",
			DefaultMode:     ModeTabs,
			DefaultShell:    ShellClaude,
			DefaultRunShell: DefaultRunShellForOS(),
			InsightsEngine:  "auto",
			InsightsModel:   "claude-sonnet-4-6",
			Language:        "en",
		},
		Projects: []Project{},
		Packages: []Package{},
		Aliases:  map[string]string{},
	}
}

// Load lê o config.json; se não existir, cria um default e o grava.
func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		cfg := Default()
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf(i18n.T("invalid config.json: %w"), err)
	}
	if cfg.Version == 0 {
		cfg.Version = currentVersion
	}
	if cfg.Projects == nil {
		cfg.Projects = []Project{}
	}
	if cfg.Packages == nil {
		cfg.Packages = []Package{}
	}
	// Backfill de defaults novos (configs criados antes da Feature 2).
	if cfg.Settings.DefaultRunShell == "" {
		cfg.Settings.DefaultRunShell = DefaultRunShellForOS()
	}
	if cfg.Settings.InsightsEngine == "" {
		cfg.Settings.InsightsEngine = "auto"
	}
	if cfg.Settings.InsightsModel == "" {
		cfg.Settings.InsightsModel = "claude-sonnet-4-6"
	}
	if cfg.Settings.Language == "" {
		cfg.Settings.Language = "en"
	}
	if cfg.Aliases == nil {
		cfg.Aliases = map[string]string{}
	}
	return &cfg, nil
}

// Save grava o config.json (criando a pasta se necessário).
func (c *Config) Save() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	p := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// --- Projetos ---

func (c *Config) FindProject(id string) (*Project, bool) {
	for i := range c.Projects {
		if c.Projects[i].ID == id {
			return &c.Projects[i], true
		}
	}
	return nil, false
}

func (c *Config) hasProjectAtPath(path string) bool {
	clean := filepath.Clean(path)
	for _, p := range c.Projects {
		if filepath.Clean(p.Path) == clean {
			return true
		}
	}
	return false
}

// AddProject valida o caminho, gera um id único e adiciona o projeto.
func (c *Config) AddProject(name, path string, shell Shell) (*Project, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("invalid path: %w"), err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf(i18n.T("path is not a directory: %s"), abs)
	}
	if name == "" {
		name = filepath.Base(abs)
	}
	id := c.uniqueProjectID(name)
	p := Project{ID: id, Name: name, Path: abs, DefaultShell: shell}
	c.Projects = append(c.Projects, p)
	return &c.Projects[len(c.Projects)-1], nil
}

// RemoveProject remove o projeto e o desreferencia de todos os pacotes.
func (c *Config) RemoveProject(id string) bool {
	idx := -1
	for i := range c.Projects {
		if c.Projects[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	c.Projects = append(c.Projects[:idx], c.Projects[idx+1:]...)
	for i := range c.Packages {
		c.Packages[i].ProjectIDs = removeString(c.Packages[i].ProjectIDs, id)
	}
	return true
}

func (c *Config) uniqueProjectID(name string) string {
	return uniqueID(slug(name), func(id string) bool {
		_, ok := c.FindProject(id)
		return ok
	})
}

// ProjectsByIDs devolve os projetos correspondentes, preservando a ordem dos ids.
func (c *Config) ProjectsByIDs(ids []string) []Project {
	var out []Project
	for _, id := range ids {
		if p, ok := c.FindProject(id); ok {
			out = append(out, *p)
		}
	}
	return out
}

// ProjectIDForPath devolve o id do projeto cujo Path contém o caminho dado
// (match por prefixo mais longo), ou "" se nenhum. Usado para taggear o diário.
func (c *Config) ProjectIDForPath(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		abs = p
	}
	abs = filepath.Clean(abs)
	best, bestLen := "", -1
	for _, pr := range c.Projects {
		base := filepath.Clean(pr.Path)
		if pathHasPrefix(abs, base) && len(base) > bestLen {
			best, bestLen = pr.ID, len(base)
		}
	}
	return best
}

// pathHasPrefix indica se p é base ou está dentro de base (case-insensitive).
func pathHasPrefix(p, base string) bool {
	if strings.EqualFold(p, base) {
		return true
	}
	withSep := base
	if !strings.HasSuffix(withSep, string(filepath.Separator)) {
		withSep += string(filepath.Separator)
	}
	return strings.HasPrefix(strings.ToLower(p), strings.ToLower(withSep))
}

// --- Aliases (Fase 1) ---

func (c *Config) Alias(name string) (string, bool) {
	v, ok := c.Aliases[name]
	return v, ok
}

func (c *Config) SetAlias(name, command string) {
	if c.Aliases == nil {
		c.Aliases = map[string]string{}
	}
	c.Aliases[name] = command
}

func (c *Config) RemoveAlias(name string) bool {
	if _, ok := c.Aliases[name]; !ok {
		return false
	}
	delete(c.Aliases, name)
	return true
}

// --- Layouts (Fase 3) ---

func (c *Config) FindLayout(id string) (*Layout, bool) {
	for i := range c.Layouts {
		if c.Layouts[i].ID == id {
			return &c.Layouts[i], true
		}
	}
	return nil, false
}

// AddLayout cria um layout a partir de um preset ("", "claude", "dev", "triple").
func (c *Config) AddLayout(name, preset string) (*Layout, error) {
	if name == "" {
		return nil, i18n.Errorf("layout name is required")
	}
	id := uniqueID(slug(name), func(id string) bool {
		_, ok := c.FindLayout(id)
		return ok
	})
	l := Layout{ID: id, Name: name, Tabs: presetTabs(preset)}
	c.Layouts = append(c.Layouts, l)
	return &c.Layouts[len(c.Layouts)-1], nil
}

func presetTabs(preset string) []Tab {
	switch preset {
	case "dev":
		return []Tab{{Panes: []Pane{
			{Shell: ShellClaude},
			{Shell: ShellPwsh, Split: "V"},
		}}}
	case "triple":
		return []Tab{{Panes: []Pane{
			{Shell: ShellClaude},
			{Shell: ShellPwsh, Split: "V"},
			{Shell: ShellPwsh, Split: "H"},
		}}}
	default: // "" / "claude"
		return []Tab{{Panes: []Pane{{Shell: ShellClaude}}}}
	}
}

// RemoveLayout remove o layout e o desreferencia dos projetos.
func (c *Config) RemoveLayout(id string) bool {
	idx := -1
	for i := range c.Layouts {
		if c.Layouts[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	c.Layouts = append(c.Layouts[:idx], c.Layouts[idx+1:]...)
	for i := range c.Projects {
		if c.Projects[i].LayoutID == id {
			c.Projects[i].LayoutID = ""
		}
	}
	return true
}

// ExpandAlias substitui {1}..{N} pelos args e {*} por todos juntos. Se o template
// não tiver placeholders, anexa os args ao final.
func ExpandAlias(tmpl string, args []string) string {
	out := tmpl
	used := false
	for i, a := range args {
		ph := fmt.Sprintf("{%d}", i+1)
		if strings.Contains(out, ph) {
			out = strings.ReplaceAll(out, ph, a)
			used = true
		}
	}
	if strings.Contains(out, "{*}") {
		out = strings.ReplaceAll(out, "{*}", strings.Join(args, " "))
		used = true
	}
	if !used && len(args) > 0 {
		out = strings.TrimSpace(out + " " + strings.Join(args, " "))
	}
	return out
}

// --- Pacotes ---

func (c *Config) FindPackage(id string) (*Package, bool) {
	for i := range c.Packages {
		if c.Packages[i].ID == id {
			return &c.Packages[i], true
		}
	}
	return nil, false
}

// AddPackage cria um pacote com id único.
func (c *Config) AddPackage(name string, projectIDs []string, mode Mode, shell Shell) (*Package, error) {
	if name == "" {
		return nil, i18n.Errorf("package name is required")
	}
	id := uniqueID(slug(name), func(id string) bool {
		_, ok := c.FindPackage(id)
		return ok
	})
	pk := Package{ID: id, Name: name, ProjectIDs: projectIDs, Mode: mode, Shell: shell}
	c.Packages = append(c.Packages, pk)
	return &c.Packages[len(c.Packages)-1], nil
}

func (c *Config) RemovePackage(id string) bool {
	for i := range c.Packages {
		if c.Packages[i].ID == id {
			c.Packages = append(c.Packages[:i], c.Packages[i+1:]...)
			return true
		}
	}
	return false
}

// --- Helpers ---

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "item"
	}
	return s
}

func uniqueID(base string, exists func(string) bool) string {
	if !exists(base) {
		return base
	}
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !exists(candidate) {
			return candidate
		}
	}
}

func removeString(s []string, v string) []string {
	out := s[:0]
	for _, x := range s {
		if x != v {
			out = append(out, x)
		}
	}
	return out
}
