package main

import (
	"reflect"
	"testing"
)

func TestProcessInput(t *testing.T) {
	lines := []string{
		"MINT u1 100 alice",
		"TX u1:alice -> u2:30:bob,u3:65:alice",
		"BALANCE alice",
		"BALANCE bob",
	}

	got := processInput(lines)
	want := []string{"OK fee=5", "65", "30"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("processInput() = %v, want %v", got, want)
	}
}

func TestRejectsDoubleSpendAndInsufficientFunds(t *testing.T) {
	lines := []string{
		"MINT u1 100 alice",
		"TX u1:alice -> u2:30:bob,u3:65:alice",
		"TX u1:alice -> u4:50:bob",
		"TX u2:bob -> u5:50:alice",
	}

	got := processInput(lines)
	want := []string{"OK fee=5", "BAD double_spend", "BAD insufficient"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("processInput() = %v, want %v", got, want)
	}
}

func TestRejectsWrongOwner(t *testing.T) {
	lines := []string{
		"MINT u1 50 alice",
		"TX u1:bob -> u2:50:carol",
	}

	got := processInput(lines)
	want := []string{"BAD wrong_owner"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("processInput() = %v, want %v", got, want)
	}
}

func TestInsufficientFunds(t *testing.T) {
	lines := []string{
		"MINT u1 10 alice",
		"TX u1:alice -> u2:50:bob",
	}

	got := processInput(lines)
	want := []string{"BAD insufficient"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("processInput() = %v, want %v", got, want)
	}
}

func TestMultipleInputsAndBalances(t *testing.T) {
	lines := []string{
		"MINT a 30 alice",
		"MINT b 20 alice",
		"TX a:alice,b:alice -> c:45:bob",
		"BALANCE alice",
		"BALANCE bob",
	}

	got := processInput(lines)
	want := []string{"OK fee=5", "0", "45"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("processInput() = %v, want %v", got, want)
	}
}

func TestNestedTransactionsAndBalances(t *testing.T) {
	lines := []string{
		"MINT cb 100 miner",
		"TX cb:miner -> p1:60:alice,p2:35:miner",
		"TX p1:alice -> x:40:bob,y:15:alice",
		"BALANCE bob",
		"BALANCE alice",
		"BALANCE miner",
	}

	got := processInput(lines)
	want := []string{"OK fee=5", "OK fee=5", "40", "15", "35"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("processInput() = %v, want %v", got, want)
	}
}
