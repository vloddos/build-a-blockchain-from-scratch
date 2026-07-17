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

// Test 1
//
// Input
// VALIDATE
//
// Expected
// VALID

// Test 2
//
// Input
// MINE first
// VALIDATE
//
// Expected
// height=1 hash=00a0fa5d7187b9320039d36ea0b324997fe9a5d890061551b67c612a219844d0 nonce=18
// VALID

// Test 3
//
// Input
// TX hello
// TX world
// MINE_TX
// VALIDATE
//
// Expected
// height=1 hash=00628d5b1f4f8d6ef536e4076fce9e6314ac1d921e80cadce91d092693680483 nonce=460
// VALID

// Test 4
//
// Input
// MINE a
// MINE b
// MINE c
// VALIDATE
//
// Expected
// height=1 hash=009d3ea2cc2ce654abb8c91cc0452a01d6aa9e0c03287daf9da5b7c5f1c96304 nonce=774
// height=2 hash=0006701e8eb8060579ec46b6b4a35275273e0aff1876b1ed77347997b26788f6 nonce=237
// height=3 hash=00b6a36905104ad1ceba41422b0ab2fa32b55d296ba27799132fa76ed6612bf0 nonce=37
// VALID

// Test 5
//
// Input
// TX one
// MINE_TX
// MINE plain
// TX two
// TX three
// MINE_TX
// VALIDATE
//
// Expected
// height=1 hash=006c11d0f8f2fb9e91ffda0ae0bef8ea83406148a76958e500e2bdffc547b4ee nonce=157
// height=2 hash=001ebd69ab3e27caf7c375a34fb4353e2ae504b87e4d449510b5940c191036da nonce=61
// height=3 hash=00f06463b3120da26793a157f7fc4c497c808194e41bef313e7135a7cffcef58 nonce=132
// VALID

// Test 6
//
// Input
// VALIDATE
// MINE only
// VALIDATE
//
// Expected
// VALID
// height=1 hash=00bef4f0b4845647c0b9a762dc23da76746261a840909e1b816ba196081a838d nonce=1
// VALID
