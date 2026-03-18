# Provably-Fair-Tool

An interactive **Provably Fair** deck verification tool.  
Enter the shuffle salt (hex) to reproduce the exact card order using the same algorithm as the server (ChaCha8 + Fisher-Yates), and compute SHA256 commitments to verify fairness yourself.

---

## Usage Guide

**No programming environment required** — just download the pre-built executable for your platform and run it.

### Step 1: Download the Executable

Choose the correct file from the `bin` directory based on your operating system:

| OS | Architecture | Filename |
|----|-------------|----------|
| **Windows** | Intel / AMD (most PCs) | `provably-fair-tool-windows-amd64.exe` |
| **macOS** | Apple Silicon (M1/M2/M3/M4) | `provably-fair-tool-darwin-arm64` |
| **macOS** | Intel | `provably-fair-tool-darwin-amd64` |
| **Linux** | x86_64 | `provably-fair-tool-linux-amd64` |
| **Linux** | ARM64 | `provably-fair-tool-linux-arm64` |

### Step 2: Run the Tool

#### Windows

1. Place the downloaded `.exe` file in any folder
2. Hold `Shift` + right-click in the folder → select **Open PowerShell window here** (or Command Prompt)
3. Run:

```
.\provably-fair-tool-windows-amd64.exe
```

#### macOS / Linux

1. Open a terminal and navigate to the directory containing the file
2. Grant execute permission (only needed once):

```bash
chmod +x ./provably-fair-tool-darwin-arm64   # replace with your filename
```

3. Run:

```bash
./provably-fair-tool-darwin-arm64   # replace with your filename
```

> **macOS shows "cannot verify the developer"?** Go to System Settings → Privacy & Security → find the blocked app → click **Allow Anyway**.

You should see the following welcome screen:

```
========== Provably Fair Deck Verification ==========
Enter 'q' to quit

Enter salt (hex, 64 chars):
```

### Step 3: Enter the Salt

Paste the **64-character hex salt** provided by the gaming platform and press Enter.

> **What does a salt look like?** A string of exactly 64 characters containing only `0-9` and `a-f`:  
> `a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2`

### Step 4: Choose Output Format

The program will ask for your preferred output format:

```
Output format [text/json] (default: text):
```

- Press **Enter** directly → human-readable text output (recommended)
- Type `json` and press Enter → JSON output (for programmatic use)

### Step 5: Review the Results

When using `text` format, you will see output similar to:

```
========== Provably Fair Verification ==========

Salt (hex):         a1b2c3d4...(your salt)
Salt Commitment:    e5f6a7b8...(= SHA256(salt), matches the platform's pre-deal commitment)
Shoe Hash:          c9d0e1f2...(= SHA256(card_order), hash of the card sequence)
Total Cards:        416

========== Shuffled Card Order ==========
#1:♠A  [203] #2:♥7  [ 58] #3:♣K  [116] ...(all 416 cards)

========== First 20 Cards ==========
  #1: ♠A (Index=203)
  #2: ♥7 (Index=58)
  ...(first 20 cards in detail)
```

### Step 6: Compare and Verify

Use the output to verify fairness as follows:

| Item | How to Verify |
|------|--------------|
| **Salt Commitment** | Compare the tool's `Salt Commitment` with the hash the platform published **before** dealing — **they must match exactly** |
| **Shoe Hash** | Compare the tool's `Shoe Hash` with the platform's published card-order hash — **they must match exactly** |
| **Card Order** | Compare the tool's card sequence with the actual dealing order to confirm consistency |

> ✅ If all values match, the platform did not cheat — the deal is **provably fair**.  
> ❌ If any value differs, the platform may have tampered with the shuffle.

### Continue or Exit

- After each verification, the program prompts for another salt so you can verify additional rounds.
- Type `q`, `quit`, or `exit` and press Enter to quit.

---

## How It Works

```
Salt
  │
  ├──→ SHA256(salt) = Salt Commitment (published before dealing)
  │
  └──→ ChaCha8(salt) + Fisher-Yates shuffle → 416-card order
                                                  │
                                                  └──→ SHA256(card_order) = Shoe Hash
```

1. Before dealing, the platform publishes the `Salt Commitment` (SHA256 hash of the salt). At this point, players do not know the salt.
2. After the game ends, the platform reveals the `Salt`.
3. Players use this tool to independently reproduce the card order and verify that the commitment and shoe hash match.

Because SHA256 is irreversible, the platform cannot forge a salt after the fact to match a previously published commitment — this guarantees fairness.

---

## FAQ

**Q: I get `salt must be 32 bytes (64 hex chars)`**  
A: The salt you entered has the wrong length. It must be exactly 64 hex characters (`0-9`, `a-f`). Check for extra spaces or missing characters.

**Q: I get `invalid salt hex`**  
A: The salt contains invalid characters. Only `0123456789abcdef` are allowed. Check for any non-hex characters.

**Q: Why 416 cards?**  
A: The shoe uses 8 standard decks, each with 52 cards (4 suits × 13 ranks), totaling 8 × 52 = 416 cards.

**Q: Does this tool require an internet connection?**  
A: No. It runs entirely offline — all computations happen locally on your machine.