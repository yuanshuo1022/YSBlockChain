package Block

import (
	. "BlockChainFinalExam/Transactions"
	. "BlockChainFinalExam/utils"
	"crypto/sha256"
	"encoding/json"
	"github.com/fatih/color"
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

	block.Hash = block.generateHash()

	return block
}
func (b *Block) generateHash() [32]byte {
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
	return json.Marshal(block)
}

// Deserialize 反序列化
func Deserialize(data []byte) (*Block, error) {
	block := &Block{}
	err := json.Unmarshal(data, block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// SaveBlock 保存Block信息
func (b *Block) SaveBlock() error {
	//TODO 保存到csv文件中
	return nil
}
func LoadBlock() ([]*Block, error) {
	file, err := os.Open(SAVEBLOCKFILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	blocks := make([]*Block, 0)

	dec := json.NewDecoder(file)

	for dec.More() {
		var block *Block
		if err := dec.Decode(&block); err != nil {
			color.Red("无法加载区块")
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
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
