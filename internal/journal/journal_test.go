package journal

import "testing"

func sampleEntries() []Entry {
	return []Entry{
		{Command: "git status", ProjectID: "repair", ExitCode: 0},
		{Command: "npm test", ProjectID: "repair", ExitCode: 1},
		{Command: "git pull", ProjectID: "beauty", ExitCode: 0},
		{Command: "docker build .", ProjectID: "", ExitCode: 2},
	}
}

func TestFilterByProject(t *testing.T) {
	got := Filter(sampleEntries(), FilterOpts{ProjectID: "repair"})
	if len(got) != 2 {
		t.Fatalf("esperava 2 de repair, veio %d", len(got))
	}
}

func TestFilterFailedOnly(t *testing.T) {
	got := Filter(sampleEntries(), FilterOpts{FailedOnly: true})
	if len(got) != 2 {
		t.Fatalf("esperava 2 com falha, veio %d", len(got))
	}
	for _, e := range got {
		if e.ExitCode == 0 {
			t.Errorf("entrada com exit 0 não deveria passar: %q", e.Command)
		}
	}
}

func TestFilterGrep(t *testing.T) {
	got := Filter(sampleEntries(), FilterOpts{Grep: "GIT"}) // case-insensitive
	if len(got) != 2 {
		t.Fatalf("esperava 2 com 'git', veio %d", len(got))
	}
}

func TestFilterCombined(t *testing.T) {
	got := Filter(sampleEntries(), FilterOpts{ProjectID: "repair", FailedOnly: true})
	if len(got) != 1 || got[0].Command != "npm test" {
		t.Fatalf("esperava só 'npm test', veio %+v", got)
	}
}
