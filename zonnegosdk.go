// Package zonnegosdk provides a comprehensive Go SDK for interacting with the Zonne Energy Marketplace on Solana.
//
// The Zonne Energy Marketplace is a decentralized platform that enables energy producers to mint and sell
// renewable energy tokens, while consumers can purchase and consume these tokens. Grid operators manage
// the energy grid and validate transactions.
//
// Key features:
//   - Complete contract integration with all Zonne contract instructions
//   - Easy account management for grids, producers, and consumers
//   - Marketplace operations for listing, buying, and canceling energy tokens
//   - Account queries for fetching states and transaction history
//   - Type-safe Go structs for all contract data
//   - High-performance RPC communication
//
// Basic usage:
//
//	// Initialize client
//	client := zonnegosdk.NewClient("http://localhost:8899")
//
//	// Setup accounts
//	gridParams := zonnegosdk.GridAccountCreationParams{
//		Grid:      gridAuthority.PublicKey(),
//		Authority: gridAuthority.PublicKey(),
//	}
//	instruction, err := client.InitializeGrid(gridParams)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create and send transaction
//	transaction, err := solana.NewTransaction(
//		[]solana.Instruction{*instruction},
//		solana.Hash{},
//		solana.TransactionPayer(gridAuthority.PublicKey()),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{gridAuthority})
//	if err != nil {
//		log.Fatal(err)
//	}
//
// For more examples, see the examples/ directory.
package zonnegosdk

import "github.com/gagliardetto/solana-go"

// Version of the SDK
const Version = "1.0.0"

// Default RPC endpoints for different Solana clusters
const (
	LocalnetRPC = "http://localhost:8899"
	DevnetRPC   = "https://api.devnet.solana.com"
	TestnetRPC  = "https://api.testnet.solana.com"
	MainnetRPC  = "https://api.mainnet-beta.solana.com"
)

// Common Solana program IDs that might be useful
var (
	SystemProgramID          = solana.SystemProgramID
	TokenProgramID           = solana.TokenProgramID
	AssociatedTokenProgramID = solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
	RentSysvarID            = solana.MustPublicKeyFromBase58("SysvarRent111111111111111111111111111111111")
	ClockSysvarID           = solana.MustPublicKeyFromBase58("SysvarC1ock11111111111111111111111111111111")
)

// Energy marketplace constants
const (
	// Minimum amounts for various operations
	MinEnergyAmount = uint64(1)
	MinPriceLamports = uint64(1)
	
	// Maximum values to prevent overflow
	MaxEnergyAmount = uint64(1<<63 - 1) // Max int64
	MaxPriceLamports = uint64(1<<63 - 1)
)

// Common error messages
const (
	ErrInvalidPublicKey = "invalid public key: cannot be zero"
	ErrInvalidAmount    = "invalid amount: must be greater than zero"
	ErrInvalidPrice     = "invalid price: must be greater than zero"
	ErrInvalidEnergyType = "invalid energy type: must be 0-3"
)

// Utility functions

// IsZeroPublicKey checks if a public key is the zero/empty key
func IsZeroPublicKey(key solana.PublicKey) bool {
	return key.IsZero()
}

// LamportsToSOL converts lamports to SOL
func LamportsToSOL(lamports uint64) float64 {
	return float64(lamports) / 1e9
}

// SOLToLamports converts SOL to lamports
func SOLToLamports(sol float64) uint64 {
	return uint64(sol * 1e9)
}
