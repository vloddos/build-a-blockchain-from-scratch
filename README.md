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
