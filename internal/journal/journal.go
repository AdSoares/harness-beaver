// Package journal mantém o diário de bordo: um log diário append-only (JSONL)
// dos comandos executados pelo bvr e a localização das análises de aprendizado.
package journal

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"harnessbeaver/internal/config"
)

// DateLayout é o formato de data usado nos nomes de arquivo (AAAA-MM-DD).
const DateLayout = "2006-01-02"

// MaxCapture limita o tamanho (bytes) de cada stream gravado por entrada.
const MaxCapture = 10 * 1024

// Entry é um comando executado e sua saída.
type Entry struct {
	Time       time.Time `json:"time"`
	Shell      string    `json:"shell"`
	Cwd        string    `json:"cwd"`
	ProjectID  string    `json:"projectId,omitempty"`
	Command    string    `json:"command"`
	ExitCode   int       `json:"exitCode"`
	DurationMs int64     `json:"durationMs"`
	Stdout     string    `json:"stdout"`
	Stderr     string    `json:"stderr"`
}

// Today devolve a data de hoje no formato do journal.
func Today() string { return time.Now().Format(DateLayout) }

// LogPath devolve o caminho do log de um dia.
func LogPath(date string) (string, error) {
	dir, err := config.LogsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, date+".jsonl"), nil
}

// LearningPath devolve o caminho da análise de um dia.
func LearningPath(date string) (string, error) {
	dir, err := config.LearningsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, date+".md"), nil
}

// Append grava uma entrada no log do dia da própria entrada (cria a pasta).
func Append(e Entry) error {
	if e.Time.IsZero() {
		e.Time = time.Now()
	}
	e.Stdout = truncate(e.Stdout)
	e.Stderr = truncate(e.Stderr)

	date := e.Time.Format(DateLayout)
	p, err := LogPath(date)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

// ReadDay lê todas as entradas de um dia (vazio se não houver log).
func ReadDay(date string) ([]Entry, error) {
	p, err := LogPath(date)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Entry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue // ignora linhas corrompidas
		}
		out = append(out, e)
	}
	return out, sc.Err()
}

// ReadRange lê as entradas de todos os dias com log no intervalo [from,to]
// (datas inclusivas, formato AAAA-MM-DD), em ordem cronológica.
func ReadRange(from, to string) ([]Entry, error) {
	days, err := Days()
	if err != nil {
		return nil, err
	}
	var out []Entry
	for _, d := range days {
		if from != "" && d < from {
			continue
		}
		if to != "" && d > to {
			continue
		}
		es, err := ReadDay(d)
		if err != nil {
			return nil, err
		}
		out = append(out, es...)
	}
	return out, nil
}

// FilterOpts define critérios de filtragem de entradas.
type FilterOpts struct {
	ProjectID  string
	Grep       string
	FailedOnly bool
}

// Filter aplica os critérios e devolve as entradas correspondentes (na mesma ordem).
func Filter(entries []Entry, o FilterOpts) []Entry {
	g := strings.ToLower(strings.TrimSpace(o.Grep))
	var out []Entry
	for _, e := range entries {
		if o.ProjectID != "" && e.ProjectID != o.ProjectID {
			continue
		}
		if o.FailedOnly && e.ExitCode == 0 {
			continue
		}
		if g != "" && !strings.Contains(strings.ToLower(e.Command), g) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// HasLog informa se existe log para a data.
func HasLog(date string) bool {
	p, err := LogPath(date)
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

// ReadLearning lê o markdown de análise já salvo para a data.
func ReadLearning(date string) (string, error) {
	p, err := LearningPath(date)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// HasLearning informa se já existe análise salva para a data.
func HasLearning(date string) bool {
	p, err := LearningPath(date)
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

// Days devolve as datas (ordenadas) que possuem log.
func Days() ([]string, error) {
	dir, err := config.LogsDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var days []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		days = append(days, strings.TrimSuffix(name, ".jsonl"))
	}
	sort.Strings(days)
	return days, nil
}

// PendingReviewDate devolve a data mais recente anterior a hoje que tem log mas
// ainda não tem análise. Devolve "" se não houver pendência.
func PendingReviewDate() (string, error) {
	days, err := Days()
	if err != nil {
		return "", err
	}
	today := Today()
	for i := len(days) - 1; i >= 0; i-- {
		d := days[i]
		if d >= today {
			continue
		}
		if !HasLearning(d) {
			return d, nil
		}
	}
	return "", nil
}

func truncate(s string) string {
	if len(s) <= MaxCapture {
		return s
	}
	return s[:MaxCapture] + "\n…[truncado]"
}
