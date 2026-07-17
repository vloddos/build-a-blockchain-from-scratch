package main

import "testing"

func TestPickConsensus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "pow scenario",
			input:    "Public, permissionless, maximize decentralization (Bitcoin-like)",
			expected: "POW",
		},
		{
			name:     "pos scenario",
			input:    "Public, eco-friendly with slashing (Ethereum-like)",
			expected: "POS",
		},
		{
			name:     "bft scenario",
			input:    "Permissioned consortium of 7 banks",
			expected: "BFT",
		},
		{
			name:     "raft scenario",
			input:    "3 trusted nodes coordinating leader election",
			expected: "RAFT",
		},
		{
			name:     "public finality",
			input:    "Public, super-fast finality (Solana-like)",
			expected: "POS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pickConsensus(tt.input); got != tt.expected {
				t.Fatalf("pickConsensus(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
