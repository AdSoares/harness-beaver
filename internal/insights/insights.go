// Package insights monta um prompt a partir do diário de um dia e gera uma
// análise de aprendizados via Claude (CLI claude -p ou API Anthropic).
package insights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/journal"
)

const apiURL = "https://api.anthropic.com/v1/messages"

// instruction é a diretiva (PT-BR) enviada ao modelo.
func instruction(date string) string {
	return fmt.Sprintf(`Você é um engenheiro sênior revisando o diário de bordo de um desenvolvedor.
Abaixo está o registro dos comandos de terminal executados em %s, com saídas e códigos de saída.

Produza uma análise em **markdown (português do Brasil)** com as seções:
- ## Resumo do dia — o que foi feito, em 2-4 frases.
- ## Comandos notáveis — os mais relevantes e por quê.
- ## Erros e como foram resolvidos — comandos que falharam (exit != 0) e a provável causa/solução.
- ## Aprendizados — lições práticas extraídas do dia.
- ## Sugestões — melhorias, atalhos ou próximos passos.

Seja específico e conciso. Ignore ruído. Não invente comandos que não aparecem no log.`, date)
}

// formatLog renderiza as entradas do dia em texto legível.
func formatLog(entries []journal.Entry) string {
	var b strings.Builder
	for i, e := range entries {
		proj := e.ProjectID
		if proj == "" {
			proj = "-"
		}
		fmt.Fprintf(&b, "### %d. [%s] (%s) proj=%s exit=%d dur=%dms\n",
			i+1, e.Time.Format("15:04:05"), e.Shell, proj, e.ExitCode, e.DurationMs)
		fmt.Fprintf(&b, "cwd: %s\n", e.Cwd)
		fmt.Fprintf(&b, "$ %s\n", e.Command)
		if out := strings.TrimSpace(e.Stdout); out != "" {
			fmt.Fprintf(&b, "stdout:\n%s\n", out)
		}
		if errOut := strings.TrimSpace(e.Stderr); errOut != "" {
			fmt.Fprintf(&b, "stderr:\n%s\n", errOut)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// BuildPrompt monta o prompt completo (instrução + log), usado no modo dry-run
// e no caminho da API.
func BuildPrompt(date string, entries []journal.Entry) string {
	return instruction(date) + "\n\n--- LOG DO DIA ---\n\n" + formatLog(entries)
}

// Analyze gera a análise do dia. Com dryRun, devolve apenas o prompt montado e
// não chama motor nem salva. Caso contrário, chama o motor configurado e salva
// o markdown em ~/.harnessbeaver/learnings/<date>.md.
func Analyze(cfg *config.Config, date string, dryRun bool) (markdown string, savedPath string, err error) {
	entries, err := journal.ReadDay(date)
	if err != nil {
		return "", "", err
	}
	if len(entries) == 0 {
		return "", "", fmt.Errorf("sem comandos registrados em %s", date)
	}

	if dryRun {
		return BuildPrompt(date, entries), "", nil
	}

	md, err := callEngine(cfg, date, entries)
	if err != nil {
		return "", "", err
	}

	path, err := journal.LearningPath(date)
	if err != nil {
		return md, "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return md, "", err
	}
	header := fmt.Sprintf("# Aprendizados — %s\n\n_Gerado pelo bvr em %s._\n\n",
		date, time.Now().Format("2006-01-02 15:04"))
	content := []byte(header + md)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return md, "", err
	}
	// Cópia extra opcional (ex.: _content/learnings da empresa). Falha aqui não
	// invalida o salvamento principal.
	if dir := strings.TrimSpace(cfg.Settings.LearningsExtraDir); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err == nil {
			_ = os.WriteFile(filepath.Join(dir, date+".md"), content, 0o644)
		}
	}
	return md, path, nil
}

// callEngine escolhe e invoca o motor de acordo com a configuração.
func callEngine(cfg *config.Config, date string, entries []journal.Entry) (string, error) {
	engine := cfg.Settings.InsightsEngine
	if engine == "" {
		engine = "auto"
	}

	claudeAvailable := false
	if _, err := exec.LookPath("claude"); err == nil {
		claudeAvailable = true
	}

	switch engine {
	case "claude":
		return callClaude(date, entries)
	case "api":
		return callAPI(cfg, date, entries)
	default: // auto
		if claudeAvailable {
			return callClaude(date, entries)
		}
		return callAPI(cfg, date, entries)
	}
}

// callClaude usa o CLI do Claude Code em modo print, com o log via stdin.
func callClaude(date string, entries []journal.Entry) (string, error) {
	cmd := exec.Command("claude", "-p", instruction(date))
	cmd.Stdin = strings.NewReader("--- LOG DO DIA ---\n\n" + formatLog(entries))
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("falha ao chamar claude: %w (%s)", err, strings.TrimSpace(errBuf.String()))
	}
	res := strings.TrimSpace(out.String())
	if res == "" {
		return "", fmt.Errorf("claude retornou vazio")
	}
	return res, nil
}

// --- API Anthropic ---

type apiRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	Messages  []apiMessage `json:"messages"`
}

type apiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type apiResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func callAPI(cfg *config.Config, date string, entries []journal.Entry) (string, error) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return "", fmt.Errorf("motor 'api' requer ANTHROPIC_API_KEY (ou instale o claude CLI e use insightsEngine=auto)")
	}
	model := cfg.Settings.InsightsModel
	if model == "" {
		model = "claude-sonnet-4-6"
	}

	body, err := json.Marshal(apiRequest{
		Model:     model,
		MaxTokens: 2000,
		Messages:  []apiMessage{{Role: "user", Content: BuildPrompt(date, entries)}},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var ar apiResponse
	if err := json.Unmarshal(raw, &ar); err != nil {
		return "", fmt.Errorf("resposta inválida da API (%d): %s", resp.StatusCode, string(raw))
	}
	if ar.Error != nil {
		return "", fmt.Errorf("API: %s", ar.Error.Message)
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	var sb strings.Builder
	for _, c := range ar.Content {
		if c.Type == "text" {
			sb.WriteString(c.Text)
		}
	}
	res := strings.TrimSpace(sb.String())
	if res == "" {
		return "", fmt.Errorf("API retornou conteúdo vazio")
	}
	return res, nil
}
