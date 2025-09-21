package main

import (
	"context"
	"fmt"
	"log"

	"github.com/akbariandev/zonnegosdk"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	// Initialize the Zonne SDK client
	client := zonnegosdk.NewClient("https://api.devnet.solana.com", "Aw4Ef9sT3VBv7FXo1qWYR4CQN7LDuTkCcQQC3mxrjwab")

	// Example parameters for minting energy tokens
	params := zonnegosdk.MintRecordCreationParams{
		Grid:          solana.MustPublicKeyFromBase58("11111111111111111111111111111112"), // Example grid pubkey
		Producer:      solana.MustPublicKeyFromBase58("11111111111111111111111111111113"), // Example producer pubkey
		Amount:        1000,                                                               // 1000 kWh
		EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
		GridAuthority: solana.MustPublicKeyFromBase58("11111111111111111111111111111114"), // Example grid authority
	}

	// Payer wallet (this would be the Crossmint wallet address)
	payer := solana.MustPublicKeyFromBase58("11111111111111111111111111111115")

	// Create a base58-encoded transaction for Crossmint
	base58Transaction, err := client.MintEnergyTokensForCrossmint(params, payer)
	if err != nil {
		log.Fatalf("Failed to create transaction for Crossmint: %v", err)
	}

	fmt.Printf("Base58 Transaction for Crossmint: %s\n", base58Transaction)

	// This is how you would use it with Crossmint service
	// (assuming you have a Crossmint client instance)
	/*
		crossmintClient := &crossmint.Crossmint{
			APIKey:    "your-crossmint-api-key",
			IsStaging: true,
		}

		walletLocator := "email:producer@example.com:solana"
		gridAuthorityEmail := "grid@example.com"

		transactionResp, err := crossmintClient.CreateZonneEnergyMintTransaction(
			walletLocator,
			base58Transaction,
			gridAuthorityEmail,
		)
		if err != nil {
			log.Fatalf("Failed to create Crossmint transaction: %v", err)
		}

		fmt.Printf("Crossmint Transaction ID: %s\n", transactionResp.ID)
		fmt.Printf("Transaction Status: %s\n", transactionResp.Status)
	*/

	// You can also create transactions for other Zonne operations
	demonstrateOtherOperations(client, payer)
}

func demonstrateOtherOperations(client *zonnegosdk.Client, payer solana.PublicKey) {
	ctx := context.Background()

	fmt.Println("\n=== Other Zonne Operations for Crossmint ===")

	// 1. Initialize Grid
	gridParams := zonnegosdk.GridAccountCreationParams{
		Grid:      solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
		Authority: solana.MustPublicKeyFromBase58("33333333333333333333333333333333"),
	}

	gridInstruction, err := client.InitializeGrid(gridParams)
	if err != nil {
		log.Printf("Failed to create grid instruction: %v", err)
		return
	}

	// Get latest blockhash for the transaction
	latestBlockhash, err := client.GetRPCClient().GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Printf("Failed to get latest blockhash: %v", err)
		return
	}

	gridTx, err := client.CreateTransactionForCrossmint(gridInstruction, payer, latestBlockhash.Value.Blockhash)
	if err != nil {
		log.Printf("Failed to create grid transaction: %v", err)
		return
	}
	fmt.Printf("Initialize Grid Transaction: %s\n", gridTx)

	// 2. List Tokens for Sale
	listingParams := zonnegosdk.ListingAccountCreationParams{
		Producer:      solana.MustPublicKeyFromBase58("22222222222222222222222222222222"),
		Amount:        500,
		PriceLamports: 1000000, // 0.001 SOL per token
		EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
	}

	listingInstruction, err := client.ListTokensForSale(listingParams)
	if err != nil {
		log.Printf("Failed to create listing instruction: %v", err)
		return
	}

	listingTx, err := client.CreateTransactionForCrossmint(listingInstruction, payer, latestBlockhash.Value.Blockhash)
	if err != nil {
		log.Printf("Failed to create listing transaction: %v", err)
		return
	}
	fmt.Printf("List Tokens Transaction: %s\n", listingTx)

	// 3. Buy Tokens
	buyInstruction, err := client.BuyTokens(
		solana.MustPublicKeyFromBase58("55555555555555555555555555555555"), // buyer
		solana.MustPublicKeyFromBase58("22222222222222222222222222222222"), // producer
		500,     // amount
		1000000, // price
		uint8(zonnegosdk.EnergyTypeSolar),
	)
	if err != nil {
		log.Printf("Failed to create buy instruction: %v", err)
		return
	}

	buyTx, err := client.CreateTransactionForCrossmint(buyInstruction, payer, latestBlockhash.Value.Blockhash)
	if err != nil {
		log.Printf("Failed to create buy transaction: %v", err)
		return
	}
	fmt.Printf("Buy Tokens Transaction: %s\n", buyTx)

	fmt.Println("\n✅ All transactions created successfully for Crossmint integration!")
}
