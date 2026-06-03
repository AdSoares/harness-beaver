package config

import (
	"path/filepath"
	"testing"
)

func TestProjectIDForPath(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "proj", "app")
	proj := filepath.Join(root, "proj")
	cfg := &Config{Projects: []Project{
		{ID: "app", Path: app},
		{ID: "proj", Path: proj},
	}}

	cases := []struct {
		query string
		want  string
	}{
		{filepath.Join(app, "src", "main.go"), "app"}, // prefixo mais longo vence
		{app, "app"},                                  // o próprio diretório
		{filepath.Join(proj, "outro"), "proj"},        // cai no pai
		{filepath.Join(root, "fora"), ""},             // nenhum
	}
	for _, c := range cases {
		if got := cfg.ProjectIDForPath(c.query); got != c.want {
			t.Errorf("ProjectIDForPath(%q) = %q, esperava %q", c.query, got, c.want)
		}
	}
}

func TestExpandAlias(t *testing.T) {
	cases := []struct {
		tmpl string
		args []string
		want string
	}{
		{"echo {1}", []string{"hi"}, "echo hi"},
		{"git commit -m {1}", []string{"msg"}, "git commit -m msg"},
		{"echo {*}", []string{"a", "b"}, "echo a b"},
		{"ls", []string{"x"}, "ls x"},   // sem placeholder: anexa
		{"ls", nil, "ls"},               // sem args
		{"run {1} {2}", []string{"a", "b"}, "run a b"},
	}
	for _, c := range cases {
		if got := ExpandAlias(c.tmpl, c.args); got != c.want {
			t.Errorf("ExpandAlias(%q, %v) = %q, esperava %q", c.tmpl, c.args, got, c.want)
		}
	}
}
