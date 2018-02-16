package blockchain

import (
	"net/http"
	"github.com/joho/godotenv"
	"os"
	"time"
	"fmt"
)

var Blockchain = []Block{
	{
		Hash:      "genesisBlockHash",
		Index:     0,
		PrevHash:  "",
		Timestamp: time.Now().String(),
		Transactions: []Transaction{
			{
				Hash: "genesisTransactionHash",
				Outcomes: []Outcome{
					{
						Target: "lilit",
						Value:  100,
					},
					{
						Target: "lilit",
						Value:  50,
					},
					{
						Target: "dsxack",
						Value:  70,
					},
				},
			},
		},

	},
	{
		Hash:      "genesisBlockHash2",
		Index:     0,
		PrevHash:  "",
		Timestamp: time.Now().String(),
		Transactions: []Transaction{
			{
				Timestamp: "2018-02-10 02:27:25.707056017 +0300 MSK m=+89.482317396",
				Hash:      "9103a244b5c292f11e8aaefbab7c91e07b2645ac10c329d67cd1eabb861fd337",
				Incomes: []Income{
					{
						PrevHashOnTrancation: "genesisTransactionHash",
						PrevOut:              0,
					},
					{
						PrevHashOnTrancation: "genesisTransactionHash",
						PrevOut:              1,
					},
				},
				Outcomes: []Outcome{
					{
						Number: 0,
						Value:  120,
						Target: "dsxack",
					},
					{
						Number: 1,
						Value:  30,
						Target: "lilit",
					},
				},
				From: "lilit",
				Target: "dsxack",
			},
		},
	},
}

var TransactionPull = make(chan Transaction, 10)

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	go listenTransactionPool(nil)

	m := makeMuxRouter()

	err = http.ListenAndServe(":"+os.Getenv("ADDR"), m)
	if err != nil {
		panic(err)
	}
}

func listenTransactionPool(results chan interface{}) {
	for tr := range TransactionPull {
		err := isTransactionValid(tr)
		if err != nil {
			results <- fmt.Errorf("transaction is bad: %s", err)
			continue
		}

		newBlock := generateBlock(Blockchain[len(Blockchain)-1], tr)
		err = isBlockValid(newBlock, Blockchain[len(Blockchain)-1])
		if err != nil {
			results <- fmt.Errorf("block is invalid: %s", err)
			continue
		}

		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)

		results <- newBlock
	}
}
