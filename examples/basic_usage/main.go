package main

import (
	"context"
	"fmt"
	"log"

	"github.com/akbariandev/zonnegosdk"
	"github.com/gagliardetto/solana-go"
)

func main() {
	// Initialize the client with a local RPC endpoint
	client := zonnegosdk.NewClient("http://localhost:8899")

	// Example keypairs (in production, load these securely)
	gridAuthority := solana.MustPrivateKeyFromBase58("your-grid-authority-private-key-here")
	producer := solana.MustPrivateKeyFromBase58("your-producer-private-key-here")
	consumer := solana.MustPrivateKeyFromBase58("your-consumer-private-key-here")

	ctx := context.Background()

	// Example 1: Initialize a grid
	fmt.Println("=== Initializing Grid ===")
	if err := initializeGrid(ctx, client, gridAuthority); err != nil {
		log.Printf("Failed to initialize grid: %v", err)
	}

	// Example 2: Initialize a producer
	fmt.Println("\n=== Initializing Producer ===")
	if err := initializeProducer(ctx, client, producer, gridAuthority); err != nil {
		log.Printf("Failed to initialize producer: %v", err)
	}

	// Example 3: Initialize a consumer
	fmt.Println("\n=== Initializing Consumer ===")
	if err := initializeConsumer(ctx, client, consumer, gridAuthority); err != nil {
		log.Printf("Failed to initialize consumer: %v", err)
	}

	// Example 4: Mint energy tokens
	fmt.Println("\n=== Minting Energy Tokens ===")
	if err := mintEnergyTokens(ctx, client, producer.PublicKey(), gridAuthority); err != nil {
		log.Printf("Failed to mint energy tokens: %v", err)
	}

	// Example 5: List tokens for sale
	fmt.Println("\n=== Listing Tokens for Sale ===")
	if err := listTokensForSale(ctx, client, producer); err != nil {
		log.Printf("Failed to list tokens for sale: %v", err)
	}

	// Example 6: Buy tokens
	fmt.Println("\n=== Buying Tokens ===")
	if err := buyTokens(ctx, client, consumer, producer.PublicKey()); err != nil {
		log.Printf("Failed to buy tokens: %v", err)
	}

	// Example 7: Get account information
	fmt.Println("\n=== Getting Account Information ===")
	if err := getAccountInfo(ctx, client, producer.PublicKey(), consumer.PublicKey()); err != nil {
		log.Printf("Failed to get account info: %v", err)
	}
}

func initializeGrid(ctx context.Context, client *zonnegosdk.Client, gridAuthority solana.PrivateKey) error {
	// Create grid initialization instruction
	params := zonnegosdk.GridAccountCreationParams{
		Grid:      gridAuthority.PublicKey(),
		Authority: gridAuthority.PublicKey(),
	}

	instruction, err := client.InitializeGrid(params)
	if err != nil {
		return fmt.Errorf("failed to create initialize grid instruction: %w", err)
	}

	// Create and send transaction
	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		solana.Hash{}, // Will be set by SendTransaction
		solana.TransactionPayer(gridAuthority.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{gridAuthority})
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Grid initialized successfully! Signature: %s\n", signature)
	return nil
}

func initializeProducer(ctx context.Context, client *zonnegosdk.Client, producer, authority solana.PrivateKey) error {
	params := zonnegosdk.ProducerAccountCreationParams{
		Producer:  producer.PublicKey(),
		Authority: authority.PublicKey(),
	}

	instruction, err := client.InitializeProducer(params)
	if err != nil {
		return fmt.Errorf("failed to create initialize producer instruction: %w", err)
	}

	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		solana.Hash{},
		solana.TransactionPayer(authority.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{authority})
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Producer initialized successfully! Signature: %s\n", signature)
	return nil
}

func initializeConsumer(ctx context.Context, client *zonnegosdk.Client, consumer, authority solana.PrivateKey) error {
	params := zonnegosdk.ConsumerAccountCreationParams{
		Consumer:  consumer.PublicKey(),
		Authority: authority.PublicKey(),
	}

	instruction, err := client.InitializeConsumer(params)
	if err != nil {
		return fmt.Errorf("failed to create initialize consumer instruction: %w", err)
	}

	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		solana.Hash{},
		solana.TransactionPayer(authority.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{authority})
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Consumer initialized successfully! Signature: %s\n", signature)
	return nil
}

func mintEnergyTokens(ctx context.Context, client *zonnegosdk.Client, producer solana.PublicKey, gridAuthority solana.PrivateKey) error {
	params := zonnegosdk.MintRecordCreationParams{
		Grid:          gridAuthority.PublicKey(),
		Producer:      producer,
		Amount:        1000, // 1000 energy tokens
		EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
		GridAuthority: gridAuthority.PublicKey(),
	}

	instruction, err := client.MintEnergyTokens(params)
	if err != nil {
		return fmt.Errorf("failed to create mint energy tokens instruction: %w", err)
	}

	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		solana.Hash{},
		solana.TransactionPayer(gridAuthority.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	transaction.String()

	signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{gridAuthority})
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Energy tokens minted successfully! Signature: %s\n", signature)
	return nil
}

func listTokensForSale(ctx context.Context, client *zonnegosdk.Client, producer solana.PrivateKey) error {
	params := zonnegosdk.ListingAccountCreationParams{
		Producer:      producer.PublicKey(),
		Amount:        500,     // 500 energy tokens
		PriceLamports: 1000000, // 0.001 SOL per token
		EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
	}

	instruction, err := client.ListTokensForSale(params)
	if err != nil {
		return fmt.Errorf("failed to create list tokens instruction: %w", err)
	}

	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		solana.Hash{},
		solana.TransactionPayer(producer.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{producer})
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Tokens listed for sale successfully! Signature: %s\n", signature)
	return nil
}

func buyTokens(ctx context.Context, client *zonnegosdk.Client, buyer solana.PrivateKey, producer solana.PublicKey) error {
	// These parameters should match the listing created above
	amount := uint64(500)
	priceLamports := uint64(1000000)
	energyType := uint8(zonnegosdk.EnergyTypeSolar)

	instruction, err := client.BuyTokens(buyer.PublicKey(), producer, amount, priceLamports, energyType)
	if err != nil {
		return fmt.Errorf("failed to create buy tokens instruction: %w", err)
	}

	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		solana.Hash{},
		solana.TransactionPayer(buyer.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{buyer})
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Tokens purchased successfully! Signature: %s\n", signature)
	return nil
}

func getAccountInfo(ctx context.Context, client *zonnegosdk.Client, producer, consumer solana.PublicKey) error {
	// Get producer account info
	producerAccount, err := client.GetProducerAccount(ctx, producer)
	if err != nil {
		return fmt.Errorf("failed to get producer account: %w", err)
	}
	fmt.Printf("Producer balance: %d energy tokens\n", producerAccount.Balance)

	// Get consumer account info
	consumerAccount, err := client.GetConsumerAccount(ctx, consumer)
	if err != nil {
		return fmt.Errorf("failed to get consumer account: %w", err)
	}
	fmt.Printf("Consumer consumption: %d energy tokens\n", consumerAccount.Consumption)

	return nil
}
