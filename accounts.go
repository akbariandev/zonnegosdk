package zonnegosdk

import (
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

// PDA derivation functions for all account types

// DeriveGridAccountPDA derives the PDA for a grid account
func (c *Client) DeriveGridAccountPDA(grid solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("grid"),
		grid.Bytes(),
	}
	return solana.FindProgramAddress(seeds, c.programID)
}

// DeriveProducerAccountPDA derives the PDA for a producer account
func (c *Client) DeriveProducerAccountPDA(producer solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("producer"),
		producer.Bytes(),
	}
	return solana.FindProgramAddress(seeds, c.programID)
}

// DeriveConsumerAccountPDA derives the PDA for a consumer account
func (c *Client) DeriveConsumerAccountPDA(consumer solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("consumer"),
		consumer.Bytes(),
	}
	return solana.FindProgramAddress(seeds, c.programID)
}

// DeriveMintRecordPDA derives the PDA for a mint record
func (c *Client) DeriveMintRecordPDA(producer solana.PublicKey, amount uint64, energyType uint8) (solana.PublicKey, uint8, error) {
	amountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBytes, amount)

	seeds := [][]byte{
		[]byte("mint"),
		producer.Bytes(),
		amountBytes,
		{energyType},
	}
	return solana.FindProgramAddress(seeds, c.programID)
}

// DeriveListingAccountPDA derives the PDA for a listing account
func (c *Client) DeriveListingAccountPDA(producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (solana.PublicKey, uint8, error) {
	amountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBytes, amount)

	priceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(priceBytes, priceLamports)

	seeds := [][]byte{
		[]byte("listing"),
		producer.Bytes(),
		amountBytes,
		priceBytes,
		{energyType},
	}
	return solana.FindProgramAddress(seeds, c.programID)
}

// Account size constants for space calculation
const (
	// Account discriminator size (8 bytes)
	AccountDiscriminatorSize = 8

	// GridAccount size: discriminator + is_active (1 byte)
	GridAccountSize = AccountDiscriminatorSize + 1

	// ProducerAccount size: discriminator + balance (8 bytes)
	ProducerAccountSize = AccountDiscriminatorSize + 8

	// ConsumerAccount size: discriminator + consumption (8 bytes)
	ConsumerAccountSize = AccountDiscriminatorSize + 8

	// MintRecord size: discriminator + grid (32) + producer (32) + amount (8) + energy_type (1) + timestamp (8)
	MintRecordSize = AccountDiscriminatorSize + 32 + 32 + 8 + 1 + 8

	// ListingAccount size: discriminator + producer (32) + amount (8) + price_lamports (8) + energy_type (1) + is_active (1) + created_at (8)
	ListingAccountSize = AccountDiscriminatorSize + 32 + 8 + 8 + 1 + 1 + 8
)

// Helper functions for account validation

// IsValidEnergyType checks if the energy type is valid
func IsValidEnergyType(energyType uint8) bool {
	return energyType <= uint8(EnergyTypeOther)
}

// ValidatePublicKey checks if a public key is valid (not zero)
func ValidatePublicKey(pubkey solana.PublicKey) bool {
	return !pubkey.IsZero()
}

// ValidateAmount checks if an amount is valid (greater than zero)
func ValidateAmount(amount uint64) bool {
	return amount > 0
}

// ValidatePrice checks if a price is valid (greater than zero)
func ValidatePrice(priceLamports uint64) bool {
	return priceLamports > 0
}

// Account creation helpers

// GridAccountCreationParams holds parameters for grid account creation
type GridAccountCreationParams struct {
	Grid      solana.PublicKey
	Authority solana.PublicKey
}

// ProducerAccountCreationParams holds parameters for producer account creation
type ProducerAccountCreationParams struct {
	Producer  solana.PublicKey
	Authority solana.PublicKey
}

// ConsumerAccountCreationParams holds parameters for consumer account creation
type ConsumerAccountCreationParams struct {
	Consumer  solana.PublicKey
	Authority solana.PublicKey
}

// MintRecordCreationParams holds parameters for mint record creation
type MintRecordCreationParams struct {
	Grid          solana.PublicKey
	Producer      solana.PublicKey
	Amount        uint64
	EnergyType    uint8
	GridAuthority solana.PublicKey
}

// ListingAccountCreationParams holds parameters for listing account creation
type ListingAccountCreationParams struct {
	Producer      solana.PublicKey
	Amount        uint64
	PriceLamports uint64
	EnergyType    uint8
}
