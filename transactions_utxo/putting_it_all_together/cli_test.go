package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIExamplesFromComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "validate genesis",
			input:    "VALIDATE\n",
			expected: []string{"VALID"},
		},
		{
			name:     "mine first",
			input:    "MINE first\nVALIDATE\n",
			expected: []string{"height=1", "VALID"},
		},
		{
			name:     "mine txs",
			input:    "TX hello\nTX world\nMINE_TX\nVALIDATE\n",
			expected: []string{"height=1", "VALID"},
		},
		{
			name:     "mine multiple sequential",
			input:    "MINE a\nMINE b\nMINE c\nVALIDATE\n",
			expected: []string{"height=1", "height=2", "height=3", "VALID"},
		},
		{
			name:     "mine txs across blocks",
			input:    "TX one\nMINE_TX\nMINE plain\nTX two\nTX three\nMINE_TX\nVALIDATE\n",
			expected: []string{"height=1", "height=2", "height=3", "VALID"},
		},
		{
			name:     "validate before and after mining",
			input:    "VALIDATE\nMINE only\nVALIDATE\n",
			expected: []string{"VALID", "height=1", "VALID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".")
			cmd.Dir = "."
			cmd.Stdin = strings.NewReader(tt.input)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("go run failed: %v\n%s", err, output)
			}
			text := string(output)
			for _, want := range tt.expected {
				if !strings.Contains(text, want) {
					t.Fatalf("expected output to contain %q\noutput was:\n%s", want, text)
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Ensure tests run from the module directory.
	os.Chdir(".")
	os.Exit(m.Run())
}
