package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

// Verify a Merkle proof.

// Use sorted concatenation: at each level, hash(min(left, right) || max(left, right)).

// Input: <tx_hex>|<root_hex>|<sibling1>,<sibling2>,...

// Output: VALID or INVALID.

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			fmt.Println("INVALID")
			continue
		}

		tx := parts[0]
		root := parts[1]
		siblings := strings.Split(parts[2], ",")
		if verifyProof(tx, root, siblings) {
			fmt.Println("VALID")
		} else {
			fmt.Println("INVALID")
		}
	}
}

func verifyProof(tx, root string, siblings []string) bool {
	h := sha256.Sum256([]byte(tx))
	current := fmt.Sprintf("%x", h)

	for _, sibling := range siblings {
		left, right := current, sibling
		if left > right {
			left, right = right, left
		}
		combined := left + right
		h = sha256.Sum256([]byte(combined))
		current = fmt.Sprintf("%x", h)
	}

	return current == root
}
