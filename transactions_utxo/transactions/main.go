package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
// Total inputs >= total outputs (no money creation), the difference is the MINER FEE
// Skipped: signatures (assume valid for this exercise)
// Output for TX: OK fee=<n> or BAD <reason>.
//
// Example:
// MINT u1 100 alice
// TX u1:alice -> u2:30:bob,u3:65:alice
// BALANCE alice    -> 65
// BALANCE bob      -> 30
// TX u1:alice -> u4:50:bob          -> BAD double_spend
// TX u2:bob -> u5:50:alice          -> BAD insufficient

type utxo struct {
	id     string
	amount int
	owner  string
	spent  bool
}

type ledger struct {
	utxos map[string]*utxo
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	lines := make([]string, 0)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	for _, out := range processInput(lines) {
		fmt.Println(out)
	}
}

func processInput(lines []string) []string {
	l := &ledger{utxos: make(map[string]*utxo)}
	outputs := make([]string, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		switch command {
		case "MINT":
			if len(parts) != 4 {
				outputs = append(outputs, "BAD invalid")
				continue
			}
			amount, err := strconv.Atoi(parts[2])
			if err != nil || amount < 0 {
				outputs = append(outputs, "BAD invalid")
				continue
			}
			l.utxos[parts[1]] = &utxo{id: parts[1], amount: amount, owner: parts[3]}
		case "BALANCE":
			if len(parts) != 2 {
				outputs = append(outputs, "BAD invalid")
				continue
			}
			outputs = append(outputs, strconv.Itoa(l.balance(parts[1])))
		case "TX":
			if len(parts) < 2 {
				outputs = append(outputs, "BAD invalid")
				continue
			}
			result, ok := l.processTransaction(trimmed)
			if !ok {
				outputs = append(outputs, result)
				continue
			}
			outputs = append(outputs, result)
		default:
			outputs = append(outputs, "BAD invalid")
		}
	}

	return outputs
}

func (l *ledger) balance(owner string) int {
	sum := 0
	for _, u := range l.utxos {
		if !u.spent && u.owner == owner {
			sum += u.amount
		}
	}
	return sum
}

func (l *ledger) processTransaction(line string) (string, bool) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return "BAD invalid", false
	}

	payload := strings.TrimSpace(parts[1])
	sections := strings.SplitN(payload, "->", 2)
	if len(sections) != 2 {
		return "BAD invalid", false
	}

	inputs, err := parseInputs(strings.TrimSpace(sections[0]))
	if err != nil {
		return "BAD invalid", false
	}

	outputs, err := parseOutputs(strings.TrimSpace(sections[1]))
	if err != nil {
		return "BAD invalid", false
	}

	for _, in := range inputs {
		u, ok := l.utxos[in.id]
		if !ok {
			return "BAD invalid", false
		}
		if u.spent {
			return "BAD double_spend", false
		}
		if u.owner != in.owner {
			return "BAD wrong_owner", false
		}
	}

	inputAmount := 0
	for _, in := range inputs {
		inputAmount += l.utxos[in.id].amount
	}

	outputAmount := 0
	for _, out := range outputs {
		outputAmount += out.amount
	}

	if inputAmount < outputAmount {
		return "BAD insufficient", false
	}

	for _, in := range inputs {
		l.utxos[in.id].spent = true
	}
	for _, out := range outputs {
		l.utxos[out.id] = &utxo{id: out.id, amount: out.amount, owner: out.owner}
	}

	fee := inputAmount - outputAmount
	return fmt.Sprintf("OK fee=%d", fee), true
}

type parsedInput struct {
	id    string
	owner string
}

type parsedOutput struct {
	id     string
	amount int
	owner  string
}

func parseInputs(raw string) ([]parsedInput, error) {
	if raw == "" {
		return nil, fmt.Errorf("empty input")
	}
	items := strings.Split(raw, ",")
	inputs := make([]parsedInput, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		parts := strings.Split(trimmed, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid input")
		}
		inputs = append(inputs, parsedInput{id: parts[0], owner: parts[1]})
	}
	return inputs, nil
}

func parseOutputs(raw string) ([]parsedOutput, error) {
	if raw == "" {
		return nil, fmt.Errorf("empty output")
	}
	items := strings.Split(raw, ",")
	outputs := make([]parsedOutput, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		parts := strings.Split(trimmed, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid output")
		}
		amount, err := strconv.Atoi(parts[1])
		if err != nil || amount < 0 {
			return nil, fmt.Errorf("invalid output amount")
		}
		outputs = append(outputs, parsedOutput{id: parts[0], amount: amount, owner: parts[2]})
	}
	return outputs, nil
}
