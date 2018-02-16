package blockchain

import (
	"time"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Block struct {
	Index        int
	Timestamp    string
	Hash         string
	PrevHash     string
	Transactions []Transaction
}

func generateBlock(oldBlock Block, transaction Transaction) (Block) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transactions = [] Transaction{transaction};
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = newBlock.calculateHash()

	return newBlock
}

func (block *Block) calculateHash() (hash string) {
	type HashBlock struct {
		Index        int
		Timestamp    string
		PrevHash     string
		Transactions []Transaction
	}

	hashTransaction := HashBlock{
		Timestamp:    block.Timestamp,
		Index:        block.Index,
		PrevHash:     block.PrevHash,
		Transactions: block.Transactions,
	}

	jsonBytes, err := json.Marshal(&hashTransaction)
	if err != nil {
		panic(err)
	}

	h := sha256.New()
	h.Write(jsonBytes)

	hashBytes := h.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func isBlockValid(newBlock, oldBlock Block) error {
	if oldBlock.Index+1 != newBlock.Index {
		return errors.New("block index is invalid")
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return errors.New("invalid old block hash")
	}

	if newBlock.calculateHash() != newBlock.Hash {
		return errors.New("invalid calculated block hash")
	}

	return nil
}
