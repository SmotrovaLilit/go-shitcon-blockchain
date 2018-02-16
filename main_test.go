package blockchain

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"time"
)

func TestTransactions(t *testing.T) {
	results := make(chan interface{})

	go listenTransactionPool(results)
	defer close(TransactionPull)

	Convey("Test transactions", t, func() {
		Convey("Empty transaction", func() {
			TransactionPull <- Transaction{}
			result := <-results

			So(result, ShouldBeError)
		})
		Convey("Income is used", func() {
			TransactionPull <- Transaction{
				Hash: "1d718cc004415418ae64176751ac3aab7c085db0b9a43e17a0ab560a9297570b",
				Timestamp: time.Now().String(),
				Incomes: []Income{
					{
						PrevHashOnTrancation: "genesisTransactionHash",
						PrevOut:              0,
					},
				},
				Outcomes: []Outcome{
					{
						Target: "dsxack",
						Value:  10,
					},
				},
				From: "lilit",
				Target: "dsxack",
			}
			result := <-results
			So(result, ShouldBeError)
			So(result.(error).Error(), ShouldContainSubstring, "income is used")
		})
		Convey("Prev outcome is not found", func() {
			TransactionPull <- Transaction{
				Hash: "1d718cc004415418ae64176751ac3aab7c085db0b9a43e17a0ab560a9297570b",
				Timestamp: time.Now().String(),
				Incomes: []Income{
					{
						PrevHashOnTrancation: "genesisTransactionHash",
						PrevOut:              15,
					},
				},
				Outcomes: []Outcome{
					{
						Target: "dsxack",
						Value:  10,
					},
				},
				From: "lilit",
				Target: "dsxack",
			}
			result := <-results

			So(result, ShouldBeError)
			So(result.(error).Error(), ShouldContainSubstring, "prev outcome is not found")
		})
		Convey("Prev outcome target is not valid", func() {
			TransactionPull <- Transaction{
				Hash: "1d718cc004415418ae64176751ac3aab7c085db0b9a43e17a0ab560a9297570b",
				Timestamp: time.Now().String(),
				Incomes: []Income{
					{
						PrevHashOnTrancation: "9103a244b5c292f11e8aaefbab7c91e07b2645ac10c329d67cd1eabb861fd337",
						PrevOut:              0,
					},
				},
				Outcomes: []Outcome{
					{
						Target: "dsxack",
						Value:  10,
					},
				},
				From: "lilit",
				Target: "dsxack",
			}
			result := <-results

			So(result, ShouldBeError)
			So(result.(error).Error(), ShouldContainSubstring, "prev outcome target is not valid")
		})
		Convey("Sum in is not equals sum out", func() {
			TransactionPull <- Transaction{
				Hash: "1d718cc004415418ae64176751ac3aab7c085db0b9a43e17a0ab560a9297570b",
				Timestamp: time.Now().String(),
				Incomes: []Income{
					{
						PrevHashOnTrancation: "9103a244b5c292f11e8aaefbab7c91e07b2645ac10c329d67cd1eabb861fd337",
						PrevOut:              1,
					},
				},
				Outcomes: []Outcome{
					{
						Target: "dsxack",
						Value:  20,
					},
					{
						Target: "lilit",
						Value:  20,
					},
				},
				From: "lilit",
				Target: "dsxack",
			}

			result := <-results

			So(result, ShouldBeError)
			So(result.(error).Error(), ShouldContainSubstring, "sum in is not equals sum out")
		})
		Convey("Invalid calculated transaction hash", func() {
			TransactionPull <- Transaction{
				Hash: "1d718cc004415418ae64176751ac3aab7c085db0b9a43e17a0ab560a9297570b",
				Timestamp: time.Now().String(),
				Incomes: []Income{
					{
						PrevHashOnTrancation: "9103a244b5c292f11e8aaefbab7c91e07b2645ac10c329d67cd1eabb861fd337",
						PrevOut:              1,
					},
				},
				Outcomes: []Outcome{
					{
						Target: "dsxack",
						Value:  20,
					},
					{
						Target: "lilit",
						Value:  10,
					},
				},
				From: "lilit",
				Target: "dsxack",
			}

			result := <-results

  			So(result, ShouldBeError)
			So(result.(error).Error(), ShouldContainSubstring, "invalid calculated transaction hash")
		})
		Convey("Valid transaction", func() {
			TransactionPull <- Transaction{
				Hash: "c5dbdce2e1dae1b261d23cdb7d15a88b0ddd7683e9efb608e5fbb9fb754bc88e",
				Timestamp: "2018-02-09 23:44:45.510986823 +0300 MSK m=+167.645687132",
				Incomes: []Income{
					{
						PrevHashOnTrancation: "9103a244b5c292f11e8aaefbab7c91e07b2645ac10c329d67cd1eabb861fd337",
						PrevOut:              1,
					},
				},
				Outcomes: []Outcome{
					{
						Target: "dsxack",
						Value:  20,
					},
					{
						Target: "lilit",
						Value:  10,
					},
				},
				From: "lilit",
				Target: "dsxack",
			}

			result := <-results

			var errType *error
			So(result, ShouldNotImplement, errType)
		})
	})
}
