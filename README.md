# Zonne Go SDK

A comprehensive Go SDK for interacting with the Zonne Energy Marketplace on Solana. This SDK provides easy-to-use functions for managing energy producers, consumers, and marketplace transactions.

## Overview

The Zonne Energy Marketplace is a decentralized platform built on Solana that enables:
- **Energy Producers** to mint and sell renewable energy tokens
- **Energy Consumers** to purchase and consume energy tokens
- **Grid Operators** to manage the energy grid and validate transactions

## Features

- 🔌 **Complete Contract Integration** - Full support for all Zonne contract instructions
- 🏗️ **Account Management** - Easy creation and management of grid, producer, and consumer accounts
- 💱 **Marketplace Operations** - List, buy, and cancel energy token listings
- 📊 **Account Queries** - Fetch account states and transaction history
- 🔐 **Type Safety** - Strongly typed Go structs for all contract data
- 🚀 **High Performance** - Efficient RPC communication with connection pooling

## Installation

```bash
go get github.com/zonne/zonnegosdk
```

## Dependencies

The SDK requires the following dependencies:

```bash
go get github.com/gagliardetto/solana-go@v1.10.0
go get github.com/gagliardetto/binary@v0.8.0
go get github.com/near/borsh-go@v0.3.2-0.20220516180422-1ff87d108454
```

## Quick Start

### 1. Initialize the Client

```go
package main

import (
    "github.com/zonne/zonnegosdk"
)

func main() {
    // Connect to local Solana cluster
    client := zonnegosdk.NewClient("http://localhost:8899")
    
    // Or connect to devnet
    // client := zonnegosdk.NewClient("https://api.devnet.solana.com")
}
```

### 2. Setup Accounts

```go
import (
    "context"
    "github.com/gagliardetto/solana-go"
)

func setupAccounts(client *zonnegosdk.Client) error {
    ctx := context.Background()
    
    // Your keypairs (load securely in production)
    gridAuthority := solana.MustPrivateKeyFromBase58("your-private-key")
    producer := solana.MustPrivateKeyFromBase58("producer-private-key")
    consumer := solana.MustPrivateKeyFromBase58("consumer-private-key")
    
    // Initialize grid
    gridParams := zonnegosdk.GridAccountCreationParams{
        Grid:      gridAuthority.PublicKey(),
        Authority: gridAuthority.PublicKey(),
    }
    gridInstruction, err := client.InitializeGrid(gridParams)
    if err != nil {
        return err
    }
    
    // Initialize producer
    producerParams := zonnegosdk.ProducerAccountCreationParams{
        Producer:  producer.PublicKey(),
        Authority: gridAuthority.PublicKey(),
    }
    producerInstruction, err := client.InitializeProducer(producerParams)
    if err != nil {
        return err
    }
    
    // Create and send transactions...
    return nil
}
```

### 3. Mint Energy Tokens

```go
func mintEnergy(client *zonnegosdk.Client, producer solana.PublicKey, gridAuth solana.PrivateKey) error {
    ctx := context.Background()
    
    params := zonnegosdk.MintRecordCreationParams{
        Grid:          gridAuth.PublicKey(),
        Producer:      producer,
        Amount:        1000, // 1000 kWh
        EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
        GridAuthority: gridAuth.PublicKey(),
    }
    
    instruction, err := client.MintEnergyTokens(params)
    if err != nil {
        return err
    }
    
    // Create and send transaction
    transaction, err := solana.NewTransaction(
        []solana.Instruction{*instruction},
        solana.Hash{},
        solana.TransactionPayer(gridAuth.PublicKey()),
    )
    if err != nil {
        return err
    }
    
    signature, err := client.SendAndConfirmTransaction(ctx, transaction, []solana.PrivateKey{gridAuth})
    if err != nil {
        return err
    }
    
    fmt.Printf("Energy minted! Signature: %s\n", signature)
    return nil
}
```

### 4. Create Energy Listings

```go
func createListing(client *zonnegosdk.Client, producer solana.PrivateKey) error {
    ctx := context.Background()
    
    params := zonnegosdk.ListingAccountCreationParams{
        Producer:      producer.PublicKey(),
        Amount:        500, // 500 kWh
        PriceLamports: 1000000, // 0.001 SOL per kWh
        EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
    }
    
    instruction, err := client.ListTokensForSale(params)
    if err != nil {
        return err
    }
    
    // Create and send transaction...
    return nil
}
```

