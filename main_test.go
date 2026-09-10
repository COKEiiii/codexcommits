package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type testRepo struct {
	path string
	head string
}

func newTestRepo(t *testing.T) testRepo {
	t.Helper()
	path := t.TempDir()
	runGit(t, path, "init", "-q")
	runGit(t, path, "config", "user.name", "Commit Test")
	runGit(t, path, "config", "user.email", "test@example.invalid")
	runGit(t, path, "config", "commit.gpgsign", "false")
	runGit(t, path, "config", "core.hooksPath", filepath.Join(path, ".git", "hooks"))
	writeFile(t, filepath.Join(path, "main.py"), "value = 1\n", 0o644)
	runGit(t, path, "add", "main.py")
	runGit(t, path, "commit", "-qm", "chore: baseline")
	writeFile(t, filepath.Join(path, "main.py"), "value = 2\n", 0o644)
	runGit(t, path, "add", "main.py")
	return testRepo{path: path, head: runGit(t, path, "rev-parse", "HEAD")}
}

func runGit(t *testing.T, repo string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func writeFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func withRepo(t *testing.T, repo string, input string, generator func([]byte, options, ioWriter) (string, error)) (string, string, error) {
	t.Helper()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldDir) })
	oldFind, oldGenerate := findExecutable, generateCommitMessage
	findExecutable = func(string) (string, error) { return "/mock/tool", nil }
	generateCommitMessage = func(diff []byte, opts options, out io.Writer) (string, error) {
		return generator(diff, opts, writerAdapter{out})
	}
	t.Cleanup(func() { findExecutable, generateCommitMessage = oldFind, oldGenerate })
	var stdout, stderr bytes.Buffer
	err := run(nil, strings.NewReader(input), &stdout, &stderr)
	return stdout.String(), stderr.String(), err
}

type ioWriter interface{ Write([]byte) (int, error) }
type writerAdapter struct{ io.Writer }

func fixedGenerator(message string) func([]byte, options, ioWriter) (string, error) {
	return func([]byte, options, ioWriter) (string, error) { return message, nil }
}

func TestOnlyStagedContentIsSentAndCommitted(t *testing.T) {
	repo := newTestRepo(t)
	writeFile(t, filepath.Join(repo.path, "main.py"), "value = 3  # unstaged marker\n", 0o644)
	var received []byte
	_, _, err := withRepo(t, repo.path, "y\n", func(diff []byte, _ options, _ ioWriter) (string, error) {
		received = append([]byte(nil), diff...)
		return "fix: update the value", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(received, []byte("+value = 2")) || bytes.Contains(received, []byte("unstaged marker")) {
		t.Fatalf("unexpected model diff:\n%s", received)
	}
	if got := runGit(t, repo.path, "show", "HEAD:main.py"); got != "value = 2" {
		t.Fatalf("committed %q", got)
	}
}

func TestCancelPreservesHeadAndIndex(t *testing.T) {
	repo := newTestRepo(t)
	tree := runGit(t, repo.path, "write-tree")
	_, _, err := withRepo(t, repo.path, "\n", fixedGenerator("fix: update the value"))
	if err != nil {
		t.Fatal(err)
	}
	if repo.head != runGit(t, repo.path, "rev-parse", "HEAD") || tree != runGit(t, repo.path, "write-tree") {
		t.Fatal("cancel changed HEAD or index")
	}
}

func TestEditIsLiteralAndRequiresConfirmation(t *testing.T) {
	repo := newTestRepo(t)
	message := "fix: preserve $(touch SHOULD_NOT_EXIST) literally"
	_, _, err := withRepo(t, repo.path, "e\n"+message+"\ny\n", fixedGenerator("fix: update the value"))
	if err != nil {
		t.Fatal(err)
	}
	if got := runGit(t, repo.path, "log", "-1", "--format=%s"); got != message {
		t.Fatalf("got %q", got)
	}
	if _, err := os.Stat(filepath.Join(repo.path, "SHOULD_NOT_EXIST")); !os.IsNotExist(err) {
		t.Fatal("message was interpreted by a shell")
	}
}

func TestIndexChangeInvalidatesResult(t *testing.T) {
	repo := newTestRepo(t)
	_, _, err := withRepo(t, repo.path, "", func([]byte, options, ioWriter) (string, error) {
		writeFile(t, filepath.Join(repo.path, "extra.txt"), "new", 0o644)
		runGit(t, repo.path, "add", "extra.txt")
		return "fix: update the value", nil
	})
	if err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("expected changed error, got %v", err)
	}
}

func TestHeadChangeInvalidatesResult(t *testing.T) {
	repo := newTestRepo(t)
	_, _, err := withRepo(t, repo.path, "", func([]byte, options, ioWriter) (string, error) {
		runGit(t, repo.path, "commit", "--allow-empty", "-qm", "chore: concurrent commit")
		return "fix: update the value", nil
	})
	if err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("expected changed error, got %v", err)
	}
}

