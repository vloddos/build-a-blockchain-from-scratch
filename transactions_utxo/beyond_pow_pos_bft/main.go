package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Pick Consensus
//
// For each scenario, recommend a consensus algorithm.
//
// Choices: POW, POS, BFT, RAFT.
//
// Examples:
//
// Public, permissionless, maximize decentralization (Bitcoin-like)   -> POW
// Public, eco-friendly with slashing (Ethereum-like)                  -> POS
// Permissioned consortium of 7 banks                                  -> BFT
// Internal microservice cluster needing strong consistency             -> RAFT
// Public, super-fast finality (Solana-like)                            -> POS
// 3 trusted nodes coordinating leader election                          -> RAFT
// 12 known financial institutions agreeing on settlements              -> BFT

func pickConsensus(s string) string {
	line := strings.ToLower(s)

	if strings.Contains(line, "leader election") ||
		strings.Contains(line, "microservice cluster") ||
		strings.Contains(line, "trusted nodes") ||
		strings.Contains(line, "internal") {
		return "RAFT"
	}

	if strings.Contains(line, "consortium") ||
		strings.Contains(line, "banks") ||
		strings.Contains(line, "financial institutions") ||
		strings.Contains(line, "known") ||
		strings.Contains(line, "settlements") {
		return "BFT"
	}

	if strings.Contains(line, "slashing") ||
		strings.Contains(line, "eco-friendly") ||
		strings.Contains(line, "finality") ||
		strings.Contains(line, "solana") ||
		strings.Contains(line, "stake") ||
		strings.Contains(line, "validator") {
		return "POS"
	}

	if strings.Contains(line, "permissionless") ||
		strings.Contains(line, "public") ||
		strings.Contains(line, "decentralization") ||
		strings.Contains(line, "bitcoin") ||
		strings.Contains(line, "open membership") {
		return "POW"
	}

	return "POW"
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		fmt.Println(pickConsensus(line))
	}
}