### 5. Purchase Energy

```go
func buyEnergy(client *zonnegosdk.Client, buyer solana.PrivateKey, producer solana.PublicKey) error {
    ctx := context.Background()
    
    // Parameters must match an existing listing
    instruction, err := client.BuyTokens(
        buyer.PublicKey(),
        producer,
        500,        // amount
        1000000,    // price in lamports
        uint8(zonnegosdk.EnergyTypeSolar),
    )
    if err != nil {
        return err
    }
    
    // Create and send transaction...
    return nil
}
```

## API Reference

### Client

#### `NewClient(rpcEndpoint string) *Client`
Creates a new Zonne SDK client connected to the specified RPC endpoint.

#### `NewClientWithCustomProgram(rpcEndpoint string, programID solana.PublicKey) *Client`
Creates a client with a custom program ID (useful for testing).

### Account Management

#### Grid Operations
- `InitializeGrid(params GridAccountCreationParams) (*solana.Instruction, error)`
- `GetGridAccount(ctx context.Context, gridPubkey solana.PublicKey) (*GridAccount, error)`

#### Producer Operations
- `InitializeProducer(params ProducerAccountCreationParams) (*solana.Instruction, error)`
- `GetProducerAccount(ctx context.Context, producerPubkey solana.PublicKey) (*ProducerAccount, error)`

#### Consumer Operations
- `InitializeConsumer(params ConsumerAccountCreationParams) (*solana.Instruction, error)`
- `GetConsumerAccount(ctx context.Context, consumerPubkey solana.PublicKey) (*ConsumerAccount, error)`

### Energy Token Operations

#### Minting
- `MintEnergyTokens(params MintRecordCreationParams) (*solana.Instruction, error)`
- `MintConsumptionTokens(consumer, grid, gridAuthority solana.PublicKey, amount uint64) (*solana.Instruction, error)`

#### Marketplace
- `ListTokensForSale(params ListingAccountCreationParams) (*solana.Instruction, error)`
- `BuyTokens(buyer, producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (*solana.Instruction, error)`
- `CancelListing(producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (*solana.Instruction, error)`

### Account Queries
- `GetListingAccount(ctx context.Context, producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (*ListingAccount, error)`
- `GetMintRecord(ctx context.Context, producer solana.PublicKey, amount uint64, energyType uint8) (*MintRecord, error)`

### PDA Derivation
- `DeriveGridAccountPDA(grid solana.PublicKey) (solana.PublicKey, uint8, error)`
- `DeriveProducerAccountPDA(producer solana.PublicKey) (solana.PublicKey, uint8, error)`
- `DeriveConsumerAccountPDA(consumer solana.PublicKey) (solana.PublicKey, uint8, error)`
- `DeriveMintRecordPDA(producer solana.PublicKey, amount uint64, energyType uint8) (solana.PublicKey, uint8, error)`
- `DeriveListingAccountPDA(producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (solana.PublicKey, uint8, error)`

### Crossmint Integration
- `MintEnergyTokensForCrossmint(params MintRecordCreationParams, payer solana.PublicKey) (string, error)`
- `CreateTransactionForCrossmint(instruction solana.Instruction, payer solana.PublicKey, latestBlockhash solana.Hash) (string, error)`

## Data Types

### Energy Types
```go
type EnergyType uint8

const (
    EnergyTypeSolar EnergyType = iota
    EnergyTypeWind
    EnergyTypeHydro
    EnergyTypeOther
)
```

### Account Structures
```go
type GridAccount struct {
    IsActive bool `borsh:"is_active"`
}

type ProducerAccount struct {
    Balance uint64 `borsh:"balance"`
}

type ConsumerAccount struct {
    Consumption uint64 `borsh:"consumption"`
}

type ListingAccount struct {
    Producer       solana.PublicKey `borsh:"producer"`
    Amount         uint64           `borsh:"amount"`
    PriceLamports  uint64           `borsh:"price_lamports"`
    EnergyType     uint8            `borsh:"energy_type"`
    IsActive       bool             `borsh:"is_active"`
    CreatedAt      int64            `borsh:"created_at"`
}

type MintRecord struct {
    Grid        solana.PublicKey `borsh:"grid"`
    Producer    solana.PublicKey `borsh:"producer"`
    Amount      uint64           `borsh:"amount"`
    EnergyType  uint8            `borsh:"energy_type"`
    Timestamp   int64            `borsh:"timestamp"`
}
```

