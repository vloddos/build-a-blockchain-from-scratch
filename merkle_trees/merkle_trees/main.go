package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Merkle Root

// Compute Merkle root of a list of transactions.

// Hash each tx (SHA-256). Build the tree bottom-up: for each pair, hash(left || right).

// If odd count at any level: duplicate the last (Bitcoin-style).

// Input: lines of transaction hex.

// Output: 64-char hex Merkle root.

// Examples:
// 01
//   -> SHA-256("01") = 938db8c9f82c8cb58d3f3ef4fd250036a48d26a712753d2fde5abd03a85cabf4
// 01,02
//   -> hash(SHA(01) || SHA(02)) = f66fc20424b839831cc300edae75967545dcee9bd33306f84c6939b53c2dbe3c
// 01,02,03
//   -> hash(hash(SHA(01) || SHA(02)) || hash(SHA(03) || SHA(03)))   (duplicate odd) = 70d5ca6f37b0c93ce7643a3446a3526de7d433d3c127f8bd3a409776850cf060
// (empty)
//   -> SHA-256("") = e3b0c44298...

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		// Each input line is a comma-separated list of transaction hex strings.
		parts := strings.Split(line, ",")
		root := merkleRoot(parts)
		fmt.Println(hex.EncodeToString(root))

		// fmt.Printf("%x\n", root)
	}
}

func merkleRoot(txHex []string) []byte {
	// Prepare leaves: SHA-256(tx)
	leaves := make([][]byte, 0, len(txHex))
	for _, tx := range txHex {
		// tx = strings.TrimSpace(tx)
		// var data []byte
		// if tx == "" {
		// 	data = []byte{}
		// } else {
		// 	d, err := hex.DecodeString(tx)
		// 	if err != nil {
		// 		// If not valid hex, treat the literal string bytes as input
		// 		data = []byte(tx)
		// 	} else {
		// 		data = d
		// 	}
		// }
		// h := sha256.Sum256(data)

		h := sha256.Sum256([]byte(tx))

		sum := make([]byte, len(h))
		copy(sum, h[:])
		leaves = append(leaves, sum)

		// fmt.Println(hex.EncodeToString(sum)) //debug
	}

	if len(leaves) == 0 {
		// No transactions -> SHA-256("")
		h := sha256.Sum256(nil)
		out := make([]byte, len(h))
		copy(out, h[:])
		return out
	}

	// Build tree bottom-up
	for len(leaves) > 1 {
		if len(leaves)%2 == 1 {
			// duplicate last when odd
			last := leaves[len(leaves)-1]
			dup := make([]byte, len(last))
			copy(dup, last)
			leaves = append(leaves, dup)
		}
		next := make([][]byte, 0, len(leaves)/2)
		for i := 0; i < len(leaves); i += 2 {
			combined := append(leaves[i], leaves[i+1]...)

			// fmt.Println(hex.EncodeToString(combined))              //debug
			// fmt.Println([]byte(hex.EncodeToString(combined)))      //debug
			// fmt.Println(len([]byte(hex.EncodeToString(combined)))) //debug

			h := sha256.Sum256([]byte(hex.EncodeToString(combined)))
			sum := make([]byte, len(h))
			copy(sum, h[:])
			next = append(next, sum)

			// fmt.Println(hex.EncodeToString(sum)) //debug
		}
		leaves = next
	}

	return leaves[0]
}
