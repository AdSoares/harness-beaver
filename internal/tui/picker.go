package tui

import (
	"fmt"
	"strings"
)

// pickItem é uma linha selecionável do picker.
type pickItem struct {
	id      string
	label   string
	desc    string
	checked bool
}

// picker é uma lista navegável com paginação, multi-seleção e filtro de busca.
// O cursor indexa a lista VISÍVEL (após aplicar o filtro).
type picker struct {
	title     string
	items     []pickItem
	cursor    int
	multi     bool
	filter    string
	filtering bool
}

func newPicker(title string, items []pickItem, multi bool) picker {
	return picker{title: title, items: items, multi: multi}
}

// visible devolve os índices (em p.items) que passam pelo filtro atual.
func (p picker) visible() []int {
	q := strings.ToLower(strings.TrimSpace(p.filter))
	idx := make([]int, 0, len(p.items))
	for i, it := range p.items {
		if q == "" || itemMatches(it, q) {
			idx = append(idx, i)
		}
	}
	return idx
}

// itemMatches faz match case-insensitive (substring) sobre label + desc.
func itemMatches(it pickItem, qLower string) bool {
	return strings.Contains(strings.ToLower(it.label+" "+it.desc), qLower)
}

// cursorReal devolve o índice real (em p.items) do item sob o cursor, ou -1.
func (p picker) cursorReal() int {
	vis := p.visible()
	if len(vis) == 0 || p.cursor < 0 || p.cursor >= len(vis) {
		return -1
	}
	return vis[p.cursor]
}

func (p *picker) up() {
	n := len(p.visible())
	if n == 0 {
		return
	}
	p.cursor--
	if p.cursor < 0 {
		p.cursor = n - 1
	}
}

func (p *picker) down() {
	n := len(p.visible())
	if n == 0 {
		return
	}
	p.cursor++
	if p.cursor >= n {
		p.cursor = 0
	}
}

// jump move o cursor por delta posições, sem wrap-around (usado na paginação).
func (p *picker) jump(delta int) {
	n := len(p.visible())
	if n == 0 {
		return
	}
	p.cursor += delta
	if p.cursor < 0 {
		p.cursor = 0
	}
	if p.cursor >= n {
		p.cursor = n - 1
	}
}

func (p *picker) toStart() { p.cursor = 0 }

func (p *picker) toEnd() {
	if n := len(p.visible()); n > 0 {
		p.cursor = n - 1
	}
}

// setFilter atualiza o filtro e reposiciona o cursor no topo da lista filtrada.
func (p *picker) setFilter(s string) {
	p.filter = s
	p.cursor = 0
}

func (p *picker) toggle() {
	if !p.multi {
		return
	}
	if r := p.cursorReal(); r >= 0 {
		p.items[r].checked = !p.items[r].checked
	}
}

// hasDesc informa se algum item tem descrição (cada um ocupa 2 linhas).
func (p picker) hasDesc() bool {
	for _, it := range p.items {
		if it.desc != "" {
			return true
		}
	}
	return false
}

func (p *picker) selectedID() string {
	if r := p.cursorReal(); r >= 0 {
		return p.items[r].id
	}
	return ""
}

func (p *picker) checked() []pickItem {
	var out []pickItem
	for _, it := range p.items {
		if it.checked {
			out = append(out, it)
		}
	}
	return out
}

// view renderiza uma janela de pageSize itens (da lista filtrada) em torno do
// cursor, com indicadores de itens acima/abaixo, filtro ativo e posição.
func (p picker) view(pageSize int) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(p.title) + "\n")
	if p.filtering || p.filter != "" {
		caret := ""
		if p.filtering {
			caret = "▏"
		}
		b.WriteString(descStyle.Render("filtro: ") + p.filter + caret + "\n")
	}
	b.WriteString("\n")

	vis := p.visible()
	if len(vis) == 0 {
		if p.filter != "" {
			b.WriteString(descStyle.Render("(nenhum item corresponde ao filtro)") + "\n")
		} else {
			b.WriteString(descStyle.Render("(vazio)") + "\n")
		}
		return b.String()
	}
	if pageSize < 1 {
		pageSize = len(vis)
	}
	start, end := windowBounds(p.cursor, len(vis), pageSize)

	if start > 0 {
		b.WriteString(descStyle.Render(fmt.Sprintf("  ▲ %d acima", start)) + "\n")
	}
	for vi := start; vi < end; vi++ {
		it := p.items[vis[vi]]
		cursor := "  "
		if vi == p.cursor {
			cursor = cursorStyle.Render("> ")
		}
		check := ""
		if p.multi {
			if it.checked {
				check = "[x] "
			} else {
				check = "[ ] "
			}
		}
		line := check + it.label
		if vi == p.cursor {
			line = selectedStyle.Render(line)
		}
		b.WriteString(cursor + line + "\n")
		if it.desc != "" {
			b.WriteString("      " + descStyle.Render(it.desc) + "\n")
		}
	}
	if end < len(vis) {
		b.WriteString(descStyle.Render(fmt.Sprintf("  ▼ %d abaixo", len(vis)-end)) + "\n")
	}

	footer := fmt.Sprintf("  %d–%d de %d", start+1, end, len(vis))
	if p.filter != "" {
		footer += fmt.Sprintf(" (filtrados de %d)", len(p.items))
	}
	if p.multi {
		footer += fmt.Sprintf(" · %d marcado(s)", len(p.checked()))
	}
	b.WriteString("\n" + descStyle.Render(footer) + "\n")
	return b.String()
}

// windowBounds calcula o intervalo [start,end) visível mantendo o cursor à vista
// (centralizado quando possível) e respeitando os limites da lista.
func windowBounds(cursor, total, pageSize int) (int, int) {
	if total <= pageSize {
		return 0, total
	}
	start := cursor - pageSize/2
	if start < 0 {
		start = 0
	}
	if start+pageSize > total {
		start = total - pageSize
	}
	return start, start + pageSize
}
