package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/akbariandev/zonnegosdk"
	"github.com/gagliardetto/solana-go"
)

// MarketplaceDemo demonstrates a complete energy marketplace workflow
func main() {
	// Initialize client
	client := zonnegosdk.NewClient("http://localhost:8899")
	ctx := context.Background()

	// Demo keypairs (replace with actual keypairs in production)
	gridAuthority := solana.MustPrivateKeyFromBase58("your-grid-authority-private-key-here")
	producer1 := solana.MustPrivateKeyFromBase58("your-producer1-private-key-here")
	producer2 := solana.MustPrivateKeyFromBase58("your-producer2-private-key-here")
	consumer1 := solana.MustPrivateKeyFromBase58("your-consumer1-private-key-here")
	consumer2 := solana.MustPrivateKeyFromBase58("your-consumer2-private-key-here")

	fmt.Println("🌞 Zonne Energy Marketplace Demo")
	fmt.Println("==================================")

	// Step 1: Setup the marketplace
	if err := setupMarketplace(ctx, client, gridAuthority, producer1, producer2, consumer1, consumer2); err != nil {
		log.Fatalf("Failed to setup marketplace: %v", err)
	}

	// Step 2: Producers generate and mint energy
	if err := generateEnergy(ctx, client, gridAuthority, producer1, producer2); err != nil {
		log.Fatalf("Failed to generate energy: %v", err)
	}

	// Step 3: Create energy listings
	if err := createListings(ctx, client, producer1, producer2); err != nil {
		log.Fatalf("Failed to create listings: %v", err)
	}

	// Step 4: Consumers purchase energy
	if err := purchaseEnergy(ctx, client, consumer1, consumer2, producer1, producer2); err != nil {
		log.Fatalf("Failed to purchase energy: %v", err)
	}

	// Step 5: Display final marketplace state
	if err := displayMarketplaceState(ctx, client, producer1, producer2, consumer1, consumer2); err != nil {
		log.Fatalf("Failed to display marketplace state: %v", err)
	}

	fmt.Println("\n✅ Marketplace demo completed successfully!")
}

func setupMarketplace(ctx context.Context, client *zonnegosdk.Client, gridAuth, prod1, prod2, cons1, cons2 solana.PrivateKey) error {
	fmt.Println("\n🏗️  Setting up marketplace...")

	// Initialize grid
	gridParams := zonnegosdk.GridAccountCreationParams{
		Grid:      gridAuth.PublicKey(),
		Authority: gridAuth.PublicKey(),
	}
	if err := executeTransaction(ctx, client, "Initialize Grid", func() (solana.Instruction, error) {
		return client.InitializeGrid(gridParams)
	}, []solana.PrivateKey{gridAuth}); err != nil {
		return err
	}

	// Initialize producers
	for i, producer := range []solana.PrivateKey{prod1, prod2} {
		params := zonnegosdk.ProducerAccountCreationParams{
			Producer:  producer.PublicKey(),
			Authority: gridAuth.PublicKey(),
		}
		if err := executeTransaction(ctx, client, fmt.Sprintf("Initialize Producer %d", i+1), func() (solana.Instruction, error) {
			return client.InitializeProducer(params)
		}, []solana.PrivateKey{gridAuth}); err != nil {
			return err
		}
	}

	// Initialize consumers
	for i, consumer := range []solana.PrivateKey{cons1, cons2} {
		params := zonnegosdk.ConsumerAccountCreationParams{
			Consumer:  consumer.PublicKey(),
			Authority: gridAuth.PublicKey(),
		}
		if err := executeTransaction(ctx, client, fmt.Sprintf("Initialize Consumer %d", i+1), func() (solana.Instruction, error) {
			return client.InitializeConsumer(params)
		}, []solana.PrivateKey{gridAuth}); err != nil {
			return err
		}
	}

	fmt.Println("   ✅ Marketplace setup complete!")
	return nil
}

func generateEnergy(ctx context.Context, client *zonnegosdk.Client, gridAuth, prod1, prod2 solana.PrivateKey) error {
	fmt.Println("\n⚡ Generating energy...")

	// Producer 1: Generate solar energy
	solarParams := zonnegosdk.MintRecordCreationParams{
		Grid:          gridAuth.PublicKey(),
		Producer:      prod1.PublicKey(),
		Amount:        2000, // 2000 kWh
		EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
		GridAuthority: gridAuth.PublicKey(),
	}
	if err := executeTransaction(ctx, client, "Producer 1: Generate Solar Energy", func() (solana.Instruction, error) {
		return client.MintEnergyTokens(solarParams)
	}, []solana.PrivateKey{gridAuth}); err != nil {
		return err
	}

	// Producer 2: Generate wind energy
	windParams := zonnegosdk.MintRecordCreationParams{
		Grid:          gridAuth.PublicKey(),
		Producer:      prod2.PublicKey(),
		Amount:        1500, // 1500 kWh
		EnergyType:    uint8(zonnegosdk.EnergyTypeWind),
		GridAuthority: gridAuth.PublicKey(),
	}
	if err := executeTransaction(ctx, client, "Producer 2: Generate Wind Energy", func() (solana.Instruction, error) {
		return client.MintEnergyTokens(windParams)
	}, []solana.PrivateKey{gridAuth}); err != nil {
		return err
	}

	fmt.Println("   ✅ Energy generation complete!")
	return nil
}

