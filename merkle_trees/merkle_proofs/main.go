package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO (merkle-proof): implement per the lesson description.

// Verify a Merkle proof.

// Use sorted concatenation: at each level, hash(min(left, right) || max(left, right)).

// Input: <tx_hex>|<root_hex>|<sibling1>,<sibling2>,...

// Output: VALID or INVALID.

// Example:
//
// 01|<root>|<sibling_for_position_1>,<sibling_for_position_0>
//
// If hash(01) combined with siblings reaches root: VALID. Else INVALID.

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
