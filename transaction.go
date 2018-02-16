package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"fmt"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"errors"
	"gopkg.in/go-playground/validator.v9"
)

type Income struct {
	PrevHashOnTrancation string  `validate:"required"`
	PrevOut              int  `validate:"required"`
}

type Outcome struct {
	Number int  `validate:"required"`
	Value  int  `validate:"required"`
	Target string  `validate:"required"`
}

type Transaction struct {
	Timestamp string    `validate:"required"`
	Hash      string    `validate:"required"`
	Incomes   []Income  `validate:"required"`
	Outcomes  []Outcome `validate:"required"`
	From      string    `validate:"required"`
	Target    string    `validate:"required"`
}

func (tr *Transaction) calculateHash() (hash string, err error) {
	type HashTransaction struct {
		Timestamp string
		Incomes   []Income
		Outcomes  []Outcome
		From      string
		Target    string
	}

	hashTransaction := HashTransaction{
		Timestamp: tr.Timestamp,
		Incomes:   tr.Incomes,
		Outcomes:  tr.Outcomes,
		From:      tr.From,
		Target:    tr.Target,
	}

	jsonBytes, err := json.Marshal(&hashTransaction)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(jsonBytes)

	hashBytes := h.Sum(nil)

	return hex.EncodeToString(hashBytes), nil
}

func newTransaction(value int, from string, target string) (*Transaction, error) {
	outs := findNotUsedOuts(from)
	sum := 0

	chosenOuts := make(map[string][]Outcome)

	for transactionHash, transactionOuts := range outs {
		for _, out := range transactionOuts {
			chosenOuts[transactionHash] = append(chosenOuts[transactionHash], out)
			sum = sum + out.Value
			if sum >= value {
				break
			}
		}
	}

	if sum < value {
		return nil, fmt.Errorf("Not enough outs. Expected: %d, actual: %d", value, sum)
	}

	var newTransactionIncomes []Income
	var newTransactionOutcomes []Outcome
	for transactionHash, transactionOuts := range chosenOuts {
		for _, out := range transactionOuts {
			newTransactionIncomes = append(newTransactionIncomes, Income{
				PrevHashOnTrancation: transactionHash,
				PrevOut:              out.Number,
			})
		}
	}

	newTransactionOutcomes = append(newTransactionOutcomes, Outcome{
		Value:  value,
		Number: 0,
		Target: target,
	})
	if sum > value {
		newTransactionOutcomes = append(newTransactionOutcomes, Outcome{
			Value:  sum - value,
			Number: 1,
			Target: from,
		})
	}

	t := time.Now()
	transaction := &Transaction{
		Timestamp: t.String(),
		Hash:      "",
		Incomes:   newTransactionIncomes,
		Outcomes:  newTransactionOutcomes,
		From:      from,
		Target:    target,
	}

	var err error
	transaction.Hash, err = transaction.calculateHash()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func findNotUsedOuts(from string) map[string][]Outcome {
	outs := make(map[string][]Outcome)

	for _, block := range Blockchain {
		for _, transaction := range block.Transactions {
			var outsByTransaction []Outcome

			for _, out := range transaction.Outcomes {
				if out.Target != from {
					continue
				}

				outsByTransaction = append(outsByTransaction, out)
			}

			outs[transaction.Hash] = outsByTransaction
		}
	}

	for _, block := range Blockchain {
		for _, transaction := range block.Transactions {
			for _, in := range transaction.Incomes {
				if _, exists := outs[in.PrevHashOnTrancation]; exists {
					for range outs[in.PrevHashOnTrancation] {
						trOuts := outs[in.PrevHashOnTrancation]

						for i, out := range trOuts {
							if out.Number == in.PrevOut {
								outs[in.PrevHashOnTrancation] = append(trOuts[:i], trOuts[i+1:]...)
							}
						}
					}

					if len(outs[in.PrevHashOnTrancation]) <= 0 {
						delete(outs, in.PrevHashOnTrancation)
					}
				}
			}
		}
	}

	spew.Dump(outs)

	return outs
}

func isTransactionValid(transaction Transaction) error {
	err := validator.New().Struct(transaction)
	if err != nil {
		return err
	}

	sumIn := 0
	sumOut := 0
	for _, in := range transaction.Incomes {
		prevOutcome, err := getOutcomeByHashAndNumber(in.PrevHashOnTrancation, in.PrevOut)
		if err != nil {
			return errors.New("prev outcome is not found")
		}

		if !isNotUsed(in) {
			return errors.New("income is used")
		}

		if prevOutcome.Target != transaction.From {
			return errors.New("prev outcome target is not valid")
		}

		sumIn += prevOutcome.Value
	}

	for _, out := range transaction.Outcomes {
		sumOut += out.Value
	}

	if sumIn != sumOut {
		return errors.New("sum in is not equals sum out")
	}

	calculatedHash, _ := transaction.calculateHash()
	if calculatedHash != transaction.Hash {
		return errors.New("invalid calculated transaction hash, expect " + calculatedHash)
	}

	return nil
}

func isNotUsed(income Income) bool {
	for _, block := range Blockchain {
		for _, transaction := range block.Transactions {

			for _, in := range transaction.Incomes {
				if in.PrevHashOnTrancation == income.PrevHashOnTrancation {
					return false
				}
			}
		}
	}

	return true
}

func getOutcomeByHashAndNumber(hash string, number int) (*Outcome, error) {
	for _, block := range Blockchain {
		for _, transaction := range block.Transactions {

			for _, out := range transaction.Outcomes {
				if out.Number == number && transaction.Hash == hash {
					return &out, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("not found outcome")
}
