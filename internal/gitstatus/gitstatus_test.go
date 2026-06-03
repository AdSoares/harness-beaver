package gitstatus

import "testing"

func TestParseCleanWithUpstream(t *testing.T) {
	out := `# branch.oid abc123
# branch.head main
# branch.upstream origin/main
# branch.ab +0 -0
`
	s := parse(out)
	if s.Branch != "main" {
		t.Errorf("branch = %q", s.Branch)
	}
	if s.Dirty {
		t.Error("não deveria estar dirty")
	}
	if !s.HasUpstream {
		t.Error("deveria ter upstream")
	}
	if s.Ahead != 0 || s.Behind != 0 {
		t.Errorf("ahead/behind = %d/%d", s.Ahead, s.Behind)
	}
}

func TestParseDirtyAheadBehind(t *testing.T) {
	out := `# branch.head feature/x
# branch.upstream origin/feature/x
# branch.ab +2 -3
1 .M N... 100644 100644 100644 aaa bbb file.go
? novo.txt
`
	s := parse(out)
	if s.Branch != "feature/x" {
		t.Errorf("branch = %q", s.Branch)
	}
	if !s.Dirty {
		t.Error("deveria estar dirty (mudança + untracked)")
	}
	if s.Ahead != 2 || s.Behind != 3 {
		t.Errorf("ahead/behind = %d/%d, esperava 2/3", s.Ahead, s.Behind)
	}
}

func TestParseDetachedNoUpstream(t *testing.T) {
	out := `# branch.oid abc123
# branch.head (detached)
`
	s := parse(out)
	if !s.Detached {
		t.Error("deveria ser detached")
	}
	if s.HasUpstream {
		t.Error("não deveria ter upstream")
	}
}

func TestLabel(t *testing.T) {
	cases := []struct {
		s    Status
		want string
	}{
		{Status{IsRepo: false}, ""},
		{Status{IsRepo: true, Branch: "main"}, "main ✓"},
		{Status{IsRepo: true, Branch: "main", Dirty: true}, "main ✗"},
		{Status{IsRepo: true, Branch: "dev", Ahead: 2, Behind: 1}, "dev ✓ ↑2 ↓1"},
	}
	for _, c := range cases {
		if got := c.s.Label(); got != c.want {
			t.Errorf("Label(%+v) = %q, esperava %q", c.s, got, c.want)
		}
	}
}
