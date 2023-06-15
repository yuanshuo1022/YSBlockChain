package BlockChain

import (
	. "BlockChainFinalExam/Block"
	. "BlockChainFinalExam/Transactions"
	"BlockChainFinalExam/utils"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/big"
	"net/http"
	"os"
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

func LoadBlock() ([]*Block, error) {
	file, err := os.Open(utils.SAVEBLOCKFILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	blocks := make([]*Block, 0, len(rows))

	for _, row := range rows {
		block, err := deserializeBlockFromCSV(row)
		if err != nil {
			color.Red("无法加载区块", err)
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func deserializeBlockFromCSV(row []string) (*Block, error) {
	jsonString := strings.Join(row, ",")
	var blockData struct {
		Nonce        *big.Int        `json:"nonce"`
		Timestamp    uint64          `json:"timestamp"`
		Number       *big.Int        `json:"number"`
		Difficulty   *big.Int        `json:"difficulty"`
		ParentHash   string          `json:"parentHash"`
		Hash         string          `json:"hash"`
		Transactions json.RawMessage `json:"transactions"`
	}

	err := json.Unmarshal([]byte(jsonString), &blockData)
	if err != nil {
		return nil, err
	}

	parentHashBytes, _ := hex.DecodeString(blockData.ParentHash)
	hashBytes, _ := hex.DecodeString(blockData.Hash)
	var parentHash, hash [32]byte
	copy(parentHash[:], parentHashBytes)
	copy(hash[:], hashBytes)
	block := &Block{
		Nonce:      blockData.Nonce,
		Timestamp:  blockData.Timestamp,
		Number:     blockData.Number,
		Difficulty: blockData.Difficulty,
		ParentHash: parentHash,
		Hash:       hash,
	}

	// 判断"transactions"字段是否存在
	if len(blockData.Transactions) > 0 {
		var transactions []*Transaction
		err := json.Unmarshal(blockData.Transactions, &transactions)
		if err != nil {
			return nil, err
		}
		block.Transactions = transactions
	}

	return block, nil
}

// NewBlockchain 创建区块链
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {

	bc := new(Blockchain)
	bc.Chain, _ = LoadBlock()
	if GetFileCount() == 0 {
		b := &Block{}
		bc.CreateBlock(big.NewInt(0), big.NewInt(0), big.NewInt(0), b.GetHash()) //创世纪块
		//bc.AddTransaction(utils.MINING_ACCOUNT_ADDRESS)
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
	//持久化
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
