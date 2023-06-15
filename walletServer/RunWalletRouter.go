package main

import (
	"github.com/gorilla/handlers"
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
	// 配置 CORS 中间件
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	// 包装处理器，添加 CORS 中间件
	handler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(http.DefaultServeMux)
	// 启动服务器
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), handler))

	//log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
