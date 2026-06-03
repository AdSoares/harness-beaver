// Package tui implementa a interface interativa rica (bubbletea) do
// HarnessBeaver: menu navegável, multi-seleção e CRUD de projetos/pacotes.
package tui

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/gitstatus"
	"harnessbeaver/internal/insights"
	"harnessbeaver/internal/journal"
	"harnessbeaver/internal/launcher"
	"harnessbeaver/internal/runner"
	"harnessbeaver/internal/scanner"
)

type screen int

const (
	scMenu screen = iota
	scOpenPick
	scOpenConfirm
	scOpenMode
	scOpenShell
	scPkgOpenPick
	scProjAddName
	scProjAddPath
	scProjRemove
	scScan
	scPkgAddName
	scPkgAddPick
	scPkgAddMode
	scPkgAddShell
	scPkgRemove
	scExecInput
	scExecResult
	scReviewOffer
	scReviewPick
	scReviewChoice
	scAnalyzing
	scConfigList
	scConfigChoice
	scConfigEdit
	scHistory
	scAliases
	scLayoutPick
	scLayoutProj
)

type model struct {
	cfg    *config.Config
	screen screen
	menu   picker
	pick   picker
	input  textinput.Model
	status string
	errMsg string

	// scratch do fluxo em andamento
	openProjects []config.Project
	openMode     config.Mode
	pkgName      string
	pkgProjects  []string
	pkgMode      config.Mode
	newProjName  string
	execTitle      string
	execOutput     string
	reviewDate     string
	configKey      string
	layoutID       string
	historyEntries []journal.Entry

	viewW int
	viewH int

	spinner  spinner.Model
	gitCache map[string]gitstatus.Status // path -> status git

	quitting bool
}

// gitStatusMsg carrega os resultados de uma consulta git assíncrona.
type gitStatusMsg struct {
	results map[string]gitstatus.Status
}

// gitLoadCmd consulta o git de vários caminhos fora do event loop.
func gitLoadCmd(paths []string) tea.Cmd {
	return func() tea.Msg {
		return gitStatusMsg{results: gitstatus.QueryMany(paths)}
	}
}

// analyzeDoneMsg sinaliza o fim de uma análise assíncrona.
type analyzeDoneMsg struct {
	date string
	md   string
	path string
	err  error
}

// analyzeCmd roda a análise fora do event loop (não congela a TUI).
func analyzeCmd(cfg *config.Config, date string) tea.Cmd {
	return func() tea.Msg {
		md, path, err := insights.Analyze(cfg, date, false)
		return analyzeDoneMsg{date: date, md: md, path: path, err: err}
	}
}

