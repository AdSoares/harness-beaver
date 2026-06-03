// Package gitstatus consulta o estado git de um diretório (branch, alterações
// pendentes e ahead/behind do upstream) via `git status --porcelain=v2`.
package gitstatus

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Status resume o estado git de um diretório.
type Status struct {
	IsRepo      bool
	Branch      string
	Detached    bool
	Dirty       bool // há alterações não commitadas (staged/unstaged/untracked)
	HasUpstream bool
	Ahead       int
	Behind      int
}

// Query roda o git no diretório e devolve seu estado. Diretórios que não são
// repositório (ou erro de git) voltam com IsRepo=false.
func Query(path string) Status {
	cmd := exec.Command("git", "-C", path, "status", "--porcelain=v2", "--branch")
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return Status{IsRepo: false}
	}
	s := parse(out.String())
	s.IsRepo = true
	return s
}

// QueryMany consulta vários caminhos em paralelo (limitado), devolvendo um mapa
// path→Status.
func QueryMany(paths []string) map[string]Status {
	out := make(map[string]Status, len(paths))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	for _, p := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			st := Query(p)
			mu.Lock()
			out[p] = st
			mu.Unlock()
		}(p)
	}
	wg.Wait()
	return out
}

// parse interpreta a saída de `git status --porcelain=v2 --branch`.
func parse(out string) Status {
	var s Status
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "# branch.head "):
			s.Branch = strings.TrimPrefix(line, "# branch.head ")
			if s.Branch == "(detached)" {
				s.Detached = true
			}
		case strings.HasPrefix(line, "# branch.upstream "):
			s.HasUpstream = true
		case strings.HasPrefix(line, "# branch.ab "):
			fields := strings.Fields(strings.TrimPrefix(line, "# branch.ab "))
			if len(fields) == 2 {
				s.Ahead = absAtoi(fields[0])  // +A
				s.Behind = absAtoi(fields[1]) // -B
			}
		case line == "" || strings.HasPrefix(line, "#"):
			// outros cabeçalhos: ignora
		default:
			// qualquer entrada de mudança (1/2/u/?) => árvore suja
			s.Dirty = true
		}
	}
	return s
}

func absAtoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimLeft(s, "+-"))
	if n < 0 {
		n = -n
	}
	return n
}

// Label devolve uma representação compacta (ex.: "main ✗ ↑2 ↓1"). Vazio se não
// for repositório.
func (s Status) Label() string {
	if !s.IsRepo {
		return ""
	}
	b := s.Branch
	if s.Detached || b == "" {
		b = "detached"
	}
	out := b
	if s.Dirty {
		out += " ✗"
	} else {
		out += " ✓"
	}
	if s.Ahead > 0 {
		out += fmt.Sprintf(" ↑%d", s.Ahead)
	}
	if s.Behind > 0 {
		out += fmt.Sprintf(" ↓%d", s.Behind)
	}
	return out
}
