package scanner

import (
	"path/filepath"
	"testing"
)

func TestGroups(t *testing.T) {
	root := t.TempDir()
	cat1 := filepath.Join(root, "cat1")
	cat2 := filepath.Join(root, "cat2")
	leaves := []Candidate{
		{Path: filepath.Join(cat1, "p1")},
		{Path: filepath.Join(cat1, "p2")},
		{Path: filepath.Join(cat2, "p3")},
	}

	groups := Groups(leaves, root)
	got := map[string]int{}
	for _, g := range groups {
		got[g.Path] = g.Count
		if !g.IsGroup {
			t.Errorf("candidato %s deveria ter IsGroup=true", g.Path)
		}
	}
	if got[cat1] != 2 {
		t.Errorf("cat1 count = %d, esperava 2", got[cat1])
	}
	if got[root] != 3 {
		t.Errorf("root count = %d, esperava 3", got[root])
	}
	if _, ok := got[cat2]; ok {
		t.Errorf("cat2 (1 projeto) não deveria ser grupo")
	}
}

func TestGroupsExcludesLeafThatIsAncestor(t *testing.T) {
	root := t.TempDir()
	// um projeto que também é pai de outros (mono-repo): não deve virar "grupo".
	mono := filepath.Join(root, "mono")
	leaves := []Candidate{
		{Path: mono},
		{Path: filepath.Join(mono, "sub1")},
		{Path: filepath.Join(mono, "sub2")},
	}
	for _, g := range Groups(leaves, root) {
		if g.Path == mono {
			t.Errorf("%s é folha; não deveria ser oferecido como grupo", mono)
		}
	}
}