// Run inicia a TUI com a configuração informada.
func Run(cfg *config.Config) error {
	runner.ResolveProjectID = cfg.ProjectIDForPath // tag de projeto no diário
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = cursorStyle
	m := model{
		cfg:      cfg,
		menu:     newPicker("HarnessBeaver — menu", mainMenuItems(), false),
		spinner:  sp,
		gitCache: map[string]gitstatus.Status{},
	}
	// Auto-oferta: se há dia anterior com log e sem análise, e ainda não
	// oferecemos hoje, abre direto na tela de oferta.
	if cfg.Settings.LastReviewOffer != journal.Today() {
		if d, _ := journal.PendingReviewDate(); d != "" {
			m.reviewDate = d
			m.pick = newPicker(
				fmt.Sprintf("Há log de %s ainda não analisado. Analisar aprendizados agora?", d),
				[]pickItem{{id: "yes", label: "Sim, analisar"}, {id: "no", label: "Agora não"}},
				false,
			)
			m.screen = scReviewOffer
		}
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd { return nil }

func mainMenuItems() []pickItem {
	return []pickItem{
		{id: "open", label: "Abrir projetos"},
		{id: "openpkg", label: "Abrir pacote"},
		{id: "openlayout", label: "Abrir com layout"},
		{id: "addproj", label: "Adicionar projeto"},
		{id: "rmproj", label: "Remover projeto"},
		{id: "scan", label: "Escanear & importar"},
		{id: "addpkg", label: "Criar pacote"},
		{id: "rmpkg", label: "Remover pacote"},
		{id: "exec", label: "Executar comando"},
		{id: "history", label: "Histórico"},
		{id: "aliases", label: "Atalhos"},
		{id: "review", label: "Analisar aprendizados"},
		{id: "config", label: "Configurações"},
		{id: "quit", label: "Sair"},
	}
}

func (m model) projectItems(multi bool) []pickItem {
	items := make([]pickItem, 0, len(m.cfg.Projects))
	for _, p := range m.cfg.Projects {
		items = append(items, pickItem{id: p.ID, label: p.Name, desc: gitDesc(m.gitCache, p)})
	}
	return items
}

// gitDesc monta a descrição do item de projeto, prefixando o status git quando
// já estiver em cache.
func gitDesc(cache map[string]gitstatus.Status, p config.Project) string {
	if st, ok := cache[p.Path]; ok && st.IsRepo {
		return "[" + st.Label() + "]  " + p.Path
	}
	return p.Path
}

// showsProjects indica se a tela atual lista projetos (recebe status git).
func (m model) showsProjects() bool {
	switch m.screen {
	case scOpenPick, scProjRemove, scPkgAddPick, scLayoutProj:
		return true
	}
	return false
}

// refreshGitDescs atualiza as descrições do picker atual com o status em cache,
// preservando seleção e cursor.
func (m *model) refreshGitDescs() {
	for i := range m.pick.items {
		if p, ok := m.cfg.FindProject(m.pick.items[i].id); ok {
			m.pick.items[i].desc = gitDesc(m.gitCache, *p)
		}
	}
}

// gitLoadForPick aplica o cache atual ao picker e dispara a consulta dos
// projetos ainda não consultados.
func (m *model) gitLoadForPick() tea.Cmd {
	m.refreshGitDescs()
	var paths []string
	for _, it := range m.pick.items {
		if p, ok := m.cfg.FindProject(it.id); ok {
			if _, cached := m.gitCache[p.Path]; !cached {
				paths = append(paths, p.Path)
			}
		}
	}
	if len(paths) == 0 {
		return nil
	}
	return gitLoadCmd(paths)
}

func (m model) packageItems() []pickItem {
	items := make([]pickItem, 0, len(m.cfg.Packages))
	for _, pk := range m.cfg.Packages {
		desc := fmt.Sprintf("%d projeto(s)", len(pk.ProjectIDs))
		items = append(items, pickItem{id: pk.ID, label: pk.Name, desc: desc})
	}
	return items
}

func modeItems(includeDefault bool) []pickItem {
	var items []pickItem
	if includeDefault {
		items = append(items, pickItem{id: "", label: "(usar default global)"})
	}
	items = append(items,
		pickItem{id: string(config.ModeTabs), label: "Abas (uma janela, uma aba por projeto)"},
		pickItem{id: string(config.ModeWindows), label: "Janelas separadas"},
	)
	return items
}

func shellItems(includeDefault bool) []pickItem {
	var items []pickItem
	if includeDefault {
		items = append(items, pickItem{id: "", label: "(usar default)"})
	}
	items = append(items,
		pickItem{id: string(config.ShellClaude), label: "Claude Code"},
		pickItem{id: string(config.ShellPwsh), label: "PowerShell"},
		pickItem{id: string(config.ShellCmd), label: "cmd"},
	)
	return items
}

func newInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 400
	ti.Width = 60
	return ti
}

func (m *model) toMenu() {
	m.screen = scMenu
	m.errMsg = ""
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewW = msg.Width
		m.viewH = msg.Height
		return m, nil
	case spinner.TickMsg:
		if m.screen != scAnalyzing {
			return m, nil // para o loop do spinner quando não está analisando
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case analyzeDoneMsg:
		if msg.err != nil {
			m.execTitle = "Falha na análise — " + msg.date
			m.execOutput = errStyle.Render(msg.err.Error())
		} else {
			m.execTitle = "Aprendizados de " + msg.date
			m.execOutput = statusStyle.Render("Salvo em: "+msg.path) + "\n\n" + msg.md
		}
		m.screen = scExecResult
		return m, nil
	case gitStatusMsg:
		if m.gitCache == nil {
			m.gitCache = map[string]gitstatus.Status{}
		}
		for k, v := range msg.results {
			m.gitCache[k] = v
		}
		if m.showsProjects() {
			m.refreshGitDescs()
		}
		return m, nil
	case tea.KeyMsg:
		// ctrl+c sempre encerra.
		if msg.Type == tea.KeyCtrlC {
			m.quitting = true
			return m, tea.Quit
		}
		if m.screen == scAnalyzing {
			return m, nil // ignora teclas enquanto analisa
		}
		if m.isInputScreen() {
			return m.updateInput(msg)
		}
		return m.updateNav(msg)
	}
	return m, nil
}

// pageSizeFor calcula quantos itens cabem na tela para o picker dado,
// considerando que itens com descrição ocupam 2 linhas.
func (m model) pageSizeFor(p *picker) int {
	avail := m.viewH - 8 // título(2) + linha + status/err + ajuda + padding(2)
	if m.viewH <= 0 {
		avail = 16 // fallback enquanto não recebemos o tamanho do terminal
	}
	lines := 1
	if p.hasDesc() {
		lines = 2
	}
	n := avail / lines
	if n < 3 {
		n = 3
	}
	return n
}

func (m model) isInputScreen() bool {
	return m.screen == scProjAddName || m.screen == scProjAddPath ||
		m.screen == scPkgAddName || m.screen == scExecInput ||
		m.screen == scConfigEdit
}

// isListScreen indica se a tela atual exibe um picker (suporta filtro/navegação).
func (m model) isListScreen() bool {
	switch m.screen {
	case scExecResult, scOpenConfirm, scAnalyzing:
		return false
	}
	return true
}

// updateNav trata telas de navegação/seleção (não-input).
func (m model) updateNav(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	active := m.activePicker()

	// Modo de filtro: as teclas alimentam o texto de busca.
	if m.isListScreen() && active.filtering {
		switch msg.String() {
		case "esc":
			active.setFilter("")
			active.filtering = false
			return m, nil
		case "enter":
			active.filtering = false // mantém o filtro, sai do modo digitação
			return m, nil
		case "backspace":
			if len(active.filter) > 0 {
				active.setFilter(active.filter[:len(active.filter)-1])
			}
			return m, nil
		case "up":
			active.up()
			return m, nil
		case "down":
			active.down()
			return m, nil
		default:
			if len(msg.Runes) > 0 {
				active.setFilter(active.filter + string(msg.Runes))
			}
			return m, nil
		}
	}

	switch msg.String() {
	case "/":
		if m.isListScreen() {
			active.filtering = true
		}
		return m, nil
	case "up", "k":
		active.up()
		return m, nil
	case "down", "j":
		active.down()
		return m, nil
	case "pgup", "left", "h":
		active.jump(-m.pageSizeFor(active))
		return m, nil
	case "pgdown", "right", "l":
		active.jump(m.pageSizeFor(active))
		return m, nil
	case "home", "g":
		active.toStart()
		return m, nil
	case "end", "G":
		active.toEnd()
		return m, nil
	case " ":
		active.toggle()
		return m, nil
	case "esc":
		if m.screen == scMenu {
			m.quitting = true
			return m, tea.Quit
		}
		if m.screen == scOpenConfirm {
			// Volta à seleção preservando as marcações.
			m.screen = scOpenPick
			m.errMsg = ""
			return m, nil
		}
		m.toMenu()
		return m, nil
	case "q":
		if m.screen == scMenu {
			m.quitting = true
			return m, tea.Quit
		}
	case "enter":
		return m.onEnter()
	}
	return m, nil
}

// activePicker devolve o ponteiro do picker em uso na tela atual.
func (m *model) activePicker() *picker {
	if m.screen == scMenu {
		return &m.menu
	}
	return &m.pick
}

func (m model) onEnter() (tea.Model, tea.Cmd) {
	m.status = ""
	m.errMsg = ""
	switch m.screen {
	case scMenu:
		return m.dispatchMenu()
	case scOpenPick:
		picked := m.pick.checked()
		if len(picked) == 0 {
			m.errMsg = "Selecione ao menos um projeto (espaço)."
			return m, nil
		}
		m.openProjects = nil
		for _, it := range picked {
			if p, ok := m.cfg.FindProject(it.id); ok {
				m.openProjects = append(m.openProjects, *p)
			}
		}
		// Mantém m.pick (a seleção) para permitir voltar com esc.
		m.screen = scOpenConfirm
		return m, nil
	case scOpenConfirm:
		m.pick = newPicker("Modo de abertura", modeItems(false), false)
		m.screen = scOpenMode
		return m, nil
	case scOpenMode:
		m.openMode = config.Mode(m.pick.selectedID())
		m.pick = newPicker("Shell em cada aba", shellItems(false), false)
		m.screen = scOpenShell
		return m, nil
	case scOpenShell:
		shell := config.Shell(m.pick.selectedID())
		items := make([]launcher.Item, 0, len(m.openProjects))
		for _, p := range m.openProjects {
			items = append(items, launcher.Item{
				Project: p,
				Shell:   launcher.ResolveShell(shell, p, m.cfg.Settings.DefaultShell),
			})
		}
		m.launch(m.openMode, items)
		m.toMenu()
		return m, nil
	case scPkgOpenPick:
		pk, ok := m.cfg.FindPackage(m.pick.selectedID())
		if !ok {
			m.errMsg = "Pacote não encontrado."
			return m, nil
		}
		mode := m.cfg.Settings.DefaultMode
		if pk.Mode != "" {
			mode = pk.Mode
		}
		projects := m.cfg.ProjectsByIDs(pk.ProjectIDs)
		items := make([]launcher.Item, 0, len(projects))
		for _, p := range projects {
			items = append(items, launcher.Item{
				Project: p,
				Shell:   launcher.ResolveShell(pk.Shell, p, m.cfg.Settings.DefaultShell),
			})
		}
		m.launch(mode, items)
		m.toMenu()
		return m, nil
	case scProjRemove:
		id := m.pick.selectedID()
		if m.cfg.RemoveProject(id) {
			m.save("Projeto removido: " + id)
		}
		m.toMenu()
		return m, nil
	case scScan:
		picked := m.pick.checked()
		n := 0
		for _, it := range picked {
			if _, err := m.cfg.AddProject(it.label, it.id, ""); err == nil {
				n++
			}
		}
		m.save(fmt.Sprintf("%d projeto(s) importado(s).", n))
		m.toMenu()
		return m, nil
	case scPkgAddPick:
		picked := m.pick.checked()
		if len(picked) == 0 {
			m.errMsg = "Selecione ao menos um projeto."
			return m, nil
		}
		m.pkgProjects = nil
		for _, it := range picked {
			m.pkgProjects = append(m.pkgProjects, it.id)
		}
		m.pick = newPicker("Modo do pacote", modeItems(true), false)
		m.screen = scPkgAddMode
		return m, nil
	case scPkgAddMode:
		m.pkgMode = config.Mode(m.pick.selectedID())
		m.pick = newPicker("Shell do pacote", shellItems(true), false)
		m.screen = scPkgAddShell
		return m, nil
	case scPkgAddShell:
		shell := config.Shell(m.pick.selectedID())
		if _, err := m.cfg.AddPackage(m.pkgName, m.pkgProjects, m.pkgMode, shell); err != nil {
			m.errMsg = err.Error()
			m.toMenu()
			return m, nil
		}
		m.save("Pacote criado: " + m.pkgName)
		m.toMenu()
		return m, nil
	case scPkgRemove:
		id := m.pick.selectedID()
		if m.cfg.RemovePackage(id) {
			m.save("Pacote removido: " + id)
		}
		m.toMenu()
		return m, nil
	case scReviewOffer:
		if m.pick.selectedID() == "yes" {
			cmd := m.beginAnalyze(m.reviewDate)
			return m, cmd
		}
		m.markOffered()
		m.toMenu()
		return m, nil
	case scReviewPick:
		date := m.pick.selectedID()
		if journal.HasLearning(date) {
			m.reviewDate = date
			label := "Dia " + date + " já tem análise salva. O que deseja?"
			if date == journal.Today() {
				label = "Hoje (" + date + ") já tem uma análise. O que deseja?"
			}
			m.pick = newPicker(label, []pickItem{
				{id: "view", label: "Ver análise existente (cache)"},
				{id: "rerun", label: "Reanalisar agora (rodar o Claude de novo)"},
			}, false)
			m.screen = scReviewChoice
			return m, nil
		}
		cmd := m.beginAnalyze(date)
		return m, cmd
	case scReviewChoice:
		if m.pick.selectedID() == "view" {
			m.showExistingLearning(m.reviewDate)
			return m, nil
		}
		cmd := m.beginAnalyze(m.reviewDate)
		return m, cmd
	case scExecResult:
		m.toMenu()
		return m, nil
	case scConfigList:
		key := m.pick.selectedID()
		def, ok := config.FindSetting(key)
		if !ok {
			m.toMenu()
			return m, nil
		}
		m.configKey = key
		if len(def.Enum) > 0 {
			var items []pickItem
			for _, opt := range def.Enum {
				items = append(items, pickItem{id: opt, label: opt})
			}
			m.pick = newPicker("Definir "+key, items, false)
			m.screen = scConfigChoice
			return m, nil
		}
		m.input = newInput(def.Help)
		m.input.SetValue(def.Get(m.cfg))
		m.screen = scConfigEdit
		return m, textinput.Blink
	case scConfigChoice:
		m.applyConfig(m.configKey, m.pick.selectedID())
		m.openConfigList()
		return m, nil
	case scHistory:
		r := m.pick.cursorReal()
		if r < 0 || r >= len(m.historyEntries) {
			m.toMenu()
			return m, nil
		}
		e := m.historyEntries[r]
		res, _ := runner.Run(config.Shell(e.Shell), e.Command, e.Cwd, false)
		m.execTitle = fmt.Sprintf("re-exec: %s   (exit %d, %dms)", e.Command, res.ExitCode, res.Duration.Milliseconds())
		m.execOutput = formatExecOutput(res)
		m.screen = scExecResult
		return m, nil
	case scAliases:
		name := m.pick.selectedID()
		tmpl, ok := m.cfg.Alias(name)
		if !ok {
			m.toMenu()
			return m, nil
		}
		command := config.ExpandAlias(tmpl, nil)
		res, _ := runner.Run(m.cfg.Settings.DefaultRunShell, command, "", false)
		m.execTitle = fmt.Sprintf("» %s   (exit %d, %dms)", command, res.ExitCode, res.Duration.Milliseconds())
		m.execOutput = formatExecOutput(res)
		m.screen = scExecResult
		return m, nil
	case scLayoutPick:
		m.layoutID = m.pick.selectedID()
		if len(m.cfg.Projects) == 0 {
			m.status = "Nenhum projeto cadastrado."
			m.toMenu()
			return m, nil
		}
		m.pick = newPicker("Abrir com layout — escolha o projeto", m.projectItems(false), false)
		m.screen = scLayoutProj
		gitCmd := m.gitLoadForPick()
		return m, gitCmd
	case scLayoutProj:
		project, okP := m.cfg.FindProject(m.pick.selectedID())
		layout, okL := m.cfg.FindLayout(m.layoutID)
		if !okP || !okL {
			m.errMsg = "Projeto ou layout não encontrado."
			m.toMenu()
			return m, nil
		}
		line, err := launcher.LaunchLayout(*project, *layout, m.cfg.Settings.TabColors, false)
		if err != nil {
			m.errMsg = "Falha ao abrir: " + err.Error()
		} else {
			m.status = fmt.Sprintf("Aberto %s com layout %s.", project.Name, layout.ID)
			_ = line
		}
		m.toMenu()
		return m, nil
	}
	return m, nil
}

// applyConfig valida/grava uma setting e define status/errMsg.
func (m *model) applyConfig(key, value string) {
	def, ok := config.FindSetting(key)
	if !ok {
		return
	}
	if err := def.Apply(m.cfg, value); err != nil {
		m.errMsg = err.Error()
		return
	}
	if err := m.cfg.Save(); err != nil {
		m.errMsg = "Erro ao salvar: " + err.Error()
		return
	}
	m.status = key + " = " + dashTUI(def.Get(m.cfg))
}

func dashTUI(s string) string {
	if s == "" {
		return "(vazio)"
	}
	return s
}

// markOffered registra que a oferta de análise já ocorreu hoje.
func (m *model) markOffered() {
	m.cfg.Settings.LastReviewOffer = journal.Today()
	_ = m.cfg.Save()
}

// beginAnalyze entra na tela de progresso e dispara a análise assíncrona.
func (m *model) beginAnalyze(date string) tea.Cmd {
	m.markOffered()
	m.reviewDate = date
	m.execTitle = "Analisando aprendizados de " + date
	m.screen = scAnalyzing
	return tea.Batch(m.spinner.Tick, analyzeCmd(m.cfg, date))
}

// showExistingLearning carrega a análise já salva (cache) e a exibe.
func (m *model) showExistingLearning(date string) {
	md, err := journal.ReadLearning(date)
	if err != nil {
		m.execTitle = "Análise de " + date
		m.execOutput = errStyle.Render("Não foi possível ler a análise existente: " + err.Error())
	} else {
		m.execTitle = "Aprendizados de " + date + " (cache)"
		m.execOutput = md
	}
	m.screen = scExecResult
}

func (m model) dispatchMenu() (tea.Model, tea.Cmd) {
	switch m.menu.selectedID() {
	case "open":
		if len(m.cfg.Projects) == 0 {
			m.status = "Nenhum projeto cadastrado. Use 'Adicionar projeto' ou 'Escanear'."
			return m, nil
		}
		m.pick = newPicker("Abrir projetos (espaço marca, enter confirma)", m.projectItems(true), true)
		m.screen = scOpenPick
		gitCmd := m.gitLoadForPick()
		return m, gitCmd
	case "openpkg":
		if len(m.cfg.Packages) == 0 {
			m.status = "Nenhum pacote cadastrado. Use 'Criar pacote'."
			return m, nil
		}
		m.pick = newPicker("Abrir pacote", m.packageItems(), false)
		m.screen = scPkgOpenPick
	case "addproj":
		m.newProjName = ""
		m.input = newInput("Nome do projeto (enter usa o nome da pasta)")
		m.screen = scProjAddName
	case "rmproj":
		if len(m.cfg.Projects) == 0 {
			m.status = "Nenhum projeto para remover."
			return m, nil
		}
		m.pick = newPicker("Remover projeto", m.projectItems(false), false)
		m.screen = scProjRemove
		gitCmd := m.gitLoadForPick()
		return m, gitCmd
	case "scan":
		return m.startScan()
	case "addpkg":
		if len(m.cfg.Projects) == 0 {
			m.status = "Cadastre projetos antes de criar um pacote."
			return m, nil
		}
		m.pkgName = ""
		m.input = newInput("Nome do pacote")
		m.screen = scPkgAddName
	case "rmpkg":
		if len(m.cfg.Packages) == 0 {
			m.status = "Nenhum pacote para remover."
			return m, nil
		}
		m.pick = newPicker("Remover pacote", m.packageItems(), false)
		m.screen = scPkgRemove
	case "exec":
		m.input = newInput("Comando a executar (ex: git status)")
		m.screen = scExecInput
	case "review":
		days, _ := journal.Days()
		if len(days) == 0 {
			m.status = "Sem logs ainda. Use 'Executar comando' para começar a registrar."
			return m, nil
		}
		var items []pickItem
		for i := len(days) - 1; i >= 0; i-- { // mais recentes primeiro
			d := days[i]
			desc := "sem análise"
			switch {
			case d == journal.Today():
				desc = "hoje — pode reanalisar"
			case journal.HasLearning(d):
				desc = "já analisado (cache disponível)"
			}
			items = append(items, pickItem{id: d, label: d, desc: desc})
		}
		m.pick = newPicker("Analisar aprendizados — escolha o dia", items, false)
		m.screen = scReviewPick
	case "config":
		m.openConfigList()
	case "history":
		m.openHistory()
	case "aliases":
		m.openAliases()
	case "openlayout":
		if len(m.cfg.Layouts) == 0 {
			m.status = "Nenhum layout. Crie com 'bvr layout add <nome> --preset dev'."
			return m, nil
		}
		var items []pickItem
		for _, l := range m.cfg.Layouts {
			panes := 0
			for _, t := range l.Tabs {
				panes += len(t.Panes)
			}
			items = append(items, pickItem{
				id: l.ID, label: l.Name,
				desc: fmt.Sprintf("%d aba(s), %d painel(éis)", len(l.Tabs), panes),
			})
		}
		m.pick = newPicker("Abrir com layout — escolha o layout", items, false)
		m.screen = scLayoutPick
	case "quit":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// openHistory monta o picker do histórico recente (mais novos primeiro).
func (m *model) openHistory() {
	to := journal.Today()
	from := time.Now().AddDate(0, 0, -6).Format(journal.DateLayout)
	entries, _ := journal.ReadRange(from, to)
	m.historyEntries = nil
	var items []pickItem
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		proj := e.ProjectID
		if proj == "" {
			proj = "-"
		}
		m.historyEntries = append(m.historyEntries, e)
		idx := len(m.historyEntries) - 1
		items = append(items, pickItem{
			id:    strconv.Itoa(idx),
			label: e.Command,
			desc:  fmt.Sprintf("%s · %s · %s · exit=%d", e.Time.Format("01-02 15:04"), e.Shell, proj, e.ExitCode),
		})
	}
	if len(items) == 0 {
		m.status = "Sem histórico nos últimos 7 dias."
		return
	}
	m.pick = newPicker("Histórico — enter re-executa · / filtra", items, false)
	m.screen = scHistory
}

// openAliases monta o picker dos atalhos cadastrados.
func (m *model) openAliases() {
	if len(m.cfg.Aliases) == 0 {
		m.status = "Nenhum atalho. Crie com 'bvr alias set <nome> <comando>'."
		return
	}
	names := make([]string, 0, len(m.cfg.Aliases))
	for n := range m.cfg.Aliases {
		names = append(names, n)
	}
	sort.Strings(names)
	var items []pickItem
	for _, n := range names {
		items = append(items, pickItem{id: n, label: n, desc: m.cfg.Aliases[n]})
	}
	m.pick = newPicker("Atalhos — enter executa · / filtra", items, false)
	m.screen = scAliases
}

// openConfigList monta o picker das settings (mostrando o valor atual).
func (m *model) openConfigList() {
	var items []pickItem
	for _, d := range config.SettingDefs {
		val := d.Get(m.cfg)
		if val == "" {
			val = "(vazio)"
		}
		items = append(items, pickItem{id: d.Key, label: d.Key, desc: val})
	}
	m.pick = newPicker("Configurações — escolha uma chave", items, false)
	m.screen = scConfigList
}

func (m model) startScan() (tea.Model, tea.Cmd) {
	root := m.cfg.Settings.ScanRoot
	if root == "" {
		if wd, err := os.Getwd(); err == nil {
			root = wd
		}
	}
	candidates, err := scanner.Scan(root, scanner.DefaultDepth)
	if err != nil {
		m.status = "Erro ao escanear: " + err.Error()
		return m, nil
	}
	// Também oferece os diretórios-pai que agrupam ≥2 projetos.
	candidates = append(candidates, scanner.Groups(candidates, root)...)
	var items []pickItem
	for _, c := range candidates {
		if m.projectExists(c.Path) {
			continue
		}
		items = append(items, pickItem{
			id:    c.Path,
			label: c.Name,
			desc:  c.Path + "  [" + strings.Join(c.Markers, ",") + "]",
		})
	}
	if len(items) == 0 {
		m.status = fmt.Sprintf("Nada novo para importar em %s", root)
		return m, nil
	}
	m.pick = newPicker(fmt.Sprintf("Escanear & importar — %s", root), items, true)
	m.screen = scScan
	return m, nil
}

func (m model) projectExists(path string) bool {
	for _, p := range m.cfg.Projects {
		if strings.EqualFold(p.Path, path) {
			return true
		}
	}
	return false
}

// updateInput trata telas de entrada de texto.
func (m model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.toMenu()
		return m, nil
	case tea.KeyEnter:
		value := strings.TrimSpace(m.input.Value())
		switch m.screen {
		case scProjAddName:
			m.newProjName = value
			m.input = newInput("Caminho do diretório (ex: F:\\proj\\meu-app)")
			m.screen = scProjAddPath
			return m, textinput.Blink
		case scProjAddPath:
			if value == "" {
				m.errMsg = "Informe um caminho."
				return m, nil
			}
			p, err := m.cfg.AddProject(m.newProjName, value, "")
			if err != nil {
				m.errMsg = err.Error()
				return m, nil
			}
			m.save("Projeto adicionado: " + p.ID)
			m.toMenu()
			return m, nil
		case scPkgAddName:
			if value == "" {
				m.errMsg = "Informe um nome."
				return m, nil
			}
			m.pkgName = value
			m.pick = newPicker("Projetos do pacote (espaço marca)", m.projectItems(true), true)
			m.screen = scPkgAddPick
			gitCmd := m.gitLoadForPick()
			return m, gitCmd
		case scExecInput:
			if value == "" {
				m.errMsg = "Informe um comando."
				return m, nil
			}
			res, _ := runner.Run(m.cfg.Settings.DefaultRunShell, value, "", false)
			m.execTitle = fmt.Sprintf("$ %s   (exit %d, %dms)", value, res.ExitCode, res.Duration.Milliseconds())
			m.execOutput = formatExecOutput(res)
			m.screen = scExecResult
			return m, nil
		case scConfigEdit:
			m.applyConfig(m.configKey, value) // valor vazio é permitido (desliga)
			m.openConfigList()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// launch dispara o launcher e registra o resultado em status/errMsg.
func (m *model) launch(mode config.Mode, items []launcher.Item) {
	lines, err := launcher.Launch(mode, items, m.cfg.Settings.TabColors, false)
	if err != nil {
		m.errMsg = "Falha ao abrir: " + err.Error()
		return
	}
	m.status = fmt.Sprintf("Aberto %d projeto(s) em modo %s. (%d comando wt)", len(items), mode, len(lines))
}

// save persiste a config e define status (ou errMsg em caso de falha).
func (m *model) save(okMsg string) {
	if err := m.cfg.Save(); err != nil {
		m.errMsg = "Erro ao salvar: " + err.Error()
		return
	}
	m.status = okMsg
}

// formatExecOutput monta a saída capturada de um comando para exibição.
func formatExecOutput(res runner.Result) string {
	var b strings.Builder
	if out := strings.TrimRight(res.Stdout, "\n"); out != "" {
		b.WriteString(out + "\n")
	}
	if errOut := strings.TrimRight(res.Stderr, "\n"); errOut != "" {
		b.WriteString(errStyle.Render(errOut) + "\n")
	}
	if b.Len() == 0 {
		return descStyle.Render("(sem saída)")
	}
	return lastLines(b.String(), 30)
}

// lastLines devolve as últimas n linhas de s (evita estourar a tela).
func lastLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) <= n {
		return strings.Join(lines, "\n")
	}
	trimmed := lines[len(lines)-n:]
	return descStyle.Render(fmt.Sprintf("…(%d linhas anteriores omitidas)", len(lines)-n)) +
		"\n" + strings.Join(trimmed, "\n")
}

// openConfirmView mostra os projetos selecionados antes de confirmar a abertura.
func (m model) openConfirmView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Confirmar projetos a abrir") + "\n\n")

	itemsFit := (m.viewH - 9) / 2
	if m.viewH <= 0 {
		itemsFit = 7
	}
	if itemsFit < 3 {
		itemsFit = 3
	}
	shown := len(m.openProjects)
	if shown > itemsFit {
		shown = itemsFit
	}
	for i := 0; i < shown; i++ {
		p := m.openProjects[i]
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, p.Name))
		b.WriteString("      " + descStyle.Render(p.Path) + "\n")
	}
	if len(m.openProjects) > shown {
		b.WriteString(descStyle.Render(fmt.Sprintf("  … e mais %d", len(m.openProjects)-shown)) + "\n")
	}
	b.WriteString("\n" + descStyle.Render(fmt.Sprintf("%d projeto(s) selecionado(s)", len(m.openProjects))) + "\n")
	return b.String()
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var b strings.Builder

	if m.screen == scMenu {
		b.WriteString(m.menu.view(m.pageSizeFor(&m.menu)))
	} else if m.isInputScreen() {
		b.WriteString(titleStyle.Render(m.inputTitle()) + "\n\n")
		b.WriteString(m.input.View() + "\n")
	} else if m.screen == scExecResult {
		b.WriteString(titleStyle.Render(m.execTitle) + "\n\n")
		b.WriteString(m.execOutput + "\n")
	} else if m.screen == scOpenConfirm {
		b.WriteString(m.openConfirmView())
	} else if m.screen == scAnalyzing {
		b.WriteString(titleStyle.Render(m.execTitle) + "\n\n")
		b.WriteString(m.spinner.View() + " chamando o Claude, aguarde…\n")
	} else {
		b.WriteString(m.pick.view(m.pageSizeFor(&m.pick)))
	}

	b.WriteString("\n")
	if m.status != "" {
		b.WriteString(statusStyle.Render(m.status) + "\n")
	}
	if m.errMsg != "" {
		b.WriteString(errStyle.Render(m.errMsg) + "\n")
	}
	b.WriteString(helpStyle.Render(m.helpLine()))

	return appStyle.Render(b.String())
}

