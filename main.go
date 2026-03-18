// verify_deck - Interactive Provably Fair deck verification tool.
// Takes a shuffle salt (hex) and reproduces the exact card order using the same
// algorithm as the server (ChaCha8 + Fisher-Yates), then computes SHA256 commitments.
//
// Usage:
//
//	go run tools/verify_deck/main.go
//
// Verification steps:
//  1. SHA256(salt) == salt_commitment
//  2. Reproduce card order with salt + ChaCha8 + Fisher-Yates
//  3. SHA256(card_order) == chain_head
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mathrand "math/rand/v2"
	"os"
	"strings"
)

type Card struct {
	Suit  uint8  // 0=♠, 1=♥, 2=♣, 3=♦
	Value uint8  // 1=A, 2-9, 10, 11=J, 12=Q, 13=K
	Index uint16 // unique index (0-415)
}

var suits = []string{"♠", "♥", "♣", "♦"}
var values = []string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func cardToString(c Card) string {
	if c.Suit < 4 && c.Value >= 1 && c.Value <= 13 {
		return suits[c.Suit] + values[c.Value]
	}
	return "??"
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("========== Provably Fair Deck Verification ==========")
	fmt.Println("Enter 'q' to quit")
	fmt.Println()

	for {
		fmt.Print("Enter salt (hex, 64 chars): ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if strings.ToLower(input) == "q" || strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			fmt.Println("Bye!")
			break
		}

		jsonOutput := false
		fmt.Print("Output format [text/json] (default: text): ")
		if scanner.Scan() {
			format := strings.TrimSpace(strings.ToLower(scanner.Text()))
			jsonOutput = (format == "json" || format == "j")
		}

		if err := verify(input, jsonOutput); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Println()
	}
}

func verify(saltHex string, jsonOutput bool) error {
	saltBytes, err := hex.DecodeString(saltHex)
	if err != nil {
		return fmt.Errorf("invalid salt hex: %v", err)
	}
	if len(saltBytes) != 32 {
		return fmt.Errorf("salt must be 32 bytes (64 hex chars), got %d bytes", len(saltBytes))
	}

	var salt [32]byte
	copy(salt[:], saltBytes)

	// ===== 1. Create 8 decks (416 cards) =====
	cards := make([]Card, 0, 416)
	idx := uint16(0)
	for range 8 {
		for suit := range uint8(4) {
			for value := uint8(1); value <= 13; value++ {
				cards = append(cards, Card{Suit: suit, Value: value, Index: idx})
				idx++
			}
		}
	}

	// ===== 2. ChaCha8 + Fisher-Yates shuffle =====
	rng := mathrand.New(mathrand.NewChaCha8(salt))
	for i := len(cards) - 1; i > 0; i-- {
		j := rng.IntN(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}

	// ===== 3. Compute commitments =====
	// C = SHA256(salt)
	saltHash := sha256.Sum256(salt[:])
	saltCommitment := hex.EncodeToString(saltHash[:])

	// H = SHA256(card_order)
	h := sha256.New()
	cardIndices := make([]uint16, len(cards))
	cardNames := make([]string, len(cards))
	for i, card := range cards {
		var buf [2]byte
		binary.BigEndian.PutUint16(buf[:], card.Index)
		h.Write(buf[:])
		cardIndices[i] = card.Index
		cardNames[i] = cardToString(card)
	}
	shoeHash := hex.EncodeToString(h.Sum(nil))

	// ===== 4. Output results =====
	if jsonOutput {
		result := map[string]interface{}{
			"salt":            saltHex,
			"salt_commitment": saltCommitment,
			"shoe_hash":       shoeHash,
			"card_count":      len(cards),
			"card_indices":    cardIndices,
			"card_names":      cardNames,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return nil
	}

	fmt.Println()
	fmt.Println("========== Provably Fair Verification ==========")
	fmt.Println()
	fmt.Printf("Salt (hex):         %s\n", saltHex)
	fmt.Printf("Salt Commitment:    %s\n", saltCommitment)
	fmt.Printf("  = SHA256(salt)\n")
	fmt.Printf("Shoe Hash:          %s\n", shoeHash)
	fmt.Printf("  = SHA256(card_order)\n")
	fmt.Printf("Total Cards:        %d\n", len(cards))
	fmt.Println()

	fmt.Println("========== Shuffled Card Order ==========")
	for i := range cards {
		if i > 0 && i%13 == 0 {
			fmt.Println()
		}
		fmt.Printf("#%d:%-4s[%3d] ", i+1, cardNames[i], cardIndices[i])
	}
	fmt.Println()
	fmt.Println()

	fmt.Println("========== First 20 Cards ==========")
	for i := 0; i < 20 && i < len(cards); i++ {
		fmt.Printf("  #%d: %s (Index=%d)\n", i+1, cardNames[i], cardIndices[i])
	}

	return nil
}
