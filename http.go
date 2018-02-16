package blockchain

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"io"
)

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleWriteTransaction).Methods("POST")
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	return muxRouter
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}


func handleWriteTransaction(w http.ResponseWriter, r *http.Request) {
	message := struct {
		From   string
		Target string
		Value  int
	}{}

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := newTransaction(message.Value, message.From, message.Target)
	if err != nil {
		respondWithJSON(w, r, http.StatusServiceUnavailable, err.Error())
		return
	}

	TransactionPull <- *transaction

	respondWithJSON(w, r, http.StatusOK, transaction)
}