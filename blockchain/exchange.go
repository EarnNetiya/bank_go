package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"goproject-bank/interfaces"
	"time"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	PrevHash []byte
	Data     interfaces.BlockchainTransaction
	Timestamp time.Time // Add timestamp to track block creation time
}

func (b *Block) DeriveHash() {
	dataBytes, _ := json.Marshal(b.Data)
	info := bytes.Join([][]byte{dataBytes, b.PrevHash, []byte(b.Timestamp.String())}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data interfaces.BlockchainTransaction, prevHash []byte) *Block {
	block := &Block{
		[]byte{},
		prevHash,
		data,
		time.Now(), // Set current timestamp
	}
	block.DeriveHash()
	return block
}

func (chain *Blockchain) AddBlock(data interfaces.BlockchainTransaction) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

func Genesis() *Block {
	genesisData := interfaces.BlockchainTransaction{
		SenderAccount:   "system",
		ReceiverAccount: "system",
		Amount:          0,
		Timestamp:       time.Now().String(),
	}
	return CreateBlock(genesisData, []byte{})
}

func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

var Chain = InitBlockChain()

// AddTransaction adds a new transaction to the blockchain
func (chain *Blockchain) AddTransaction(sender, receiver string, amount float64) error {
	transaction := interfaces.BlockchainTransaction{
		SenderAccount:   sender,
		ReceiverAccount: receiver,
		Amount:          amount,
		Timestamp:       time.Now().String(),
	}
	chain.AddBlock(transaction)
	return nil
}

// GetTransactionHistory retrieves all transactions from the blockchain
func (chain *Blockchain) GetTransactionHistory() []interfaces.BlockchainTransaction {
	var transactions []interfaces.BlockchainTransaction
	for _, block := range chain.blocks {
		transactions = append(transactions, block.Data)
	}
	return transactions
}

// VerifyChain checks the integrity of the blockchain
func (chain *Blockchain) VerifyChain() bool {
	for i := 1; i < len(chain.blocks); i++ {
		currentBlock := chain.blocks[i]
		previousBlock := chain.blocks[i-1]

		// Verify current block's hash
		currentBlock.DeriveHash()
		if !bytes.Equal(currentBlock.Hash, currentBlock.Hash) { // Self-check
			return false
		}

		// Verify link to previous block
		if !bytes.Equal(currentBlock.PrevHash, previousBlock.Hash) {
			return false
		}
	}
	return true
}

func (chain *Blockchain) GetBlockchainWithHashes() []interfaces.BlockWithHash {
    var blockchainWithHashes []interfaces.BlockWithHash
    for _, block := range chain.blocks {
        blockWithHash := interfaces.BlockWithHash{
            Hash:      hex.EncodeToString(block.Hash),
            PrevHash:  hex.EncodeToString(block.PrevHash),
            Data:      block.Data,
            Timestamp: block.Timestamp,
        }
        blockchainWithHashes = append(blockchainWithHashes, blockWithHash)
    }
    return blockchainWithHashes
}