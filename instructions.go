package zonnegosdk

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
)

// Instruction discriminators (first 8 bytes of each instruction)
var (
	InitializeGridDiscriminator        = [8]byte{30, 224, 57, 59, 226, 1, 253, 219}
	InitializeProducerDiscriminator    = [8]byte{168, 203, 156, 75, 245, 233, 21, 201}
	InitializeConsumerDiscriminator    = [8]byte{228, 176, 96, 150, 209, 224, 65, 98}
	MintEnergyTokensDiscriminator      = [8]byte{147, 199, 3, 69, 8, 89, 72, 226}
	ListTokensForSaleDiscriminator     = [8]byte{213, 6, 33, 225, 91, 199, 59, 195}
	CancelListingDiscriminator         = [8]byte{41, 183, 50, 232, 230, 233, 157, 70}
	BuyTokensDiscriminator             = [8]byte{189, 21, 230, 133, 247, 2, 110, 42}
	MintConsumptionTokensDiscriminator = [8]byte{75, 241, 244, 71, 205, 59, 169, 126}
)

// InitializeGrid creates an instruction to initialize a grid account
func (c *Client) InitializeGrid(params GridAccountCreationParams) (solana.Instruction, error) {
	if !ValidatePublicKey(params.Grid) {
		return nil, fmt.Errorf("invalid grid public key")
	}
	if !ValidatePublicKey(params.Authority) {
		return nil, fmt.Errorf("invalid authority public key")
	}

	gridAccountPDA, _, err := c.DeriveGridAccountPDA(params.Grid)
	if err != nil {
		return nil, fmt.Errorf("failed to derive grid account PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: gridAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: params.Grid, IsWritable: false, IsSigner: false},
		{PublicKey: params.Authority, IsWritable: true, IsSigner: true},
		{PublicKey: solana.SystemProgramID, IsWritable: false, IsSigner: false},
	}

	data := InitializeGridDiscriminator[:]

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// InitializeProducer creates an instruction to initialize a producer account
func (c *Client) InitializeProducer(params ProducerAccountCreationParams) (solana.Instruction, error) {
	if !ValidatePublicKey(params.Producer) {
		return nil, fmt.Errorf("invalid producer public key")
	}
	if !ValidatePublicKey(params.Authority) {
		return nil, fmt.Errorf("invalid authority public key")
	}

	producerAccountPDA, _, err := c.DeriveProducerAccountPDA(params.Producer)
	if err != nil {
		return nil, fmt.Errorf("failed to derive producer account PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: producerAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: params.Producer, IsWritable: false, IsSigner: false},
		{PublicKey: params.Authority, IsWritable: true, IsSigner: true},
		{PublicKey: solana.SystemProgramID, IsWritable: false, IsSigner: false},
	}

	data := InitializeProducerDiscriminator[:]

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// InitializeConsumer creates an instruction to initialize a consumer account
func (c *Client) InitializeConsumer(params ConsumerAccountCreationParams) (solana.Instruction, error) {
	if !ValidatePublicKey(params.Consumer) {
		return nil, fmt.Errorf("invalid consumer public key")
	}
	if !ValidatePublicKey(params.Authority) {
		return nil, fmt.Errorf("invalid authority public key")
	}

	consumerAccountPDA, _, err := c.DeriveConsumerAccountPDA(params.Consumer)
	if err != nil {
		return nil, fmt.Errorf("failed to derive consumer account PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: consumerAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: params.Consumer, IsWritable: false, IsSigner: false},
		{PublicKey: params.Authority, IsWritable: true, IsSigner: true},
		{PublicKey: solana.SystemProgramID, IsWritable: false, IsSigner: false},
	}

	data := InitializeConsumerDiscriminator[:]

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// MintEnergyTokens creates an instruction to mint energy tokens
func (c *Client) MintEnergyTokens(params MintRecordCreationParams) (solana.Instruction, error) {
	if !ValidatePublicKey(params.Grid) {
		return nil, fmt.Errorf("invalid grid public key")
	}
	if !ValidatePublicKey(params.Producer) {
		return nil, fmt.Errorf("invalid producer public key")
	}
	if !ValidatePublicKey(params.GridAuthority) {
		return nil, fmt.Errorf("invalid grid authority public key")
	}
	if !ValidateAmount(params.Amount) {
		return nil, fmt.Errorf("invalid amount")
	}
	if !IsValidEnergyType(params.EnergyType) {
		return nil, fmt.Errorf("invalid energy type")
	}

	producerAccountPDA, _, err := c.DeriveProducerAccountPDA(params.Producer)
	if err != nil {
		return nil, fmt.Errorf("failed to derive producer account PDA: %w", err)
	}

	gridAccountPDA, _, err := c.DeriveGridAccountPDA(params.Grid)
	if err != nil {
		return nil, fmt.Errorf("failed to derive grid account PDA: %w", err)
	}

	mintRecordPDA, _, err := c.DeriveMintRecordPDA(params.Producer, params.Amount, params.EnergyType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive mint record PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: producerAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: gridAccountPDA, IsWritable: false, IsSigner: false},
		{PublicKey: mintRecordPDA, IsWritable: true, IsSigner: false},
		{PublicKey: params.Producer, IsWritable: false, IsSigner: false},
		{PublicKey: params.GridAuthority, IsWritable: true, IsSigner: true},
		{PublicKey: solana.SystemProgramID, IsWritable: false, IsSigner: false},
	}

	// Serialize instruction data
	type InstructionData struct {
		Amount     uint64 `borsh:"amount"`
		EnergyType uint8  `borsh:"energy_type"`
	}

	instructionData := InstructionData{
		Amount:     params.Amount,
		EnergyType: params.EnergyType,
	}

	serializedData, err := borsh.Serialize(instructionData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize instruction data: %w", err)
	}

	data := append(MintEnergyTokensDiscriminator[:], serializedData...)

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// ListTokensForSale creates an instruction to list tokens for sale
func (c *Client) ListTokensForSale(params ListingAccountCreationParams) (solana.Instruction, error) {
	if !ValidatePublicKey(params.Producer) {
		return nil, fmt.Errorf("invalid producer public key")
	}
	if !ValidateAmount(params.Amount) {
		return nil, fmt.Errorf("invalid amount")
	}
	if !ValidatePrice(params.PriceLamports) {
		return nil, fmt.Errorf("invalid price")
	}
	if !IsValidEnergyType(params.EnergyType) {
		return nil, fmt.Errorf("invalid energy type")
	}

	producerAccountPDA, _, err := c.DeriveProducerAccountPDA(params.Producer)
	if err != nil {
		return nil, fmt.Errorf("failed to derive producer account PDA: %w", err)
	}

	listingAccountPDA, _, err := c.DeriveListingAccountPDA(params.Producer, params.Amount, params.PriceLamports, params.EnergyType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive listing account PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: producerAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: listingAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: params.Producer, IsWritable: true, IsSigner: true},
		{PublicKey: solana.SystemProgramID, IsWritable: false, IsSigner: false},
	}

	// Serialize instruction data
	type InstructionData struct {
		Amount        uint64 `borsh:"amount"`
		PriceLamports uint64 `borsh:"price_lamports"`
		EnergyType    uint8  `borsh:"energy_type"`
	}

	instructionData := InstructionData{
		Amount:        params.Amount,
		PriceLamports: params.PriceLamports,
		EnergyType:    params.EnergyType,
	}

	serializedData, err := borsh.Serialize(instructionData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize instruction data: %w", err)
	}

	data := append(ListTokensForSaleDiscriminator[:], serializedData...)

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// CancelListing creates an instruction to cancel a listing
func (c *Client) CancelListing(producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (solana.Instruction, error) {
	if !ValidatePublicKey(producer) {
		return nil, fmt.Errorf("invalid producer public key")
	}

	listingAccountPDA, _, err := c.DeriveListingAccountPDA(producer, amount, priceLamports, energyType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive listing account PDA: %w", err)
	}

	producerAccountPDA, _, err := c.DeriveProducerAccountPDA(producer)
	if err != nil {
		return nil, fmt.Errorf("failed to derive producer account PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: listingAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: producerAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: producer, IsWritable: true, IsSigner: true},
	}

	data := CancelListingDiscriminator[:]

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// BuyTokens creates an instruction to buy tokens from a listing
func (c *Client) BuyTokens(buyer, producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (solana.Instruction, error) {
	if !ValidatePublicKey(buyer) {
		return nil, fmt.Errorf("invalid buyer public key")
	}
	if !ValidatePublicKey(producer) {
		return nil, fmt.Errorf("invalid producer public key")
	}

	listingAccountPDA, _, err := c.DeriveListingAccountPDA(producer, amount, priceLamports, energyType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive listing account PDA: %w", err)
	}

	consumerAccountPDA, _, err := c.DeriveConsumerAccountPDA(buyer)
	if err != nil {
		return nil, fmt.Errorf("failed to derive consumer account PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: listingAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: consumerAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: producer, IsWritable: true, IsSigner: false},
		{PublicKey: buyer, IsWritable: true, IsSigner: true},
		{PublicKey: solana.SystemProgramID, IsWritable: false, IsSigner: false},
	}

	// Serialize instruction data (listing_id parameter)
	listingIDBytes := listingAccountPDA.Bytes()
	data := append(BuyTokensDiscriminator[:], listingIDBytes...)

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// MintConsumptionTokens creates an instruction to mint consumption tokens
func (c *Client) MintConsumptionTokens(consumer, grid, gridAuthority solana.PublicKey, amount uint64) (solana.Instruction, error) {
	if !ValidatePublicKey(consumer) {
		return nil, fmt.Errorf("invalid consumer public key")
	}
	if !ValidatePublicKey(grid) {
		return nil, fmt.Errorf("invalid grid public key")
	}
	if !ValidatePublicKey(gridAuthority) {
		return nil, fmt.Errorf("invalid grid authority public key")
	}
	if !ValidateAmount(amount) {
		return nil, fmt.Errorf("invalid amount")
	}

	consumerAccountPDA, _, err := c.DeriveConsumerAccountPDA(consumer)
	if err != nil {
		return nil, fmt.Errorf("failed to derive consumer account PDA: %w", err)
	}

	gridAccountPDA, _, err := c.DeriveGridAccountPDA(grid)
	if err != nil {
		return nil, fmt.Errorf("failed to derive grid account PDA: %w", err)
	}

	accounts := []*solana.AccountMeta{
		{PublicKey: consumerAccountPDA, IsWritable: true, IsSigner: false},
		{PublicKey: gridAccountPDA, IsWritable: false, IsSigner: false},
		{PublicKey: gridAuthority, IsWritable: true, IsSigner: true},
	}

	// Serialize instruction data
	type InstructionData struct {
		Amount uint64 `borsh:"amount"`
	}

	instructionData := InstructionData{
		Amount: amount,
	}

	serializedData, err := borsh.Serialize(instructionData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize instruction data: %w", err)
	}

	data := append(MintConsumptionTokensDiscriminator[:], serializedData...)

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// MintEnergyTokensForCrossmint creates a complete transaction for minting energy tokens that can be used with Crossmint
func (c *Client) MintEnergyTokensForCrossmint(params MintRecordCreationParams, payer solana.PublicKey) (string, error) {
	// Create the instruction
	instruction, err := c.MintEnergyTokens(params)
	if err != nil {
		return "", fmt.Errorf("failed to create mint energy tokens instruction: %w", err)
	}

	// Get latest blockhash from RPC
	latestBlockhash, err := c.rpcClient.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get latest blockhash: %w", err)
	}

	// Create the transaction
	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		latestBlockhash.Value.Blockhash,
		solana.TransactionPayer(payer),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// Serialize the transaction to bytes
	txBytes, err := transaction.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// Encode to base58 for Crossmint
	base58Transaction := base58.Encode(txBytes)
	return base58Transaction, nil
}

// CreateTransactionForCrossmint creates a base58-encoded transaction from any instruction for Crossmint
func (c *Client) CreateTransactionForCrossmint(instruction solana.Instruction, payer solana.PublicKey, latestBlockhash solana.Hash) (string, error) {
	// Create the transaction
	transaction, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		latestBlockhash,
		solana.TransactionPayer(payer),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// Serialize the transaction to bytes
	txBytes, err := transaction.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// Encode to base58 for Crossmint
	base58Transaction := base58.Encode(txBytes)
	return base58Transaction, nil
}
