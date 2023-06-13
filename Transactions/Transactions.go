package Transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/fatih/color"
	"strings"
)

// Transaction 交易结构
type Transaction struct {
	SenderAddress  string   `json:"senderAddress"`
	ReceiveAddress string   `json:"receiveAddress"`
	Value          int64    `json:"value"`
	Hash           [32]byte `json:"hash omitempty"`
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress *string `json:"recipientBlockchainAddress"`
	SenderPublicKey            *string `json:"senderPublicKey"`
	Value                      *uint64 `json:"value"`
	Signature                  *string `json:"signature"`
}

// NewTransaction 创建交易
func NewTransaction(sender string, receive string, value int64) *Transaction {
	t := &Transaction{
		SenderAddress:  sender,
		ReceiveAddress: receive,
		Value:          value,
	}
	t.Hash = t.generateTrsHash()
	return t
}

// 生成交易hash
func (t *Transaction) generateTrsHash() [32]byte {
	m, _ := json.Marshal(&t)
	return sha256.Sum256(m)
}

// Serialize 方法将 Transaction 结构体序列化为字节数组
func (t *Transaction) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

// UnmarshalJSON 方法实现了对 JSON 数据的自定义反序列化操作
func (t *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction // 定义一个别名，以避免递归调用
	aux := &struct {
		*Alias
		Hash *string `json:"hash"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Hash != nil {
		h, _ := hex.DecodeString(*aux.Hash)
		copy(t.Hash[:], h[:32])
	}
	return nil
}

// Print 打印交易
func (t *Transaction) Print() {
	//TODO 打印交易信息
	color.Red("%s\n", strings.Repeat("~", 30))
	color.Cyan("发送地址             %s\n", t.SenderAddress)
	color.Cyan("接受地址             %s\n", t.ReceiveAddress)
	color.Cyan("金额                 %d\n", t.Value)
	color.Cyan("哈希                 %x\n", t.Hash)

}

// Validate 验证数据，只要有一个为空则为false
func (tr *TransactionRequest) Validate() bool {
	return tr.SenderBlockchainAddress != nil &&
		tr.RecipientBlockchainAddress != nil &&
		tr.SenderPublicKey != nil &&
		tr.Value != nil &&
		tr.Signature != nil
}
