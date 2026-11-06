package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSubmissionLogAppends(t *testing.T) {
	path := filepath.Join(t.TempDir(), "submissions.jsonl")

	want := []submission{
		{Time: time.Date(2026, 11, 6, 9, 0, 0, 0, time.UTC), Name: "A", Phone: "253-555-0100", Email: "a@example.com", Address: "1 Main St, Bonney Lake", Service: "Gutter Installation", Message: "line one\nline two"},
		{Time: time.Date(2026, 11, 6, 9, 5, 0, 0, time.UTC), Name: "B", Phone: "253-555-0101", Email: "b@example.com", Address: "2 Main St, Sumner", Service: "Gutter Cleaning"},
	}

	// Two opens: a restart must append, not truncate.
	for _, s := range want {
		l, err := openSubmissionLog(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := l.append(s); err != nil {
			t.Fatal(err)
		}
		l.f.Close()
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("file mode = %o, want 600", perm)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var got []submission
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var s submission
		if err := json.Unmarshal(sc.Bytes(), &s); err != nil {
			t.Fatalf("line %q: %v", sc.Text(), err)
		}
		got = append(got, s)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestSubmissionLogDisabled(t *testing.T) {
	l, err := openSubmissionLog("")
	if err != nil || l != nil {
		t.Fatalf("empty path: got %v, %v; want nil, nil", l, err)
	}
	if err := l.append(submission{Name: "x"}); err != nil {
		t.Fatalf("nil log append: %v", err)
	}
}

func TestSubmissionLogBadPath(t *testing.T) {
	if _, err := openSubmissionLog(filepath.Join(t.TempDir(), "missing-dir", "s.jsonl")); err == nil {
		t.Fatal("expected error for unwritable path")
	}
}
