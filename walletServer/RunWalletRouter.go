package main

import (
	"log"
	"net/http"
	"strconv"
)

func (ws *WalletServer) Run() {

	fs := http.FileServer(http.Dir("walletServer/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/walletByPrivatekey", ws.walletByPrivatekey)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
