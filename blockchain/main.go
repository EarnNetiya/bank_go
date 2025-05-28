package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"goproject-bank/interfaces"
	"time"
	// "fmt"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	PrevHash []byte
	Data     interfaces.BlockchainTransaction
}

func (b *Block) DeriveHash() {
	dataBytes, _ := json.Marshal(b.Data)
	info := bytes.Join([][]byte{dataBytes, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data interfaces.BlockchainTransaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, prevHash, data}
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
// func GetChain() *Blockchain {
// 	if globalChain == nil {
// 		globalChain = InitBlockChain()
// 	}
// 	return globalChain
// }
