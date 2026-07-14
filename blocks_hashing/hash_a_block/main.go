package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
)

// Hash a Block

// Compute a block's hash.

// Format the header as concatenation: <prev_hash>|<merkle_root>|<timestamp>|<nonce>.

// Hash with SHA-256 (single, not double — simpler for the lesson).

// Input: lines of <prev_hash>|<merkle_root>|<timestamp>|<nonce>.

// Output: 64-char hex SHA-256.

// Example:
// 0000000000000000000000000000000000000000000000000000000000000000|abcd|1700000000|0
//   -> SHA-256 of "0000...|abcd|1700000000|0" = some 64-char hex

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		fmt.Printf("%x\n", sha256.Sum256([]byte(line)))
	}
}
