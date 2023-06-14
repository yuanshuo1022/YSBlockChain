package main

import (
	"BlockChainFinalExam/Transactions"
	"BlockChainFinalExam/utils"
	"BlockChainFinalExam/wallet"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      *string `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil || strings.TrimSpace(*tr.RecipientBlockchainAddress) == "" ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil || len(*tr.Value) == 0 {
		return false
	}
	return true
}
func (ws *WalletServer) CreateTransaction(
	w http.ResponseWriter,
	req *http.Request) {
	defer req.Body.Close()
	switch req.Method {
	case http.MethodPost:
		var t TransactionRequest
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&t)
		log.Printf("\n\n\n")
		log.Println("发送人公钥SenderPublicKey ==", *t.SenderPublicKey)
		log.Println("发送人私钥SenderPrivateKey ==", *t.SenderPrivateKey)
		log.Println("发送人地址SenderBlockchainAddress ==", *t.SenderBlockchainAddress)
		log.Println("接收人地址RecipientBlockchainAddress ==", *t.RecipientBlockchainAddress)
		log.Println("金额Value ==", *t.Value)
		log.Printf("\n\n\n")

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseUint(*t.Value, 10, 64)
		if err != nil {
			log.Println("ERROR: parse error")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("Validate fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")

		// 交易签名
		transaction := wallet.NewTransaction(privateKey, publicKey,
			*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()
		color.Red("signature:%s", signature)

		bt := &Transactions.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value,
			Signature:                  &signatureStr,
		}
		m, _ := json.Marshal(bt)
		color.Green("提交给BlockServer交易:%s", m)
		buf := bytes.NewBuffer(m)

		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		fmt.Println("=======================resp:", resp)
		if resp.StatusCode == 201 {
			// 201是哪里来的？请参见blockserver  Transactions方法的  w.WriteHeader(http.StatusCreated)语句
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: 非法的HTTP请求方式")
	}
}
