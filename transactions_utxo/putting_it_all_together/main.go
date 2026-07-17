package main

import (
	"bufio"
	"fmt"
	"os"
)

// Capstone Blockchain
//
// Build a mini-blockchain: blocks, hashing, PoW, transactions, merkle root, validation.
//
// Commands:
//
// MINE <data> - mine a block with given data; print height + hash + nonce
// TX <data> - add tx to mempool
// MINE_TX - mine a block including all mempool txs (compute Merkle root); empty mempool
// VALIDATE - validate the entire chain; print VALID or INVALID @ <height>
// Difficulty: 2 leading zeros. Genesis block at height 0 with prev_hash all zeros.
//
// Capstone — combine the lesson primitives.

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
