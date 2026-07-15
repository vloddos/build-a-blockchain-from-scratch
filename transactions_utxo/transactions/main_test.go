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
