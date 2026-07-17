package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Confirmation Tracking
//
// Track confirmations.
//
// Commands:
//
// BLOCK <height> <txs> - block at height containing comma-separated tx IDs
// TX_CONFIRMATIONS <tx> -> number of blocks deep, or UNCONFIRMED if not in any block
// TIP_HEIGHT -> current tip
// REORG <new_height> - chain reorgs to this height (older blocks removed)
//
// For TX_CONFIRMATIONS: confirmations = tip_height - block_height + 1 (the block containing tx counts as 1).
//
// Example:
// BLOCK 1 a
// BLOCK 2 b,c
// BLOCK 3 d
// TX_CONFIRMATIONS a    -> 3
// TX_CONFIRMATIONS d    -> 1
// TX_CONFIRMATIONS x    -> UNCONFIRMED
// REORG 1
// TX_CONFIRMATIONS a    -> 1
// TX_CONFIRMATIONS d    -> UNCONFIRMED

type confirmationState struct {
	blocks map[int][]string
}

func newConfirmationState() *confirmationState {
	return &confirmationState{blocks: make(map[int][]string)}
}

func (s *confirmationState) handleBlock(height int, txIDs []string) {
	s.blocks[height] = txIDs
}

func (s *confirmationState) handleReorg(newHeight int) {
	for height := range s.blocks {
		if height > newHeight {
			delete(s.blocks, height)
		}
	}
}

func (s *confirmationState) tipHeight() int {
	if len(s.blocks) == 0 {
		return 0
	}

	maxHeight := 0
	for height := range s.blocks {
		if height > maxHeight {
			maxHeight = height
		}
	}
	return maxHeight
}

func (s *confirmationState) confirmations(txID string) (int, bool) {
	tip := s.tipHeight()
	if tip == 0 {
		return 0, false
	}

	var matchingHeight int
	found := false
	for height, txs := range s.blocks {
		for _, candidate := range txs {
			if candidate == txID {
				if !found || height > matchingHeight {
					matchingHeight = height
					found = true
				}
			}
		}
	}

	if !found {
		return 0, false
	}

	return tip - matchingHeight + 1, true
}

func processLine(line string, state *confirmationState) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return "", false
	}

	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return "", false
	}

	command := parts[0]
	switch command {
	case "BLOCK":
		if len(parts) < 2 {
			return "", false
		}
		height, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", false
		}
		txText := strings.Join(parts[2:], " ")
		txIDs := parseTxIDs(txText)
		state.handleBlock(height, txIDs)
		return "", false
	case "TX_CONFIRMATIONS":
		if len(parts) < 2 {
			return "", false
		}
		confirmations, ok := state.confirmations(parts[1])
		if !ok {
			return "UNCONFIRMED", true
		}
		return strconv.Itoa(confirmations), true
	case "TIP_HEIGHT":
		return strconv.Itoa(state.tipHeight()), true
	case "REORG":
		if len(parts) < 2 {
			return "", false
		}
		newHeight, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", false
		}
		state.handleReorg(newHeight)
		return "", false
	default:
		return "", false
	}
}

func parseTxIDs(txText string) []string {
	if strings.TrimSpace(txText) == "" {
		return nil
	}

	parts := strings.Split(txText, ",")
	txIDs := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			txIDs = append(txIDs, trimmed)
		}
	}
	return txIDs
}

func main() {
	state := newConfirmationState()
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		output, shouldPrint := processLine(line, state)
		if shouldPrint {
			fmt.Println(output)
		}
	}
}
