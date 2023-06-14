package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Printf("ERROR: 非法的HTTP请求方式")
	}
}
