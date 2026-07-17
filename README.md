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

## Transactions & UTXO

### Transactions

A **transaction** moves value from inputs to outputs. Bitcoin uses the **UTXO** model:

- **Output (UTXO)**: a chunk of bitcoin owned by some address. Spendable until consumed.
- **Input**: references a previous output to spend it. Consumed once used.

```txt
Tx1 outputs:
  output_0:  0.5 BTC -> Alice
  output_1:  1.5 BTC -> Bob

Tx2 inputs:
  input_0: refers to (Tx1, output_0)   -- spend Alice's 0.5

Tx2 outputs:
  output_0: 0.3 BTC -> Charlie
  output_1: 0.2 BTC -> Alice  (change)
```

The total inputs must >= total outputs. The difference is the MINER FEE.

To prove the input is YOURS, you sign with the private key matching the output's address. Bitcoin uses ECDSA over secp256k1.

Validation:

- Every input is signed validly
- Every input refers to an unspent output
- Total inputs ≥ total outputs (fees go to miner)
- Other rules: locktime, scripts, etc.

Ethereum uses an account model instead: each address has a balance. Transactions debit one balance, credit another. Simpler but requires more state per node.

Bitcoin's UTXO is more privacy-friendly (each output is independent) but harder to track balance. Modern wallets aggregate via address scanning.

### Double-Spend Defense

Without a central authority, what stops me from spending the same coin twice?

```txt
I have 1 BTC.
Tx_A: send 1 BTC to merchant_A
Tx_B: send 1 BTC to merchant_B    (different output, same UTXO!)

Both broadcast at once. Which is real?
```

Answer: the chain decides. Whichever tx is included in the LONGER VALID CHAIN wins. The other is rejected (its input was already spent).

But what about temporary forks? Block 100 mined by miner X contains tx_A. Block 100 mined by miner Y contains tx_B. Both are valid; both extend block 99.

The chain "RESOLVES" when block 101 builds on top of one of them. The other becomes an orphan.

**Confirmations**: how many blocks deep your tx is. After 6 confirmations (Bitcoin convention), the chain is considered final — reorganizing 6 blocks would require >50% hashrate AND lots of luck.

Merchant policy:

- Coffee: 0 confirmations (low value, low risk)
- Online purchase: 1 confirmation (~10 min)
- Large transfer: 6 confirmations (~1 hour)
- Custodial deposit: 100+ confirmations (paranoid)

This is why Bitcoin transactions are slow. Other chains optimize:

- Litecoin: 2.5 min blocks (4× faster)
- Solana: ~400ms blocks (different consensus)
- Lightning Network: off-chain, instant

Trade-off: faster confirmation = more risk of reorg.

### Beyond PoW: PoS, BFT

Proof of Work is one point in a design space, optimized for one extreme: *nobody knows who the participants are, and anyone may join or leave at will*. Relax that assumption and much cheaper consensus becomes available. An engineer's job is matching the algorithm to the trust model — which is literally this lesson's exercise.

#### Proof of Stake

Replace "one hash, one lottery ticket" with "one staked coin, one lottery ticket." Validators lock coins as collateral; the protocol pseudo-randomly selects a proposer, weighted by stake; provable misbehavior (like signing two conflicting blocks) gets the stake **slashed**. Security still costs something an attacker must acquire — but it's capital at risk instead of burned electricity. Ethereum switched in the 2022 Merge, cutting its energy use by ~99.95%; Cardano, Solana, Polkadot, and Cosmos chains are PoS-native. The honest trade-offs: the classic *nothing-at-stake* objection (voting on every fork is free unless slashing punishes it — that's precisely what slashing is for), long-range rewrites of old history using once-valid keys (countered by checkpointing / "weak subjectivity"), and a rich-get-richer tilt, since stake earns yield.

#### The BFT family

Byzantine Fault Tolerance comes from a different lineage — distributed-systems theory, not cryptocurrency. Lamport, Shostak, and Pease formalized the Byzantine generals problem in 1982; Castro and Liskov's **PBFT** (1999) made it practical: among `n = 3f + 1` known validators, up to `f` can be arbitrarily malicious and the rest still agree, with **immediate finality** — a committed block is final, no confirmations, no reorgs, ever. Cost: validators are a fixed, known set, and communication is O(n²) messages per decision, so counts stay in the dozens-to-hundreds. Descendants refine it: **Tendermint** (started by Jae Kwon in 2014, later developed with Ethan Buchman; the engine of Cosmos) integrates BFT with PoS validator selection; **HotStuff** (2019) got messages to O(n) and powered Facebook's Diem. Modern "PoS" chains are mostly hybrids: PoS chooses *who may vote*, a BFT-style protocol decides *how votes commit*.

#### And plain old Raft

If nodes can crash but never lie — one company's internal cluster — Byzantine defenses are pure overhead. **Raft** (and Paxos) tolerate crash faults only: `f` failures among `2f + 1` nodes, a simple elected leader, no signatures, no adversary model. This is etcd, Consul, and every "strongly consistent" database you've deployed. Using PoW for a microservice cluster is malpractice; so is using Raft between mutually distrusting banks.

