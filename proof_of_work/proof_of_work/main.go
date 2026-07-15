package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Mine a Block
//
// Find a nonce such that SHA-256 of <data>:<nonce> starts with <difficulty> zero hex chars.
//
// Input: <data>|<difficulty> per line.
// Output: nonce (smallest non-negative integer).
//
// Examples:
// hello|2
//   -> Some nonce N where SHA-256("hello:N") starts with "00"
//
// hello|0
//   -> 0   (any nonce works; "0" is smallest)
//
//Test by iterating: nonce=0,1,2,... until hash matches.

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}

		data := parts[0]
		difficulty, err := strconv.Atoi(parts[1])
		if err != nil || difficulty < 0 {
			continue
		}

		prefix := strings.Repeat("0", difficulty)
		for nonce := 0; ; nonce++ {
			txt := fmt.Sprintf("%s:%d", data, nonce)
			hash := sha256.Sum256([]byte(txt))
			hexHash := hex.EncodeToString(hash[:])
			if strings.HasPrefix(hexHash, prefix) {
				fmt.Println(nonce)
				break
			}
		}
	}
}
