package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

// Test 1
//
// Input
// BLOCK 1 a
// BLOCK 2 b,c
// BLOCK 3 d
// TX_CONFIRMATIONS a
// TX_CONFIRMATIONS d
// TX_CONFIRMATIONS x
// REORG 1
// TX_CONFIRMATIONS a
// TX_CONFIRMATIONS d
//
// Expected
// 3
// 1
// UNCONFIRMED
// 1
// UNCONFIRMED

// Test 2
//
// Input
// BLOCK 1 a
// TIP_HEIGHT
//
// Expected
// 1

// Test 3
//
// Input
// BLOCK 1 a,b
// BLOCK 2 c
// BLOCK 3 d,e,f
// TX_CONFIRMATIONS a
// TX_CONFIRMATIONS d
// TX_CONFIRMATIONS f
// TIP_HEIGHT
//
// Expected
// 3
// 1
// 1
// 3

// Test 4
//
// Input
// BLOCK 1 a
// BLOCK 2 b
// REORG 0
// TX_CONFIRMATIONS a
// TIP_HEIGHT
//
// Expected
// UNCONFIRMED
// 0

// Test 5
//
// Input
// BLOCK 5 x
// BLOCK 6 y
// BLOCK 7 z
// TX_CONFIRMATIONS x
// TX_CONFIRMATIONS y
// TX_CONFIRMATIONS z
// REORG 6
// TX_CONFIRMATIONS z
// TX_CONFIRMATIONS y
//
// Expected
// 3
// 2
// 1
// UNCONFIRMED
// 1

// Test 6
//
// Input
// TIP_HEIGHT
// TX_CONFIRMATIONS unknown
//
// Expected
// 0
// UNCONFIRMED

func runScenario(input string) []string {
	state := newConfirmationState()
	var outputs []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		output, shouldPrint := processLine(line, state)
		if shouldPrint {
			outputs = append(outputs, output)
		}
	}

	return outputs
}

func TestConfirmationTrackingScenarios(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "test 1",
			input: strings.Join([]string{
				"BLOCK 1 a",
				"BLOCK 2 b,c",
				"BLOCK 3 d",
				"TX_CONFIRMATIONS a",
				"TX_CONFIRMATIONS d",
				"TX_CONFIRMATIONS x",
				"REORG 1",
				"TX_CONFIRMATIONS a",
				"TX_CONFIRMATIONS d",
			}, "\n"),
			expected: []string{"3", "1", "UNCONFIRMED", "1", "UNCONFIRMED"},
		},
		{
			name: "test 2",
			input: strings.Join([]string{
				"BLOCK 1 a",
				"TIP_HEIGHT",
			}, "\n"),
			expected: []string{"1"},
		},
		{
			name: "test 3",
			input: strings.Join([]string{
				"BLOCK 1 a,b",
				"BLOCK 2 c",
				"BLOCK 3 d,e,f",
				"TX_CONFIRMATIONS a",
				"TX_CONFIRMATIONS d",
				"TX_CONFIRMATIONS f",
				"TIP_HEIGHT",
			}, "\n"),
			expected: []string{"3", "1", "1", "3"},
		},
		{
			name: "test 4",
			input: strings.Join([]string{
				"BLOCK 1 a",
				"BLOCK 2 b",
				"REORG 0",
				"TX_CONFIRMATIONS a",
				"TIP_HEIGHT",
			}, "\n"),
			expected: []string{"UNCONFIRMED", "0"},
		},
		{
			name: "test 5",
			input: strings.Join([]string{
				"BLOCK 5 x",
				"BLOCK 6 y",
				"BLOCK 7 z",
				"TX_CONFIRMATIONS x",
				"TX_CONFIRMATIONS y",
				"TX_CONFIRMATIONS z",
				"REORG 6",
				"TX_CONFIRMATIONS z",
				"TX_CONFIRMATIONS y",
			}, "\n"),
			expected: []string{"3", "2", "1", "UNCONFIRMED", "1"},
		},
		{
			name: "test 6",
			input: strings.Join([]string{
				"TIP_HEIGHT",
				"TX_CONFIRMATIONS unknown",
			}, "\n"),
			expected: []string{"0", "UNCONFIRMED"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runScenario(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("unexpected output\ninput:\n%s\nexpected: %q\nactual:   %q", tt.input, tt.expected, got)
			}
		})
	}
}