func (m model) inputTitle() string {
	switch m.screen {
	case scProjAddName:
		return "Adicionar projeto — nome"
	case scProjAddPath:
		return "Adicionar projeto — caminho"
	case scPkgAddName:
		return "Criar pacote — nome"
	case scExecInput:
		return "Executar comando"
	case scConfigEdit:
		return "Definir " + m.configKey
	}
	return ""
}

func (m model) helpLine() string {
	if m.isInputScreen() {
		return "enter: confirmar · esc: cancelar · ctrl+c: sair"
	}
	if m.screen == scAnalyzing {
		return "analisando… (ctrl+c sai)"
	}
	if m.screen == scExecResult {
		return "enter/esc: voltar ao menu"
	}
	if m.screen == scOpenConfirm {
		return "enter: abrir · esc: voltar à seleção"
	}
	if m.isListScreen() && m.activePicker().filtering {
		return "filtrando: digite · ↑/↓ navega · enter aplica · esc limpa"
	}
	if m.screen == scMenu {
		return "↑/↓ navegar · / filtrar · enter selecionar · q sair"
	}
	if m.pick.multi {
		return "↑/↓ nav · ←/→ págs · / filtrar · espaço marcar · enter ok · esc voltar"
	}
	return "↑/↓ nav · ←/→ págs · / filtrar · enter selecionar · esc voltar"
}
