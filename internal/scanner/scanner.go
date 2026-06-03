// Package scanner varre um diretório raiz em busca de subdiretórios que
// pareçam projetos de código, detectados por marcadores conhecidos.
package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Candidate é um diretório candidato a ser registrado: um projeto de código
// (folha) ou um diretório-pai que agrupa vários projetos (IsGroup).
type Candidate struct {
	Name    string
	Path    string
	Markers []string
	IsGroup bool
	Count   int // nº de projetos agrupados (quando IsGroup)
}

// DefaultDepth é a profundidade máxima padrão de varredura a partir da raiz.
const DefaultDepth = 3

// fileMarkers são arquivos/pastas cuja presença indica um projeto.
var fileMarkers = []string{
	".git", "package.json", "go.mod", "pyproject.toml", "Cargo.toml",
	"src", "CLAUDE.md", "AGENTS.md", "requirements.txt",
}

// globMarkers são padrões glob (ex.: *.sln) testados na raiz do diretório.
var globMarkers = []string{"*.sln", "*.csproj"}

// Scan caminha em root até maxDepth níveis. Ao encontrar um diretório com
// marcadores, ele é registrado como candidato e o ramo é podado (não descemos
// mais para não listar subpastas do mesmo repositório).
func Scan(root string, maxDepth int) ([]Candidate, error) {
	if maxDepth <= 0 {
		maxDepth = DefaultDepth
	}
	rootClean := filepath.Clean(root)
	if _, err := os.Stat(rootClean); err != nil {
		return nil, err
	}

	var out []Candidate
	err := filepath.WalkDir(rootClean, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // ignora diretórios ilegíveis
		}
		if !d.IsDir() {
			return nil
		}
		if path == rootClean {
			return nil
		}
		if depth(rootClean, path) > maxDepth {
			return filepath.SkipDir
		}
		if markers := detectMarkers(path); len(markers) > 0 {
			out = append(out, Candidate{
				Name:    filepath.Base(path),
				Path:    path,
				Markers: markers,
			})
			return filepath.SkipDir // poda o ramo
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

// Groups deriva, a partir das folhas (projetos detectados), os diretórios-pai
// dentro de root que são ancestrais de ≥2 projetos. Útil para registrar e abrir
// uma ferramenta no diretório que agrupa vários projetos.
func Groups(leaves []Candidate, root string) []Candidate {
	root = filepath.Clean(root)
	leafSet := make(map[string]bool, len(leaves))
	for _, lf := range leaves {
		leafSet[filepath.Clean(lf.Path)] = true
	}

	counts := map[string]int{}
	for _, lf := range leaves {
		dir := filepath.Dir(filepath.Clean(lf.Path))
		for within(dir, root) {
			counts[dir]++
			if dir == root {
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	var out []Candidate
	for dir, n := range counts {
		if n < 2 || leafSet[dir] {
			continue
		}
		out = append(out, Candidate{
			Name:    filepath.Base(dir),
			Path:    dir,
			IsGroup: true,
			Count:   n,
			Markers: []string{fmt.Sprintf("grupo: %d projetos", n)},
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// within indica se p é root ou está dentro de root.
func within(p, root string) bool {
	if strings.EqualFold(p, root) {
		return true
	}
	withSep := root
	if !strings.HasSuffix(withSep, string(filepath.Separator)) {
		withSep += string(filepath.Separator)
	}
	return strings.HasPrefix(strings.ToLower(p), strings.ToLower(withSep))
}

func depth(root, path string) int {
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." {
		return 0
	}
	return len(strings.Split(rel, string(os.PathSeparator)))
}

func detectMarkers(dir string) []string {
	var found []string
	for _, m := range fileMarkers {
		if _, err := os.Stat(filepath.Join(dir, m)); err == nil {
			found = append(found, m)
		}
	}
	for _, g := range globMarkers {
		matches, _ := filepath.Glob(filepath.Join(dir, g))
		if len(matches) > 0 {
			found = append(found, strings.TrimPrefix(g, "*"))
		}
	}
	return found
}
