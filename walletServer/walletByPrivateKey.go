package main

import (
	"BlockChainFinalExam/wallet"
	"io"
	"log"
	"net/http"
)

func (ws *WalletServer) walletByPrivatekey(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:

		w.Header().Add("Content-Type", "application/json")
		privatekey := req.FormValue("privatekey")
		myWallet := wallet.LoadWallet(privatekey)
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}
