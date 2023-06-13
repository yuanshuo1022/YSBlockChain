package BlockChain

import (
	. "BlockChainFinalExam/Block"
	. "BlockChainFinalExam/Transactions"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

// Blockchain 区块链结构
type Blockchain struct {
	TransactionPool   []*Transaction `json:"transactionPool"`
	Chain             []*Block       `json:"chain"`
	BlockchainAddress string         `json:"blockchainAddress"`
	Port              uint16         `json:"-"`
	Mux               sync.Mutex     `json:"-"`
	Neighbors         []string       `json:"-"`
	MuxNeighbors      sync.Mutex     `json:"-"`
}

// NewBlockchain 创建区块链
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {

	bc := new(Blockchain)
	blocks, _ := LoadBlock()
	bc.Chain = blocks
	if len(blocks) == 0 {
		b := &Block{}
		bc.CreateBlock(big.NewInt(0), big.NewInt(0), big.NewInt(1), b.GetHash()) //创世纪块
	}
	bc.BlockchainAddress = blockchainAddress
	bc.Port = port
	return bc
}

// CreateBlock 在区块链上创建新区块
func (bc *Blockchain) CreateBlock(difficulty, nonce, number *big.Int, previousHash [32]byte) *Block {
	b := NewBlock(previousHash, difficulty, nonce, number, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = []*Transaction{}
	err := b.SaveBlock()
	if err != nil {
		log.Fatal("向区块链文件追加块失败:", err)
	}
	// 删除其他节点交易
	for _, n := range bc.Neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return b
}
func (bc *Blockchain) Print() {
	for i, block := range bc.Chain {
		color.Green("%s BLOCK %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	color.Yellow("%s\n\n\n", strings.Repeat("*", 50))
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.Chain,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	v := &struct {
		Blocks *[]*Block `json:"chain"`
	}{
		Blocks: &bc.Chain,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
