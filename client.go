package zonnegosdk

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

// ProgramID is the Zonne energy marketplace program ID
var ProgramID = solana.MustPublicKeyFromBase58("AhgnNzTBJRiCUdaFYQsziXKVxfJSLHqc6AqXwfE5zDUa")

// Client represents a client for interacting with the Zonne energy marketplace program
type Client struct {
	rpcClient *rpc.Client
	programID solana.PublicKey
}

// NewClient creates a new Zonne SDK client
func NewClient(rpcEndpoint string) *Client {
	return &Client{
		rpcClient: rpc.New(rpcEndpoint),
		programID: ProgramID,
	}
}

// NewClientWithCustomProgram creates a new client with a custom program ID
func NewClientWithCustomProgram(rpcEndpoint string, programID solana.PublicKey) *Client {
	return &Client{
		rpcClient: rpc.New(rpcEndpoint),
		programID: programID,
	}
}

// GetRPCClient returns the underlying RPC client
func (c *Client) GetRPCClient() *rpc.Client {
	return c.rpcClient
}

// GetProgramID returns the program ID
func (c *Client) GetProgramID() solana.PublicKey {
	return c.programID
}

// Account fetching methods

// GetGridAccount fetches a grid account
func (c *Client) GetGridAccount(ctx context.Context, gridPubkey solana.PublicKey) (*GridAccount, error) {
	gridAccountPDA, _, err := c.DeriveGridAccountPDA(gridPubkey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive grid account PDA: %w", err)
	}

	accountInfo, err := c.rpcClient.GetAccountInfo(ctx, gridAccountPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get grid account info: %w", err)
	}

	if accountInfo.Value == nil {
		return nil, fmt.Errorf("grid account not found")
	}

	var gridAccount GridAccount
	if err := borsh.Deserialize(&gridAccount, accountInfo.Value.Data.GetBinary()[8:]); err != nil {
		return nil, fmt.Errorf("failed to deserialize grid account: %w", err)
	}

	return &gridAccount, nil
}

// GetProducerAccount fetches a producer account
func (c *Client) GetProducerAccount(ctx context.Context, producerPubkey solana.PublicKey) (*ProducerAccount, error) {
	producerAccountPDA, _, err := c.DeriveProducerAccountPDA(producerPubkey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive producer account PDA: %w", err)
	}

	accountInfo, err := c.rpcClient.GetAccountInfo(ctx, producerAccountPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get producer account info: %w", err)
	}

	if accountInfo.Value == nil {
		return nil, fmt.Errorf("producer account not found")
	}

	var producerAccount ProducerAccount
	if err := borsh.Deserialize(&producerAccount, accountInfo.Value.Data.GetBinary()[8:]); err != nil {
		return nil, fmt.Errorf("failed to deserialize producer account: %w", err)
	}

	return &producerAccount, nil
}

// GetConsumerAccount fetches a consumer account
func (c *Client) GetConsumerAccount(ctx context.Context, consumerPubkey solana.PublicKey) (*ConsumerAccount, error) {
	consumerAccountPDA, _, err := c.DeriveConsumerAccountPDA(consumerPubkey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive consumer account PDA: %w", err)
	}

	accountInfo, err := c.rpcClient.GetAccountInfo(ctx, consumerAccountPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get consumer account info: %w", err)
	}

	if accountInfo.Value == nil {
		return nil, fmt.Errorf("consumer account not found")
	}

	var consumerAccount ConsumerAccount
	if err := borsh.Deserialize(&consumerAccount, accountInfo.Value.Data.GetBinary()[8:]); err != nil {
		return nil, fmt.Errorf("failed to deserialize consumer account: %w", err)
	}

	return &consumerAccount, nil
}

// GetListingAccount fetches a listing account
func (c *Client) GetListingAccount(ctx context.Context, producer solana.PublicKey, amount, priceLamports uint64, energyType uint8) (*ListingAccount, error) {
	listingAccountPDA, _, err := c.DeriveListingAccountPDA(producer, amount, priceLamports, energyType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive listing account PDA: %w", err)
	}

	accountInfo, err := c.rpcClient.GetAccountInfo(ctx, listingAccountPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing account info: %w", err)
	}

	if accountInfo.Value == nil {
		return nil, fmt.Errorf("listing account not found")
	}

	var listingAccount ListingAccount
	if err := borsh.Deserialize(&listingAccount, accountInfo.Value.Data.GetBinary()[8:]); err != nil {
		return nil, fmt.Errorf("failed to deserialize listing account: %w", err)
	}

	return &listingAccount, nil
}

// GetMintRecord fetches a mint record
func (c *Client) GetMintRecord(ctx context.Context, producer solana.PublicKey, amount uint64, energyType uint8) (*MintRecord, error) {
	mintRecordPDA, _, err := c.DeriveMintRecordPDA(producer, amount, energyType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive mint record PDA: %w", err)
	}

	accountInfo, err := c.rpcClient.GetAccountInfo(ctx, mintRecordPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to get mint record info: %w", err)
	}

	if accountInfo.Value == nil {
		return nil, fmt.Errorf("mint record not found")
	}

	var mintRecord MintRecord
	if err := borsh.Deserialize(&mintRecord, accountInfo.Value.Data.GetBinary()[8:]); err != nil {
		return nil, fmt.Errorf("failed to deserialize mint record: %w", err)
	}

	return &mintRecord, nil
}

// Transaction building and sending helper
func (c *Client) SendTransaction(ctx context.Context, transaction *solana.Transaction, signers []solana.PrivateKey) (solana.Signature, error) {
	// Get recent blockhash
	recent, err := c.rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	transaction.Message.RecentBlockhash = recent.Value.Blockhash

	// Sign transaction
	_, err = transaction.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		for _, signer := range signers {
			if signer.PublicKey().Equals(key) {
				return &signer
			}
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	sig, err := c.rpcClient.SendTransaction(ctx, transaction)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	return sig, nil
}

// SendAndConfirmTransaction sends a transaction and waits for confirmation
func (c *Client) SendAndConfirmTransaction(ctx context.Context, transaction *solana.Transaction, signers []solana.PrivateKey) (solana.Signature, error) {
	sig, err := c.SendTransaction(ctx, transaction, signers)
	if err != nil {
		return solana.Signature{}, err
	}

	// Wait for confirmation by polling transaction status
	for i := 0; i < 30; i++ { // Try for up to 30 seconds
		status, err := c.rpcClient.GetSignatureStatuses(ctx, true, sig)
		if err == nil && len(status.Value) > 0 && status.Value[0] != nil {
			if status.Value[0].ConfirmationStatus == rpc.ConfirmationStatusFinalized {
				break
			}
		}

		// Wait 1 second before checking again
		select {
		case <-ctx.Done():
			return sig, ctx.Err()
		case <-time.After(time.Second):
		}
	}

	return sig, nil
}
