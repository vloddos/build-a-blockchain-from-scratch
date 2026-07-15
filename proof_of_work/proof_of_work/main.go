package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO (pow): implement per the lesson description.

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
		fmt.Println("TODO")
	}
}
