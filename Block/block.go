package Block

import (
	. "BlockChainFinalExam/Transactions"
	. "BlockChainFinalExam/utils"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

type Block struct {
	Nonce        *big.Int       `json:"nonce"`        //随机数
	Timestamp    uint64         `json:"timestamp"`    //区块时间戳
	Number       *big.Int       `json:"number"`       //区块高度
	Difficulty   *big.Int       `json:"difficulty"`   //区块难度
	ParentHash   [32]byte       `json:"parentHash"`   //父区块哈希
	Hash         [32]byte       `json:"hash"`         //本区块hash
	Transactions []*Transaction `json:"transactions"` // 交易列表
}

// NewBlock 新建Block
func NewBlock(previousHash [32]byte, difficulty, nonce, number *big.Int, transaction []*Transaction) *Block {
	block := new(Block)
	block = &Block{
		Timestamp:    uint64(time.Now().UnixNano()),
		Nonce:        nonce,
		ParentHash:   previousHash,
		Number:       number,
		Difficulty:   difficulty,
		Transactions: transaction,
	}

	block.Hash = block.GenerateHash()

	return block
}
func (b *Block) GenerateHash() [32]byte {
	m, _ := json.Marshal(&Block{
		Nonce:      b.Nonce,
		Timestamp:  b.Timestamp,
		Number:     b.Number,
		Difficulty: b.Difficulty,
	})
	return sha256.Sum256([]byte(m))
}

// Serialize 序列化
func Serialize(block *Block) ([]byte, error) {
	type transactionSerialize struct {
		SenderAddress  string `json:"senderAddress"`
		ReceiveAddress string `json:"receiveAddress"`
		Value          int64  `json:"value"`
		Hash           string `json:"hash"`
	}
	b := struct {
		Nonce        *big.Int                `json:"nonce"`        //随机数
		Timestamp    uint64                  `json:"timestamp"`    //区块时间戳
		Number       *big.Int                `json:"number"`       //区块高度
		Difficulty   *big.Int                `json:"difficulty"`   //区块难度
		ParentHash   string                  `json:"parentHash"`   //父区块哈希
		Hash         string                  `json:"hash"`         //本区块hash
		Transactions []*transactionSerialize `json:"transactions"` // 交易列表
	}{
		Nonce:      block.Nonce,
		Timestamp:  block.Timestamp,
		Number:     block.Number,
		Difficulty: block.Difficulty,
		ParentHash: fmt.Sprintf("%x", block.ParentHash),
		Hash:       fmt.Sprintf("%x", block.Hash),
	}
	for _, transaction := range block.Transactions {
		b.Transactions = append(b.Transactions, &transactionSerialize{
			SenderAddress:  transaction.SenderAddress,
			ReceiveAddress: transaction.ReceiveAddress,
			Value:          transaction.Value,
			Hash:           fmt.Sprintf("%x", transaction.Hash),
		})
	}
	result, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Deserialize 反序列化
func Deserialize(data []byte) (*Block, error) {
	type transactionDeserialize struct {
		SenderAddress  string `json:"senderAddress"`
		ReceiveAddress string `json:"receiveAddress"`
		Value          int64  `json:"value"`
		Hash           string `json:"hash"`
	}

	var b struct {
		Nonce        *big.Int                  `json:"nonce"`
		Timestamp    uint64                    `json:"timestamp"`
		Number       *big.Int                  `json:"number"`
		Difficulty   *big.Int                  `json:"difficulty"`
		ParentHash   string                    `json:"parentHash"`
		Hash         string                    `json:"hash"`
		Transactions []*transactionDeserialize `json:"transactions"`
	}

	err := json.Unmarshal(data, &b)
	if err != nil {
		return nil, err
	}

	block := &Block{
		Nonce:        b.Nonce,
		Timestamp:    b.Timestamp,
		Number:       b.Number,
		Difficulty:   b.Difficulty,
		ParentHash:   [32]byte{},
		Hash:         [32]byte{},
		Transactions: make([]*Transaction, len(b.Transactions)),
	}

	// 反序列化父区块哈希
	parentHashBytes, err := hex.DecodeString(b.ParentHash)
	if err != nil {
		return nil, err
	}
	copy(block.ParentHash[:], parentHashBytes)

	// 反序列化本区块哈希
	hashBytes, err := hex.DecodeString(b.Hash)
	if err != nil {
		return nil, err
	}
	copy(block.Hash[:], hashBytes)

	// 反序列化交易列表
	for i, t := range b.Transactions {
		hashTrBytes, err := hex.DecodeString(t.Hash)
		if err != nil {
			return nil, err
		}
		var hash [32]byte
		copy(hash[:], hashTrBytes)

		block.Transactions[i] = &Transaction{
			SenderAddress:  t.SenderAddress,
			ReceiveAddress: t.ReceiveAddress,
			Value:          t.Value,
			Hash:           hash,
		}
	}

	return block, nil
}

// SaveBlock 保存Block信息
func (b *Block) SaveBlock() error {
	file, err := os.OpenFile(SAVEBLOCKFILE, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	serialized, err := Serialize(b)
	if err != nil {
		return err
	}

	err = writer.Write([]string{string(serialized)})
	if err != nil {
		return err
	}

	return nil
}

// GetFileCount 获取文件区块高度
func GetFileCount() int {
	// 打开 CSV 文件
	file, err := os.Open(SAVEBLOCKFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// 创建 CSV Reader
	reader := csv.NewReader(file)
	// 读取所有行
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// 计算行数（数据条数）
	rowCount := len(rows)
	return rowCount
}

// Print 打印区块信息
func (b *Block) Print() {
	log.Printf("%-15v:%30d\n", "number", b.Number)
	log.Printf("%-15v:%30d\n", "nonce", b.Nonce)
	log.Printf("%-15v:%30d\n", "difficulty", b.Difficulty)
	log.Printf("%-15v:%30x\n", "hash", b.Hash)
	log.Printf("%-15v:%30d\n", "timestamp", b.Timestamp)
	log.Printf("%-15v:%30x\n", "parentHash", b.ParentHash)
	for _, t := range b.Transactions {
		t.Print()
	}
}

/*
*Get方法
 */
func (b *Block) GetNonce() *big.Int {
	return b.Nonce
}

func (b *Block) GetTimestamp() uint64 {
	return b.Timestamp
}

func (b *Block) GetNumber() *big.Int {
	return b.Number
}

func (b *Block) GetDifficulty() *big.Int {
	return b.Difficulty
}

func (b *Block) GetParentHash() [32]byte {
	return b.ParentHash
}

func (b *Block) GetHash() [32]byte {
	return b.Hash
}

func (b *Block) GetTransactions() []*Transaction {
	return b.Transactions
}
