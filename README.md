# Blockchain

## Block Structure

A blockchain is a linked list where each node (BLOCK) cryptographically references the previous:

```text
Block N:
  prev_hash:  hash of block N-1
  timestamp:  when this block was created
  merkle_root: hash of all transactions in this block
  nonce:      arbitrary number (used for PoW)
  ... transactions ...

Block N+1:
  prev_hash:  hash of block N
  ...
```

Each block's hash includes its `prev_hash`. To modify block N, you'd need to recompute every subsequent block's hash. The longer the chain, the more expensive — that's the immutability mechanism.

Bitcoin's block header (80 bytes):

```text
4 bytes:  version
32 bytes: prev_block_hash (big-endian)
32 bytes: merkle_root
4 bytes:  timestamp (Unix seconds)
4 bytes:  difficulty bits
4 bytes:  nonce
The hash of the block = SHA-256(SHA-256(header)). The "double SHA-256" is Bitcoin-specific. Other chains (Ethereum) use Keccak-256.
```

A blockchain is fundamentally just data structures + cryptography. No consensus algorithm needed for the data structure itself — that comes in later (PoW, PoS).

## Chain Validation

To verify a chain, walk from genesis (block 0) forward:

```python
def validate(chain):
    for i in range(1, len(chain)):
        if chain[i].prev_hash != hash(chain[i-1]):
            return False
    return True
```

If anyone tampers with block N, its hash changes. Block N+1's `prev_hash` no longer matches. Validation fails at N+1. To "fix": recompute N+1's hash, which changes its own hash, propagating to N+2, ad infinitum.

Tampering thus requires recomputing EVERY hash from N to the tip. With Proof of Work, each hash takes serious compute → tampering becomes prohibitively expensive.

Genesis block: special. It has no `prev_hash` (or all zeros). Hardcoded in the protocol. Bitcoin's genesis block has a famous Times newspaper headline embedded.

Forks: when two miners produce blocks at the same height. Each is a valid chain. Nodes pick the longest (most work) chain — others become orphans and their transactions go back to the mempool.

## Merkle Trees

A **Merkle tree** is a binary tree where:

- Leaves = hashes of data (transactions)
- Internal nodes = hash of their two children's hashes concatenated
- Root = single hash representing all data

```txt
              ROOT
             /    \
           h12    h34
          /  \   /  \
        h1  h2 h3  h4
        |   |  |   |
       tx1 tx2 tx3 tx4
```

Properties:

- **Compact summary**: 32-byte root represents N transactions.
- **Efficient proofs**: prove tx_i is in the tree by giving log(N) sibling hashes — the **Merkle proof**.

Merkle proof for tx2:

```txt
[h1, h34]   (siblings on the path)

verify:
  h2 = hash(tx2)
  h12 = hash(h1 || h2)
  root_computed = hash(h12 || h34)
  return root_computed == known_root
```

Used by:

- **Bitcoin**: SPV clients verify their transactions are in a block without downloading the whole block
- **Git**: commit objects reference tree objects which reference blobs — content-addressed Merkle DAG
- **IPFS**: same Merkle DAG idea
- **Certificate transparency**: log of issued TLS certs; anyone can verify cert is logged
- **Filecoin / Storj**: prove data was stored without revealing it

## Merkle Proofs

Given the root and a transaction, prove the tx is in the tree by providing the SIBLING HASHES along the path:

```txt
        ROOT
       /    \
      A      B
     / \    / \
    C   D  E   F
    |   |  |   |
   t1  t2 t3  t4

Proof for t2:
  [hash(t1), B]      ← siblings, in order

verify(t2, proof, root):
  h = hash(t2)
  for sibling in proof:
    h = hash(min(h, sibling) || max(h, sibling))
  return h == root
```

The order matters! Some impls always concat (sibling || h) when sibling is on left. Bitcoin uses position bit per level. Some use sorted (min, max) — cleaner but loses position info.

Proof size: O(log N) hashes. For 1 million transactions: 20 hashes (640 bytes). Massively smaller than the full data.

This is what makes Bitcoin SPV (Simplified Payment Verification) viable — phone wallets verify their txs in blocks without downloading 600 GB of chain.

Verifiers must know:

- The Merkle root (from block header — small)
- The proof (sibling hashes)
- The transaction itself

Then they can confidently claim "this tx WAS in this block."

## Proof of Work

How do you decide WHO gets to add the next block? In a decentralized system without authorities:

**Proof of Work** (Nakamoto, 2008): you must find a nonce such that `hash(block) < target`. Computing requires lots of trial-and-error.

```python
def mine(block):
    while True:
        if hash(block) < target:
            return block
        block.nonce += 1
```

The `target` controls difficulty. Lower = harder. Bitcoin adjusts target every 2016 blocks (~2 weeks) so blocks come every ~10 minutes regardless of total network hashrate.

Average tries to find a valid hash: `2^256 / target`. For Bitcoin's current difficulty: ~10^25 hashes per block. At 800 EH/s network hashrate, that's ~10 minutes.

Why this works:

- **Honest miners** spend electricity finding valid blocks. They're rewarded with new bitcoins (block reward) + transaction fees.
- **Attackers** would need MORE compute than honest miners combined to mine a fraudulent chain that grows faster than the honest one. Currently impossible at Bitcoin scale.

The 51% attack: if you have >50% of hashrate, you can OUTPACE the honest chain. You can double-spend (spend coins, then mine an alternate chain that excludes the spend). Hard but theoretically possible.

Critique: PoW wastes ENORMOUS amounts of energy. Bitcoin mining ~0.5% of global electricity. Ethereum (now PoS) used to be ~1% but switched.

## Difficulty Adjustment

Bitcoin re-targets difficulty every 2016 blocks (~2 weeks):

```txt
new_target = old_target * (actual_time / expected_time)

Where:
  actual_time = time taken for the last 2016 blocks
  expected_time = 2016 * 10 minutes = 20160 minutes

Capped: new_target can't change by more than 4× in either direction
```

If blocks came faster than 10 min average: actual_time < expected_time → ratio < 1 → target decreases → mining harder.

If slower: target increases → mining easier.

This stabilizes block rate against:

- Hardware advances (ASICs)
- Miners joining/leaving
- Network growth/shrinkage

Bitcoin difficulty has grown ~10^14 since launch. The original Satoshi software target = 2^224. Today: 2^176 or so.

Some chains (Ethereum pre-merge) adjusted every block. Bitcoin's 2016-block window smooths short-term hashrate variance.

Difficulty bombs: Ethereum had a "difficulty bomb" — gradually increasing difficulty to force the chain to upgrade. Several were defused as upgrade timelines slipped.