#### The decision rubric

|Trust model|Pick|
|-----------|----|
|Open membership, anyone can join, maximize censorship-resistance|POW|
|Open membership, capital-based security, slashing, energy-efficient|POS|
|Known members who don't trust each other (consortium of banks)|BFT|
|Known members you fully trust (your own infra), need consistency|RAFT|

Two axes decide everything: *is membership open or closed*? and *can participants be malicious or merely crash*? Open + malicious → PoW/PoS. Closed + malicious → BFT. Closed + honest-but-crashy → Raft.

#### Your exercise

Each input line describes a scenario; print exactly one token: `POW`, `POS`, `BFT`, or `RAFT` (uppercase). The graded confusion is BFT vs RAFT — both run among known nodes, so scan for *trust*: "Permissioned consortium of 7 banks" and "12 known financial institutions" are mutually distrusting parties → `BFT`; "Internal microservice cluster needing strong consistency" and "3 trusted nodes coordinating leader election" are one administrative domain → `RAFT`. Second catch: "Public, super-fast finality (Solana-like)" is `POS` — fast public finality is a stake-plus-BFT hybrid, not raw BFT, because the validator set is open.

### Putting It All Together

You've built every layer that makes a real blockchain work:

|Layer|Lesson|
|-----|------|
|Block structure|Block Structure|
|Chain validation|Chain Validation|
|Merkle trees|Merkle Trees|
|Merkle proofs (SPV)|Merkle Proofs|
|Proof of Work|Proof of Work|
|Difficulty retargeting|Difficulty Adjustment|
|UTXO transactions|Transactions|
|Double-spend defense|Double-Spend Defense|
|Consensus alternatives|Beyond PoW: PoS, BFT|
|**Real Bitcoin block hashing**|Bitcoin Block Header|
|**ECDSA on secp256k1**|ECDSA on secp256k1|
|**HD wallets (BIP32)**|BIP32 HD Wallets|
|**Bitcoin Script**|Bitcoin Script|
|**Mempool & fee selection**|Mempool & Fee Selection|
|**Fork choice / heaviest chain**|Fork Resolution|
|**P2P gossip**|P2P Gossip Propagation|

Every column in the table corresponds to a real component in Bitcoin Core, Geth, Reth, or any major node implementation. You've hand-built each in Go and verified your code against canonical test vectors (Bitcoin's actual genesis hash, the secp256k1 generator point, sorted-Merkle proofs, etc).

#### Reference implementations to study next

- **Bitcoin Core** (`bitcoin/bitcoin`, ~1M lines of C++) — the reference. Read `src/validation.cpp` for chain validation; `src/script/interpreter.cpp` for Script.
- **btcd** (`btcsuite/btcd`, Go) — friendlier to read than Bitcoin Core if you prefer Go.
- **rust-bitcoin** (`rust-bitcoin/rust-bitcoin`) — clean Rust crates for the primitives (`hashes`, `bitcoin`, `secp256k1`).
- **Geth** (`ethereum/go-ethereum`, Go) — most popular Ethereum client.
- **Reth** (`paradigmxyz/reth`, Rust) — modern Ethereum rewrite, very clean.
- **Solana validator** (`solana-labs/solana`, Rust) — PoH + PoS, very different trade-offs.

#### What you've intentionally skipped

These are the next big topics if you want to go deeper:

- **Smart contracts**: Ethereum's EVM (Geth, Reth) or Solana's BPF VM. Programs that run on-chain with deterministic gas metering.
- **SegWit / Taproot** (BIP-141, BIP-341): newer Bitcoin script formats with Schnorr signatures and Merkleized script trees.
- **Layer 2**: Lightning Network (off-chain channels), Optimism / Arbitrum (optimistic rollups), zkSync / Starknet (zk-rollups).
- **Cross-chain bridges**: Wormhole, IBC. Move assets between chains. Frequently hacked — bridge security is an open research area.
- **Privacy coins**: Monero (ring signatures + RingCT + stealth addresses), Zcash (zk-SNARKs).
- **MEV (Maximal Extractable Value)**: extracting value via transaction ordering. A multi-billion-dollar industry built on mempool observation.
- **PBFT / HotStuff / Tendermint**: deeper dive into BFT consensus families for permissioned and PoS networks.

#### Where blockchain shines

- Permissionless value transfer (Bitcoin)
- Programmable money / DeFi (Ethereum)
- Censorship resistance (places where governments restrict speech / finance)
- Auditable supply chains (some niches)
- Bearer instruments (NFTs, regardless of opinion)

#### Where blockchain doesn't

- High throughput (use a database)
- Privacy (most chains are public ledgers; opt-in privacy is complex)
- Cheap fees during demand spikes (gas wars)
- Reversibility (send to the wrong address = gone)
- Complex business logic (smart-contract bugs cost billions yearly)

You now understand the foundations that every blockchain shares. The exercise below is your capstone — chain together blocks, transactions, mempool, and validation into one mini-blockchain.