func TestFailingHookIsRespected(t *testing.T) {
	repo := newTestRepo(t)
	hook := filepath.Join(repo.path, ".git", "hooks", "pre-commit")
	writeFile(t, hook, "#!/bin/sh\nexit 1\n", 0o755)
	_, _, err := withRepo(t, repo.path, "y\n", fixedGenerator("fix: update the value"))
	if err == nil || !strings.Contains(err.Error(), "git commit failed") {
		t.Fatalf("expected hook failure, got %v", err)
	}
	if repo.head != runGit(t, repo.path, "rev-parse", "HEAD") {
		t.Fatal("failed hook created a commit")
	}
}

func TestEmptyIndexNeverCallsGenerator(t *testing.T) {
	repo := newTestRepo(t)
	runGit(t, repo.path, "reset", "-q", "HEAD")
	called := false
	_, _, err := withRepo(t, repo.path, "", func([]byte, options, ioWriter) (string, error) {
		called = true
		return "fix: impossible", nil
	})
	if err == nil || !strings.Contains(err.Error(), "no staged changes") || called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}

func TestGenerationFailureDoesNotCommit(t *testing.T) {
	repo := newTestRepo(t)
	_, _, err := withRepo(t, repo.path, "", func([]byte, options, ioWriter) (string, error) {
		return "", fmt.Errorf("unavailable")
	})
	if err == nil || !strings.Contains(err.Error(), "unavailable") {
		t.Fatalf("got %v", err)
	}
	if repo.head != runGit(t, repo.path, "rev-parse", "HEAD") {
		t.Fatal("generation failure created a commit")
	}
}

func TestValidateSubject(t *testing.T) {
	for _, good := range []string{"feat: add parser", "fix(api): handle empty input", "feat!: remove old API"} {
		if err := validateSubject(good); err != nil {
			t.Errorf("%q: %v", good, err)
		}
	}
	for _, bad := range []string{"Update files", "feat: ", "feat: line one\nline two"} {
		if err := validateSubject(bad); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

func TestSanitizedEnvironment(t *testing.T) {
	t.Setenv("CODEX_API_KEY", "secret")
	t.Setenv("OPENAI_API_KEY", "secret")
	t.Setenv("OPENAI_BASE_URL", "https://example.invalid")
	t.Setenv("SAFE_MARKER", "kept")
	joined := strings.Join(sanitizedEnvironment(), "\n")
	if strings.Contains(joined, "secret") || strings.Contains(joined, "OPENAI_BASE_URL") || !strings.Contains(joined, "SAFE_MARKER=kept") {
		t.Fatalf("unexpected environment: %s", joined)
	}
}

func TestOptionsUseCodexDefaultModel(t *testing.T) {
	t.Setenv("CODEXCOMMITS_MODEL", "")
	opts, _, err := parseOptions(nil, &bytes.Buffer{})
	if err != nil || opts.model != "" || opts.timeout != 180*time.Second {
		t.Fatalf("opts=%+v err=%v", opts, err)
	}
}

func TestHelpIsSuccessful(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := run([]string{"-h"}, strings.NewReader(""), &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(errOut.String(), "Usage: codexcommits") {
		t.Fatalf("missing help: %s", errOut.String())
	}
}
