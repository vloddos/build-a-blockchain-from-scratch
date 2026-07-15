package main

import (
	"bufio"
	"fmt"
	"os"
)

// UTXO Validation
//
// Track UTXOs and validate transactions.
//
// Commands:
//
// MINT <id> <amount> <owner> - create a new UTXO (block reward)
// TX <inputs_csv> -> <outputs_csv> - spend inputs (<utxo_id>:<owner>), create outputs (<id>:<amount>:<owner>)
// BALANCE <owner> - sum unspent UTXOs for owner
// Validation rules:
//
// Each input must exist and not be already spent
// Total inputs >= total outputs (no money creation)
// Skipped: signatures (assume valid for this exercise)
// Output for TX: OK fee=<n> or BAD <reason>.
//
// Example:
// MINT u1 100 alice
// TX u1:alice -> u2:30:bob,u3:65:alice
// BALANCE alice    -> 65
// BALANCE bob      -> 30
// TX u1:alice -> u4:50:bob          -> BAD double_spend
// TX u2:bob -> u5:50:alice          -> BAD insufficient (50 > 30 — wait, 50 > 30 doesn't add up; spend 30 to make 50 = invalid)

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
