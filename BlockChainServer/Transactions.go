package main

import (
	. "BlockChainFinalExam/Transactions"
	"BlockChainFinalExam/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		//返回交易池中的交易列表
		bcs.MethodGetFunc(w)
	case http.MethodPost:
		//接收来自wallet的交易请求
		bcs.MethodPOSTFunc(w, req)
	case http.MethodPut:
		// PUT方法 用于在另据节点同步交易
		bcs.MethodPutFunc(w, req)
	case http.MethodDelete:
		bc := bcs.GetBlockchain()
		bc.ClearTransactionPool()
		io.WriteString(w, string(utils.JsonStatus("success")))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// MethodGetFunc Get请求返回交易池中的交易列表
func (bcs *BlockchainServer) MethodGetFunc(w http.ResponseWriter) {
	// Get:显示交易池的内容，Mine成功后清空交易池
	w.Header().Add("Content-Type", "application/json")
	bc := bcs.GetBlockchain()
	transactions := bc.TransactionPool
	m, _ := json.Marshal(struct {
		Transactions []*Transaction `json:"transactions"`
		Length       int            `json:"length"`
	}{
		Transactions: transactions,
		Length:       len(transactions),
	})
	io.WriteString(w, string(m[:]))
}

// MethodPOSTFunc 接收来自wallet的交易请求
func (bcs *BlockchainServer) MethodPOSTFunc(w http.ResponseWriter, req *http.Request) {
	log.Printf("\n\n\n")
	log.Println("接受到wallet发送的交易")
	//解码请求体中的JSON数据到TransactionRequest结构体
	decoder := json.NewDecoder(req.Body)
	var t TransactionRequest
	err := decoder.Decode(&t)
	if err != nil {
		log.Printf("ERROR: %v", err)
		io.WriteString(w, string(utils.JsonStatus("Decode Transaction失败")))
		return
	}
	log.Println("发送人公钥SenderPublicKey:", *t.SenderPublicKey)
	log.Println("发送人私钥SenderPrivateKey:", *t.SenderBlockchainAddress)
	log.Println("接收人地址RecipientBlockchainAddress:", *t.RecipientBlockchainAddress)
	log.Println("金额Value:", *t.Value)
	log.Println("交易Signature:", *t.Signature)
	//验证交易请求的字段是否完整
	if !t.Validate() {
		log.Println("ERROR: missing field(s)")
		io.WriteString(w, string(utils.JsonStatus("fail")))
		return
	}
	//将发送方的公钥和签名转换为相应的数据类型，获取区块链实例
	publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
	signature := utils.SignatureFromString(*t.Signature)
	bc := bcs.GetBlockchain()
	//根据交易是否成功创建，返回相应的JSON响应
	isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress,
		*t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

	w.Header().Add("Content-Type", "application/json")
	var m []byte
	if !isCreated {
		w.WriteHeader(http.StatusBadRequest)
		m = utils.JsonStatus("fail[from:blockchainServer]")
	} else {
		w.WriteHeader(http.StatusCreated)
		m = utils.JsonStatus("success[from:blockchainServer]")
	}
	io.WriteString(w, string(m))
}

func (bcs *BlockchainServer) MethodPutFunc(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t TransactionRequest
	err := decoder.Decode(&t)
	if err != nil {
		log.Printf("ERROR: %v", err)
		io.WriteString(w, string(utils.JsonStatus("fail")))
		return
	}
	if !t.Validate() {
		log.Println("ERROR: missing field(s)")
		io.WriteString(w, string(utils.JsonStatus("fail")))
		return
	}
	publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
	signature := utils.SignatureFromString(*t.Signature)
	bc := bcs.GetBlockchain()

	isUpdated := bc.AddTransaction(*t.SenderBlockchainAddress,
		*t.RecipientBlockchainAddress, int64(*t.Value), publicKey, signature)

	w.Header().Add("Content-Type", "application/json")
	var m []byte
	if !isUpdated {
		w.WriteHeader(http.StatusBadRequest)
		m = utils.JsonStatus("fail")
	} else {
		m = utils.JsonStatus("success")
	}
	io.WriteString(w, string(m))
}
