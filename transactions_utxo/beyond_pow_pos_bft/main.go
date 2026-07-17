package main

import (
	"bufio"
	"fmt"
	"os"
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

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		fmt.Println("TODO")
	}
}