func createListings(ctx context.Context, client *zonnegosdk.Client, prod1, prod2 solana.PrivateKey) error {
	fmt.Println("\n📝 Creating energy listings...")

	// Producer 1: List solar energy
	solarListing := zonnegosdk.ListingAccountCreationParams{
		Producer:      prod1.PublicKey(),
		Amount:        1000,   // 1000 kWh
		PriceLamports: 500000, // 0.0005 SOL per kWh
		EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
	}
	if err := executeTransaction(ctx, client, "Producer 1: List Solar Energy", func() (solana.Instruction, error) {
		return client.ListTokensForSale(solarListing)
	}, []solana.PrivateKey{prod1}); err != nil {
		return err
	}

	// Producer 2: List wind energy
	windListing := zonnegosdk.ListingAccountCreationParams{
		Producer:      prod2.PublicKey(),
		Amount:        800,    // 800 kWh
		PriceLamports: 400000, // 0.0004 SOL per kWh (cheaper than solar)
		EnergyType:    uint8(zonnegosdk.EnergyTypeWind),
	}
	if err := executeTransaction(ctx, client, "Producer 2: List Wind Energy", func() (solana.Instruction, error) {
		return client.ListTokensForSale(windListing)
	}, []solana.PrivateKey{prod2}); err != nil {
		return err
	}

	fmt.Println("   ✅ Energy listings created!")
	return nil
}

func purchaseEnergy(ctx context.Context, client *zonnegosdk.Client, cons1, cons2, prod1, prod2 solana.PrivateKey) error {
	fmt.Println("\n💰 Purchasing energy...")

	// Consumer 1: Buy solar energy from Producer 1
	if err := executeTransaction(ctx, client, "Consumer 1: Buy Solar Energy", func() (solana.Instruction, error) {
		return client.BuyTokens(cons1.PublicKey(), prod1.PublicKey(), 1000, 500000, uint8(zonnegosdk.EnergyTypeSolar))
	}, []solana.PrivateKey{cons1}); err != nil {
		return err
	}

	// Consumer 2: Buy wind energy from Producer 2
	if err := executeTransaction(ctx, client, "Consumer 2: Buy Wind Energy", func() (solana.Instruction, error) {
		return client.BuyTokens(cons2.PublicKey(), prod2.PublicKey(), 800, 400000, uint8(zonnegosdk.EnergyTypeWind))
	}, []solana.PrivateKey{cons2}); err != nil {
		return err
	}

	fmt.Println("   ✅ Energy purchases complete!")
	return nil
}

func displayMarketplaceState(ctx context.Context, client *zonnegosdk.Client, prod1, prod2, cons1, cons2 solana.PrivateKey) error {
	fmt.Println("\n📊 Final Marketplace State")
	fmt.Println("==========================")

	// Display producer balances
	fmt.Println("\n🏭 Producers:")
	for i, producer := range []solana.PrivateKey{prod1, prod2} {
		account, err := client.GetProducerAccount(ctx, producer.PublicKey())
		if err != nil {
			return fmt.Errorf("failed to get producer %d account: %w", i+1, err)
		}
		energyType := []string{"Solar", "Wind"}[i]
		fmt.Printf("   Producer %d (%s): %d kWh remaining\n", i+1, energyType, account.Balance)
	}

	// Display consumer consumption
	fmt.Println("\n🏠 Consumers:")
	for i, consumer := range []solana.PrivateKey{cons1, cons2} {
		account, err := client.GetConsumerAccount(ctx, consumer.PublicKey())
		if err != nil {
			return fmt.Errorf("failed to get consumer %d account: %w", i+1, err)
		}
		fmt.Printf("   Consumer %d: %d kWh consumed\n", i+1, account.Consumption)
	}

	return nil
}

func executeTransaction(ctx context.Context, client *zonnegosdk.Client, description string, createInstruction func() (solana.Instruction, error), signers []solana.PrivateKey) error {
	fmt.Printf("   🔄 %s...", description)

	instruction, err := createInstruction()
	if err != nil {
		fmt.Printf(" ❌\n")
		return fmt.Errorf("failed to create instruction for %s: %w", description, err)
	}

	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		solana.Hash{},
		solana.TransactionPayer(signers[0].PublicKey()),
	)
	if err != nil {
		fmt.Printf(" ❌\n")
		return fmt.Errorf("failed to create transaction for %s: %w", description, err)
	}

	signature, err := client.SendAndConfirmTransaction(ctx, transaction, signers)
	if err != nil {
		fmt.Printf(" ❌\n")
		return fmt.Errorf("failed to send transaction for %s: %w", description, err)
	}

	fmt.Printf(" ✅ (Sig: %s)\n", signature.String()[:8]+"...")

	// Small delay to prevent rate limiting
	time.Sleep(100 * time.Millisecond)

	return nil
}
