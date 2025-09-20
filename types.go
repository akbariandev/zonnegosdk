package zonnegosdk

import (
	"time"

	"github.com/gagliardetto/solana-go"
)

// EnergyType represents different types of energy
type EnergyType uint8

const (
	EnergyTypeSolar EnergyType = iota
	EnergyTypeWind
	EnergyTypeHydro
	EnergyTypeOther
)

// GridAccount represents the grid state account
type GridAccount struct {
	IsActive bool `borsh:"is_active"`
}

// ProducerAccount represents a producer's energy balance
type ProducerAccount struct {
	Balance uint64 `borsh:"balance"`
}

// ConsumerAccount represents a consumer's energy consumption
type ConsumerAccount struct {
	Consumption uint64 `borsh:"consumption"`
}

// MintRecord represents a record of energy token minting
type MintRecord struct {
	Grid        solana.PublicKey `borsh:"grid"`
	Producer    solana.PublicKey `borsh:"producer"`
	Amount      uint64           `borsh:"amount"`
	EnergyType  uint8            `borsh:"energy_type"`
	Timestamp   int64            `borsh:"timestamp"`
}

// ListingAccount represents an active energy token listing for sale
type ListingAccount struct {
	Producer       solana.PublicKey `borsh:"producer"`
	Amount         uint64           `borsh:"amount"`
	PriceLamports  uint64           `borsh:"price_lamports"`
	EnergyType     uint8            `borsh:"energy_type"`
	IsActive       bool             `borsh:"is_active"`
	CreatedAt      int64            `borsh:"created_at"`
}

// Event types for contract events

// GridInitializedEvent represents a grid initialization event
type GridInitializedEvent struct {
	Grid solana.PublicKey `json:"grid"`
}

// ProducerInitializedEvent represents a producer initialization event
type ProducerInitializedEvent struct {
	Producer solana.PublicKey `json:"producer"`
}

// ConsumerInitializedEvent represents a consumer initialization event
type ConsumerInitializedEvent struct {
	Consumer solana.PublicKey `json:"consumer"`
}

// TokensMintedEvent represents an energy token minting event
type TokensMintedEvent struct {
	Producer   solana.PublicKey `json:"producer"`
	Amount     uint64           `json:"amount"`
	EnergyType uint8            `json:"energy_type"`
}

// TokensListedEvent represents an energy token listing event
type TokensListedEvent struct {
	ListingID     solana.PublicKey `json:"listing_id"`
	Producer      solana.PublicKey `json:"producer"`
	Amount        uint64           `json:"amount"`
	PriceLamports uint64           `json:"price_lamports"`
	EnergyType    uint8            `json:"energy_type"`
}

// ListingCancelledEvent represents a listing cancellation event
type ListingCancelledEvent struct {
	ListingID solana.PublicKey `json:"listing_id"`
	Producer  solana.PublicKey `json:"producer"`
	Amount    uint64           `json:"amount"`
}

// TokensPurchasedEvent represents a token purchase event
type TokensPurchasedEvent struct {
	ListingID     solana.PublicKey `json:"listing_id"`
	Buyer         solana.PublicKey `json:"buyer"`
	Producer      solana.PublicKey `json:"producer"`
	Amount        uint64           `json:"amount"`
	PriceLamports uint64           `json:"price_lamports"`
}

// ConsumptionMintedEvent represents a consumption token minting event
type ConsumptionMintedEvent struct {
	Consumer solana.PublicKey `json:"consumer"`
	Amount   uint64           `json:"amount"`
}

// Helper methods for time conversion
func (m *MintRecord) GetTimestamp() time.Time {
	return time.Unix(m.Timestamp, 0)
}

func (l *ListingAccount) GetCreatedAt() time.Time {
	return time.Unix(l.CreatedAt, 0)
}

// Helper methods for energy type
func (e EnergyType) String() string {
	switch e {
	case EnergyTypeSolar:
		return "Solar"
	case EnergyTypeWind:
		return "Wind"
	case EnergyTypeHydro:
		return "Hydro"
	default:
		return "Other"
	}
}

// ParseEnergyType converts a string to EnergyType
func ParseEnergyType(s string) EnergyType {
	switch s {
	case "Solar", "solar":
		return EnergyTypeSolar
	case "Wind", "wind":
		return EnergyTypeWind
	case "Hydro", "hydro":
		return EnergyTypeHydro
	default:
		return EnergyTypeOther
	}
}
