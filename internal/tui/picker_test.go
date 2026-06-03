package tui

import (
	"strings"
	"testing"
)

func makeItems(n int) []pickItem {
	items := make([]pickItem, n)
	for i := 0; i < n; i++ {
		items[i] = pickItem{id: string(rune('a' + i)), label: "proj-" + string(rune('0'+i%10))}
	}
	return items
}

func TestWindowBounds(t *testing.T) {
	cases := []struct {
		name                         string
		cursor, total, page          int
		wantStart, wantEnd           int
	}{
		{"cabe tudo", 0, 5, 10, 0, 5},
		{"topo", 0, 28, 10, 0, 10},
		{"fim", 27, 28, 10, 18, 28},
		{"meio centraliza", 14, 28, 10, 9, 19},
		{"perto do fim clampa", 25, 28, 10, 18, 28},
	}
	for _, c := range cases {
		gotStart, gotEnd := windowBounds(c.cursor, c.total, c.page)
		if gotStart != c.wantStart || gotEnd != c.wantEnd {
			t.Errorf("%s: windowBounds(%d,%d,%d) = (%d,%d), quer (%d,%d)",
				c.name, c.cursor, c.total, c.page, gotStart, gotEnd, c.wantStart, c.wantEnd)
		}
	}
}

// O bug relatado: com o cursor no início, o primeiro item deve aparecer.
func TestViewShowsFirstItemAtTop(t *testing.T) {
	p := newPicker("teste", makeItems(28), true)
	p.cursor = 0
	out := p.view(10)
	if !strings.Contains(out, "proj-0") {
		t.Fatalf("primeiro item não apareceu com cursor no topo:\n%s", out)
	}
	if !strings.Contains(out, "▼ 18 abaixo") {
		t.Errorf("faltou indicador de itens abaixo:\n%s", out)
	}
	if strings.Contains(out, "▲") {
		t.Errorf("não deveria haver indicador 'acima' no topo:\n%s", out)
	}
	if !strings.Contains(out, "1–10 de 28") {
		t.Errorf("faltou rodapé de posição:\n%s", out)
	}
}

func TestViewAtBottomShowsAboveIndicator(t *testing.T) {
	p := newPicker("teste", makeItems(28), false)
	p.cursor = 27
	out := p.view(10)
	if !strings.Contains(out, "▲ 18 acima") {
		t.Errorf("faltou indicador 'acima' no fim:\n%s", out)
	}
	if strings.Contains(out, "▼") {
		t.Errorf("não deveria haver 'abaixo' no fim:\n%s", out)
	}
}

func TestFilterReducesVisibleAndSelectsCorrectly(t *testing.T) {
	items := []pickItem{
		{id: "a", label: "repair", desc: "smb"},
		{id: "b", label: "beauty", desc: "smb"},
		{id: "c", label: "builder-map", desc: "engenharia"},
		{id: "d", label: "content-forge", desc: "midias"},
	}
	p := newPicker("t", items, false)

	// Sem filtro: todos visíveis.
	if got := len(p.visible()); got != 4 {
		t.Fatalf("sem filtro esperava 4 visíveis, veio %d", got)
	}

	// Filtro por "b" casa repair(desc smb? não) — 'b' aparece em beauty e builder-map e... 'b' em "smb"?
	// Para ser determinístico, filtra por "build".
	p.setFilter("build")
	vis := p.visible()
	if len(vis) != 1 {
		t.Fatalf("filtro 'build' esperava 1 visível, veio %d", len(vis))
	}
	if p.selectedID() != "c" {
		t.Errorf("selectedID com filtro = %q, esperava 'c'", p.selectedID())
	}

	// Filtro por categoria no desc.
	p.setFilter("smb")
	if got := len(p.visible()); got != 2 {
		t.Errorf("filtro 'smb' (desc) esperava 2, veio %d", got)
	}
}

func TestFilterNavigationStaysInSubset(t *testing.T) {
	items := []pickItem{
		{id: "a", label: "alpha"},
		{id: "b", label: "beta"},
		{id: "c", label: "gamma-beta"},
	}
	p := newPicker("t", items, false)
	p.setFilter("beta") // casa "beta" e "gamma-beta"
	if len(p.visible()) != 2 {
		t.Fatalf("esperava 2 visíveis, veio %d", len(p.visible()))
	}
	p.cursor = 0
	p.down()
	if p.selectedID() != "c" {
		t.Errorf("após down no subconjunto, selectedID = %q, esperava 'c'", p.selectedID())
	}
	p.down() // wrap dentro do subconjunto
	if p.selectedID() != "b" {
		t.Errorf("após wrap, selectedID = %q, esperava 'b'", p.selectedID())
	}
}

func TestCheckedSurvivesFilterChange(t *testing.T) {
	items := []pickItem{
		{id: "a", label: "alpha"},
		{id: "b", label: "beta"},
	}
	p := newPicker("t", items, true)
	p.cursor = 0
	p.toggle() // marca alpha
	p.setFilter("beta")
	p.cursor = 0
	p.toggle() // marca beta (único visível)
	p.setFilter("")
	if got := len(p.checked()); got != 2 {
		t.Errorf("esperava 2 marcados após mudar filtro, veio %d", got)
	}
}

func TestJumpClamps(t *testing.T) {
	p := newPicker("t", makeItems(28), false)
	p.cursor = 5
	p.jump(-10)
	if p.cursor != 0 {
		t.Errorf("jump(-10) de 5 deveria clampar em 0, deu %d", p.cursor)
	}
	p.jump(100)
	if p.cursor != 27 {
		t.Errorf("jump(100) deveria clampar em 27, deu %d", p.cursor)
	}
}
