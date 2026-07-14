package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

// Validate a Chain

// Read blocks; check the chain is valid.

// Each block: <index> <prev_hash> <data> <hash> per line.

// Genesis block has prev_hash = 0 (64 zeros).

// A chain is valid if:

// 1. Block 0 has prev_hash all-zero
// 2. For each i >= 1: block[i].prev_hash == block[i-1].hash
// 3. Each block's hash matches SHA-256 of <index>|<prev_hash>|<data>

// Output: VALID or INVALID block <index> for first issue.

type Block struct {
	index     string
	prev_hash string
	data      string
	hash      string
}

func newBlock(line string) *Block {
	blockData := strings.Split(line, " ")
	index, prev_hash, data, hash := blockData[0], blockData[1], blockData[2], blockData[3]
	return &Block{index, prev_hash, data, hash}
}

const genesisBlockPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"

func validate(prevBlock *Block, block *Block) bool {
	var prevHashEqualityCondition bool

	if prevBlock == nil {
		prevHashEqualityCondition = block.prev_hash == genesisBlockPrevHash
	} else {
		prevHashEqualityCondition = block.prev_hash == prevBlock.hash
	}

	// return prevHashEqualityCondition && block.hash == fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", block.index, block.prev_hash, block.data))))
	return prevHashEqualityCondition && block.hash == fmt.Sprintf("%x", sha256.Sum256(fmt.Appendf(nil, "%s|%s|%s", block.index, block.prev_hash, block.data)))
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)

	var block *Block
	var prevBlock *Block = nil

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		block = newBlock(line)
		if !validate(prevBlock, block) {
			fmt.Println("INVALID block", block.index)
			return
		}
		prevBlock = block
	}
	fmt.Println("VALID")
}
