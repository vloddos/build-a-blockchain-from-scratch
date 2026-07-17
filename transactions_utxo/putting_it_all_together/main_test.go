package main

import "testing"

func TestGenesisAndValidation(t *testing.T) {
	chain := newBlockchain()
	if chain.blocks[0].height != 0 {
		t.Fatalf("expected genesis height 0, got %d", chain.blocks[0].height)
	}
	if chain.blocks[0].prevHash != genesisPrevHash {
		t.Fatalf("expected genesis prev hash %q, got %q", genesisPrevHash, chain.blocks[0].prevHash)
	}
	valid, _ := chain.validate()
	if !valid {
		t.Fatal("expected genesis chain to validate")
	}
}

func TestMiningAddsBlockAndKeepsChainValid(t *testing.T) {
	chain := newBlockchain()
	block := chain.mine("hello", nil)
	if block.height != 1 {
		t.Fatalf("expected mined block height 1, got %d", block.height)
	}
	if block.prevHash != chain.blocks[0].hash {
		t.Fatalf("expected prev hash to reference previous block hash")
	}
	valid, _ := chain.validate()
	if !valid {
		t.Fatal("expected chain with mined block to validate")
	}
}
