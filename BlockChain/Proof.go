package BlockChain

import (
	. "BlockChainFinalExam/Block"
	. "BlockChainFinalExam/Transactions"
	. "BlockChainFinalExam/utils"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/fatih/color"
)

var MINING_DIFFICULT = 0x8000

func (bc *Blockchain) StartMining() {
	bc.Mining()
	// 使用time.AfterFunc函数创建了一个定时器，它在指定的时间间隔后执行bc.StartMining函数（自己调用自己）。
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
	color.Yellow("mine time: %v\n", time.Now())
}

// Mining 将交易池的交易打包成新的区块
func (bc *Blockchain) Mining() bool {
	bc.Mux.Lock()
	defer bc.Mux.Unlock()

	// 判断交易池是否有交易，如果没有交易，直接返回
	//if len(bc.TransactionPool) == 0 {
	//	//fmt.Println("没有交易，打包失败")
	//	return false
	//}
	// 将挖矿奖励的交易加入交易池
	bc.AddTransaction(MINING_ACCOUNT_ADDRESS, bc.BlockchainAddress, MINING_REWARD, nil, nil)
	nonce, difficulty := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash
	bc.CreateBlock(difficulty, nonce, big.NewInt(int64(GetFileCount())), previousHash)
	log.Println("action=mining, status=success")
	// 向邻居节点发送共识请求
	for _, n := range bc.Neighbors {
		endpoint := fmt.Sprintf("http://%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return true
}

// ProofOfWork 执行共识机制的工作量证明
func (bc *Blockchain) ProofOfWork() (*big.Int, *big.Int) {
	transactions := bc.CopyTransactionPool() // 获取交易池中的交易列表（选择交易？控制交易数量？）
	previousHash := bc.LastBlock().GetHash() // 获取最新区块的哈希值
	nonce := big.NewInt(0)                   // 初始随机数为0
	begin := time.Now()                      // 记录开始时间
	number := bc.LastBlock().GetNumber()     // 获取最新区块的高度
	// 调整挖矿难度
	if getBlockTimestampDifference(bc, len(bc.Chain)-1) < 6e+9 {
		MINING_DIFFICULT += 32
	} else {
		if MINING_DIFFICULT >= 130000 {
			MINING_DIFFICULT -= 32
		}
	}
	// 不断尝试生成满足难度要求的随机数
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULT, number) {
		nonce.Add(nonce, big.NewInt(1)) // 尝试下一个随机数
	}
	// 输出优化信息
	end := time.Now()
	log.Printf("POW spend Time:%f Second", end.Sub(begin).Seconds())
	log.Printf("POW spend Time:%s", end.Sub(begin))
	return nonce, big.NewInt(int64(MINING_DIFFICULT)) // 返回生成的随机数和挖矿难度
}

// ValidProof 难度验证
func (bc *Blockchain) ValidProof(
	nonce *big.Int,
	previousHash [32]byte,
	transactions []*Transaction,
	difficulty int,
	number *big.Int,
) bool {
	// 计算目标值（target）
	target := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil) // 计算2^256
	target = target.Div(target, big.NewInt(int64(difficulty)))      // 将结果除以难度值算出目标值
	// 创建一个临时区块用于计算哈希值
	tmpBlock := Block{
		Nonce:        nonce,
		Number:       number,
		ParentHash:   previousHash,
		Transactions: transactions,
		Timestamp:    0,
	}
	// 计算临时区块的哈希值，并将其转换为大整数
	result := convertBytesToBigInt(tmpBlock.GenerateHash())
	// 比较目标值和计算结果
	return target.Cmp(result) > 0
}

// 将字节数组转换为大整数
func convertBytesToBigInt(b [32]byte) *big.Int {
	return new(big.Int).SetBytes(b[:])
}

// 计算区块之间的时间差
func getBlockTimestampDifference(bc *Blockchain, num int) int {
	if num == 0 {
		return 0
	}
	return int(bc.Chain[num].GetTimestamp() - bc.Chain[num-1].GetTimestamp())
}

// LastBlock 获取链上最后一个区块
func (bc *Blockchain) LastBlock() *Block {

	return bc.Chain[len(bc.Chain)-1]
}

// CopyTransactionPool 复制当前区块链中的交易池（transactionPool）并返回一个新的交易池副本。
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.TransactionPool {
		transactions = append(transactions,
			NewTransaction(t.SenderAddress,
				t.ReceiveAddress,
				t.Value))
	}
	return transactions
}
