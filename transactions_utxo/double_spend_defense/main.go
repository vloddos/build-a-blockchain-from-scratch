package main

import (
	"bufio"
	"fmt"
	"os"
)

// Confirmation Tracking
//
// Track confirmations.
//
// Commands:
//
// BLOCK <height> <txs> - block at height containing comma-separated tx IDs
// TX_CONFIRMATIONS <tx> -> number of blocks deep, or UNCONFIRMED if not in any block
// TIP_HEIGHT -> current tip
// REORG <new_height> - chain reorgs to this height (older blocks removed)
//
// For TX_CONFIRMATIONS: confirmations = tip_height - block_height + 1 (the block containing tx counts as 1).
//
// Example:
// BLOCK 1 a
// BLOCK 2 b,c
// BLOCK 3 d
// TX_CONFIRMATIONS a    -> 3
// TX_CONFIRMATIONS d    -> 1
// TX_CONFIRMATIONS x    -> UNCONFIRMED
// REORG 1
// TX_CONFIRMATIONS a    -> 1
// TX_CONFIRMATIONS d    -> UNCONFIRMED

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
