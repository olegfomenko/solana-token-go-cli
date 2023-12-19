package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/near/borsh-go"
	solana "github.com/olegfomenko/solana-go"
	"github.com/olegfomenko/solana-go/rpc"
	"github.com/olegfomenko/solana-token-cli/types"
)

const (
	PrvKeyPathENV = "SENDER_PRV"
	SolanaRPCENV  = "SOLANA_RPC"
)

func main() {
	var prvPath string
	if prvPath = os.Getenv(PrvKeyPathENV); prvPath == "" {
		panic("Invalid keygen file path: should not be empty")
	}

	prv, err := solana.PrivateKeyFromSolanaKeygenFile(prvPath)
	if err != nil {
		panic(err)
	}

	var (
		networkType  = (os.Args[1:])[0]
		metadataPath = (os.Args[1:])[1]
		mintPubKey   = (os.Args[1:])[2]
	)

	meta, err := readMetadata(metadataPath)
	if err != nil {
		panic(err)
	}

	mint, err := solana.PublicKeyFromBase58(mintPubKey)
	if err != nil {
		panic(err)
	}

	instruction, err := createMetadataInstruction(prv, *meta, mint)
	if err != nil {
		panic(err)
	}

	client := GetRPC(networkType)

	blockhash, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			instruction,
		},
		blockhash.Value.Blockhash,
		solana.TransactionPayer(prv.PublicKey()),
	)
	if err != nil {
		panic(err)
	}

	if _, err = tx.AddSignature(prv); err != nil {
		panic(err)
	}

	binTx, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}

	sig, err := client.SendRawTransactionWithOpts(context.TODO(), binTx, rpc.TransactionOpts{
		SkipPreflight: false,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Metadata created. Transaction hash: %s\n", sig.String())
}

func createMetadataInstruction(prv solana.PrivateKey, metadata types.Data, mint solana.PublicKey) (solana.Instruction, error) {
	const CreateMetadataAccountV3Instruction = 33

	data, err := borsh.Serialize(struct {
		Instruction       uint8
		Data              types.Data
		IsMutable         bool
		CollectionDetails *borsh.Enum
	}{
		Instruction: CreateMetadataAccountV3Instruction,
		Data:        metadata,
		IsMutable:   true,
	})
	if err != nil {
		return nil, err
	}

	metadataAccount, _, err := solana.FindTokenMetadataAddress(mint)
	if err != nil {
		return nil, err
	}

	instruction := solana.NewInstruction(
		solana.TokenMetadataProgramID,
		solana.AccountMetaSlice{
			&solana.AccountMeta{
				PublicKey:  metadataAccount,
				IsWritable: true,
				IsSigner:   false,
			},
			&solana.AccountMeta{
				PublicKey:  mint,
				IsWritable: false,
				IsSigner:   false,
			},
			&solana.AccountMeta{
				PublicKey:  prv.PublicKey(),
				IsWritable: false,
				IsSigner:   true,
			},
			&solana.AccountMeta{
				PublicKey:  prv.PublicKey(),
				IsWritable: false,
				IsSigner:   false,
			},
			&solana.AccountMeta{
				PublicKey:  prv.PublicKey(),
				IsWritable: false,
				IsSigner:   true,
			},
			&solana.AccountMeta{
				PublicKey:  solana.SystemProgramID,
				IsWritable: false,
				IsSigner:   false,
			},
			&solana.AccountMeta{
				PublicKey:  solana.SysVarRentPubkey,
				IsWritable: false,
				IsSigner:   false,
			},
		},
		data,
	)

	return instruction, nil
}

func readMetadata(path string) (metadata *types.Data, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	metadata = &types.Data{}
	if err = json.Unmarshal(data, metadata); err != nil {
		panic(err)
	}

	return metadata, nil
}

func GetRPC(networkType string) *rpc.Client {
	if url := os.Getenv(SolanaRPCENV); url != "" {
		return rpc.New(url)
	}

	const (
		NetworkDevnet  = "devnet"
		NetworkTestnet = "testnet"
		NetworkMainnet = "mainnet"
	)

	switch networkType {
	case NetworkMainnet:
		return rpc.New("https://api.mainnet-beta.solana.com")
	case NetworkTestnet:
		return rpc.New("https://api.testnet.solana.com")
	case NetworkDevnet:
		return rpc.New("https://api.devnet.solana.com")
	default:
		panic("Invalid network type, should be: mainnet|testnet|devnet")
	}
}
