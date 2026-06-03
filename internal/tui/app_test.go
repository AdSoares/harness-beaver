package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"harnessbeaver/internal/config"
)

// setupTempHome isola ~/.harnessbeaver num diretório temporário, para os testes
// não tocarem na config/learnings reais do usuário.
func setupTempHome(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("USERPROFILE", tmp) // Windows: os.UserHomeDir lê daqui
	t.Setenv("HOME", tmp)        // outros SOs
}

// writeLearning cria um arquivo de análise (cache) para um dia.
func writeLearning(t *testing.T, date, content string) {
	dir, err := config.LearningsDir()
	if err != nil {
		t.Fatalf("LearningsDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, date+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func dayPicker(date string) picker {
	return newPicker("", []pickItem{{id: date, label: date}}, false)
}

func choicePicker(cursor int) picker {
	p := newPicker("", []pickItem{
		{id: "view", label: "ver"},
		{id: "rerun", label: "rerun"},
	}, false)
	p.cursor = cursor
	return p
}

// Dia já analisado: escolher na lista deve levar à tela de escolha (ver/reanalisar),
// sem disparar análise.
func TestReviewPick_AnalyzedDay_GoesToChoice(t *testing.T) {
	setupTempHome(t)
	date := "2026-05-30"
	writeLearning(t, date, "# cache")

	m := model{cfg: config.Default(), screen: scReviewPick, pick: dayPicker(date)}
	res, cmd := m.onEnter()
	rm := res.(model)

	if rm.screen != scReviewChoice {
		t.Fatalf("esperava scReviewChoice, veio screen=%d", rm.screen)
	}
	if cmd != nil {
		t.Error("não deveria disparar cmd ao apenas abrir a escolha")
	}
	if rm.reviewDate != date {
		t.Errorf("reviewDate = %q, esperava %q", rm.reviewDate, date)
	}
	if len(rm.pick.items) != 2 || rm.pick.items[0].id != "view" || rm.pick.items[1].id != "rerun" {
		t.Errorf("opções inesperadas: %+v", rm.pick.items)
	}
}

// Dia sem análise: escolher na lista deve iniciar a análise (tela de progresso + cmd).
func TestReviewPick_UnanalyzedDay_StartsAnalysis(t *testing.T) {
	setupTempHome(t)
	date := "2026-05-29" // sem arquivo de learning

	m := model{cfg: config.Default(), screen: scReviewPick, pick: dayPicker(date)}
	res, cmd := m.onEnter()
	rm := res.(model)

	if rm.screen != scAnalyzing {
		t.Fatalf("esperava scAnalyzing, veio screen=%d", rm.screen)
	}
	if cmd == nil {
		t.Error("esperava cmd de análise não-nil")
	}
}

// Na tela de escolha, "ver" deve mostrar o cache (scExecResult) sem disparar cmd.
func TestReviewChoice_View_ShowsCache(t *testing.T) {
	setupTempHome(t)
	date := "2026-05-28"
	writeLearning(t, date, "# Aprendizados — cache de teste\nlinha")

	m := model{
		cfg:        config.Default(),
		screen:     scReviewChoice,
		reviewDate: date,
		pick:       choicePicker(0), // "view"
	}
	res, cmd := m.onEnter()
	rm := res.(model)

	if rm.screen != scExecResult {
		t.Fatalf("esperava scExecResult, veio screen=%d", rm.screen)
	}
	if cmd != nil {
		t.Error("ver cache não deve disparar cmd (sem chamar a IA)")
	}
	if !strings.Contains(rm.execOutput, "cache de teste") {
		t.Errorf("execOutput não contém o conteúdo do cache:\n%s", rm.execOutput)
	}
}

// Na tela de escolha, "rerun" deve iniciar a análise (scAnalyzing + cmd).
func TestReviewChoice_Rerun_StartsAnalysis(t *testing.T) {
	setupTempHome(t)
	date := "2026-05-27"
	writeLearning(t, date, "# análise antiga")

	m := model{
		cfg:        config.Default(),
		screen:     scReviewChoice,
		reviewDate: date,
		pick:       choicePicker(1), // "rerun"
	}
	res, cmd := m.onEnter()
	rm := res.(model)

	if rm.screen != scAnalyzing {
		t.Fatalf("esperava scAnalyzing, veio screen=%d", rm.screen)
	}
	if cmd == nil {
		t.Error("rerun deve disparar cmd de análise")
	}
}