## Examples

### Complete Marketplace Demo

See `examples/marketplace_demo/main.go` for a comprehensive example that demonstrates:
1. Setting up the marketplace
2. Initializing producers and consumers
3. Minting energy tokens
4. Creating and managing listings
5. Purchasing energy tokens
6. Querying account states

### Basic Usage

See `examples/basic_usage/main.go` for simple examples of each SDK function.

### Crossmint Integration

See `examples/crossmint_integration/main.go` for examples of creating base58-encoded transactions for Crossmint:
1. Creating transactions for Crossmint's API
2. Minting energy tokens via Crossmint
3. Using Crossmint smart wallets for transactions

## Error Handling

The SDK provides detailed error messages for common issues:

```go
// Validate inputs before creating instructions
if !zonnegosdk.ValidatePublicKey(producer) {
    return fmt.Errorf("invalid producer public key")
}

if !zonnegosdk.ValidateAmount(amount) {
    return fmt.Errorf("amount must be greater than zero")
}

if !zonnegosdk.IsValidEnergyType(energyType) {
    return fmt.Errorf("invalid energy type")
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

For integration tests with a local Solana cluster:

```bash
# Start local Solana test validator
solana-test-validator

# Run tests
go test -tags=integration ./...
```

## Crossmint Integration

The SDK provides built-in support for [Crossmint](https://docs.crossmint.com/api-reference/wallets/create-transaction) smart wallets and transaction creation.

### Creating Transactions for Crossmint

```go
// Create a base58-encoded transaction for Crossmint
params := zonnegosdk.MintRecordCreationParams{
    Grid:          gridPubkey,
    Producer:      producerPubkey,
    Amount:        1000,
    EnergyType:    uint8(zonnegosdk.EnergyTypeSolar),
    GridAuthority: gridAuthorityPubkey,
}

payer := crossmintWalletAddress
base58Transaction, err := client.MintEnergyTokensForCrossmint(params, payer)
if err != nil {
    log.Fatal(err)
}

// Use with Crossmint API
fmt.Printf("Base58 Transaction: %s\n", base58Transaction)
```

### Generic Transaction Creation

```go
// Create any instruction and convert to base58 for Crossmint
instruction, err := client.InitializeGrid(gridParams)
if err != nil {
    log.Fatal(err)
}

// Get latest blockhash
latestBlockhash, err := client.GetRPCClient().GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
if err != nil {
    log.Fatal(err)
}

// Create base58 transaction
base58Tx, err := client.CreateTransactionForCrossmint(
    instruction, 
    payer, 
    latestBlockhash.Value.Blockhash,
)
if err != nil {
    log.Fatal(err)
}
```

### Crossmint Service Integration

The SDK works with Crossmint's transaction creation API:

```go
// This would be in your Crossmint service
transactionResp, err := crossmintClient.CreateTransaction(
    "email:user@example.com:solana",  // wallet locator
    base58Transaction,                 // from SDK
    []string{"email:signer@example.com"}, // signers
)
```

## Configuration

### Environment Variables

- `SOLANA_RPC_URL` - Default RPC endpoint
- `ZONNE_PROGRAM_ID` - Custom program ID for testing

### Connection Settings

```go
// Configure custom RPC settings
client := zonnegosdk.NewClient("https://api.mainnet-beta.solana.com")

// Set custom commitment level
rpcClient := client.GetRPCClient()
// Use rpcClient with custom settings...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📧 Email: support@zonne.energy
- 💬 Discord: [Zonne Community](https://discord.gg/zonne)
- 🐛 Issues: [GitHub Issues](https://github.com/zonne/zonnegosdk/issues)

## Changelog

### v1.0.0
- Initial release
- Full contract integration
- Complete SDK functionality
- Examples and documentation
