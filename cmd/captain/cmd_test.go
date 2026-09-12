package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// captureStdout temporarily redirects os.Stdout while f runs, returning
// everything written to it.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w

	f()

	os.Stdout = old
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestHandleCmdVersion builds the full cobra command tree via handleCmd and
// dispatches the harmless "version" subcommand, which needs neither docker
// nor network access - unlike build/test/push/pull/purge/self-update.
func TestHandleCmdVersion(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"captain", "version"}

	output := captureStdout(t, func() {
		handleCmd()
	})

	if output != "v1.1.3\n" {
		t.Fatalf("expected version output %q, got %q", "v1.1.3\n", output)
	}
}

func TestGetNamespaceUsesUserEnv(t *testing.T) {
	oldUser, hadUser := os.LookupEnv("USER")
	defer func() {
		if hadUser {
			os.Setenv("USER", oldUser)
		} else {
			os.Unsetenv("USER")
		}
	}()

	os.Setenv("USER", "captain-test-user")
	if got := getNamespace(); got != "captain-test-user" {
		t.Fatalf("expected namespace %q, got %q", "captain-test-user", got)
	}
}
