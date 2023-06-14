package BlockChain

import (
	. "BlockChainFinalExam/Transactions"
	. "BlockChainFinalExam/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// AddTransaction 添加交易
func (bc *Blockchain) AddTransaction(
	sender string,
	recipient string,
	value int64,
	senderPublicKey *ecdsa.PublicKey,
	s *Signature) bool {

	t := NewTransaction(sender, recipient, value)
	//如果是挖矿得到的奖励交易，不验证
	if sender == MINING_ACCOUNT_ADDRESS {
		bc.TransactionPool = append(bc.TransactionPool, t)
		return true
	}

	// 判断有没有足够的余额
	if bc.CalculateTotalAmount(sender) < uint64(value) {
		log.Printf("ERROR: %s ，你的钱包里没有足够的钱", sender)
		return false
	}
	//验证签名
	if bc.VerifyTransactionSignature(senderPublicKey, *s, t) {
		bc.TransactionPool = append(bc.TransactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false
}

// CalculateTotalAmount  获取账户余额(通过循环每一个区块的方式获取每个阶段的余额)
func (bc *Blockchain) CalculateTotalAmount(accountAddress string) uint64 {
	var totalAmount uint64 = 0
	for _, _chain := range bc.Chain {
		for _, _tx := range _chain.GetTransactions() {
			if accountAddress == _tx.ReceiveAddress {
				totalAmount = totalAmount + uint64(_tx.Value)
			}
			if accountAddress == _tx.SenderAddress {
				totalAmount = totalAmount - uint64(_tx.Value)
			}
		}
	}
	return totalAmount
}

// VerifyTransactionSignature 验证交易签名
func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s Signature, t *Transaction) bool {
	m, _ := json.Marshal(struct {
		Sender    string `json:"sender_blockchain_address"`
		Recipient string `json:"recipient_blockchain_address"`
		Value     int64  `json:"value"`
	}{
		Sender:    t.SenderAddress,
		Recipient: t.ReceiveAddress,
		Value:     t.Value,
	})
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

// CreateTransaction 调用addTransaction方法并广播到区块链网络中
func (bc *Blockchain) CreateTransaction(sender string, recipient string, value uint64,
	senderPublicKey *ecdsa.PublicKey, s *Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, int64(value), senderPublicKey, s)

	if isTransacted {
		for _, n := range bc.Neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(),
				senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{
				SenderBlockchainAddress:    &sender,
				RecipientBlockchainAddress: &recipient,
				SenderPublicKey:            &publicKeyStr,
				Value:                      &value,
				Signature:                  &signatureStr}
			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)
			endpoint := fmt.Sprintf("http://%s/transactions", n)
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", endpoint, buf)
			resp, _ := client.Do(req)
			log.Printf("   **  **  **  CreateTransaction : %v", resp)
		}
	}

	return isTransacted
}

// ClearTransactionPool 清空交易池
func (bc *Blockchain) ClearTransactionPool() {
	bc.TransactionPool = bc.TransactionPool[:0]
	blocks, _ := LoadBlock()
	bc.Chain = blocks
}
