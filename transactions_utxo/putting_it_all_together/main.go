package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
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

const (
	difficulty      = 2
	genesisPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"
	genesisData     = "GENESIS"
)

type block struct {
	height       int
	prevHash     string
	data         string
	transactions []string
	merkleRoot   string
	timestamp    int64
	nonce        int
	hash         string
}

type blockchain struct {
	blocks  []*block
	mempool []string
}

func newBlockchain() *blockchain {
	bc := &blockchain{}
	genesis := bc.mineBlock(genesisData, nil, 0, genesisPrevHash)
	bc.blocks = append(bc.blocks, genesis)
	return bc
}

func (bc *blockchain) mineBlock(data string, txs []string, height int, prevHash string) *block {
	merkleRoot := merkleRootFor(txs)
	b := &block{
		height:       height,
		prevHash:     prevHash,
		data:         data,
		transactions: append([]string(nil), txs...),
		merkleRoot:   merkleRoot,
		timestamp:    time.Now().Unix(),
	}

	prefix := strings.Repeat("0", difficulty)
	for nonce := 0; ; nonce++ {
		b.nonce = nonce
		b.hash = computeBlockHash(b)
		if strings.HasPrefix(b.hash, prefix) {
			return b
		}
	}
}

func (bc *blockchain) appendBlock(data string, txs []string) *block {
	prev := bc.blocks[len(bc.blocks)-1]
	b := bc.mineBlock(data, txs, prev.height+1, prev.hash)
	bc.blocks = append(bc.blocks, b)
	return b
}

func (bc *blockchain) mine(data string, txs []string) *block {
	return bc.appendBlock(data, txs)
}

func (bc *blockchain) validate() (bool, int) {
	if len(bc.blocks) == 0 {
		return false, -1
	}

	if bc.blocks[0].height != 0 {
		return false, 0
	}
	if bc.blocks[0].prevHash != genesisPrevHash {
		return false, 0
	}

	prefix := strings.Repeat("0", difficulty)
	for i, b := range bc.blocks {
		if i == 0 {
			if !strings.HasPrefix(b.hash, prefix) {
				return false, 0
			}
			if b.hash != computeBlockHash(b) {
				return false, 0
			}
			continue
		}

		prev := bc.blocks[i-1]
		if b.prevHash != prev.hash {
			return false, b.height
		}
		if b.height != i {
			return false, b.height
		}
		if !strings.HasPrefix(b.hash, prefix) {
			return false, b.height
		}
		if b.hash != computeBlockHash(b) {
			return false, b.height
		}
	}

	return true, -1
}

func computeBlockHash(b *block) string {
	// header := fmt.Sprintf("%d|%s|%s|%s|%d|%d", b.height, b.prevHash, b.data, b.merkleRoot, b.timestamp, b.nonce)
	header := fmt.Sprintf("%s|%s|%d|%d", b.prevHash, b.merkleRoot, b.timestamp, b.nonce)
	sum := sha256.Sum256([]byte(header))
	return hex.EncodeToString(sum[:])
}

func merkleRootFor(txs []string) string {
	if len(txs) == 0 {
		sum := sha256.Sum256(nil)
		return hex.EncodeToString(sum[:])
	}

	leaves := make([][]byte, 0, len(txs))
	for _, tx := range txs {
		h := sha256.Sum256([]byte(tx))
		leaf := make([]byte, len(h))
		copy(leaf, h[:])
		leaves = append(leaves, leaf)
	}

	for len(leaves) > 1 {
		if len(leaves)%2 == 1 {
			last := leaves[len(leaves)-1]
			leaves = append(leaves, append([]byte(nil), last...))
		}
		next := make([][]byte, 0, len(leaves)/2)
		for i := 0; i < len(leaves); i += 2 {
			combined := append(append([]byte(nil), leaves[i]...), leaves[i+1]...)
			h := sha256.Sum256(combined)
			node := make([]byte, len(h))
			copy(node, h[:])
			next = append(next, node)
		}
		leaves = next
	}

	return hex.EncodeToString(leaves[0])
}

func main() {
	bc := newBlockchain()

	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "MINE":
			data := strings.Join(parts[1:], " ")
			block := bc.appendBlock(data, nil)
			fmt.Printf("height=%d hash=%s nonce=%d\n", block.height, block.hash, block.nonce)
		case "TX":
			data := strings.Join(parts[1:], " ")
			bc.mempool = append(bc.mempool, data)
		case "MINE_TX":
			block := bc.appendBlock(strings.Join(bc.mempool, ","), bc.mempool)
			fmt.Printf("height=%d hash=%s nonce=%d\n", block.height, block.hash, block.nonce)
			bc.mempool = nil
		case "VALIDATE":
			valid, height := bc.validate()
			if valid {
				fmt.Println("VALID")
			} else {
				fmt.Printf("INVALID @ %d\n", height)
			}
		}
	}
}
