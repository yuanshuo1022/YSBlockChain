package BlockChain

import (
	. "BlockChainFinalExam/Block"
	. "BlockChainFinalExam/utils"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
	"time"
)

// Run 启动区块链
func (bc *Blockchain) Run() {
	// 设置邻居节点
	bc.StartSyncNeighbors()
	// 解决邻居节点冲突问题
	bc.ResolveConflicts()
	// 开始挖矿
	bc.StartMining()
}
func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}

func (bc *Blockchain) SyncNeighbors() {
	bc.MuxNeighbors.Lock()
	defer bc.MuxNeighbors.Unlock()
	bc.SetNeighbors()
}

func (bc *Blockchain) SetNeighbors() {
	bc.Neighbors = FindNeighbors(
		GetHost(), bc.Port,
		NEIGHBOR_IP_RANGE_START, NEIGHBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END)

	color.Blue("邻居节点：%v", bc.Neighbors)
}

// ResolveConflicts 解决区块链之间的冲突
func (bc *Blockchain) ResolveConflicts() bool {
	longestChain := make([]*Block, 0)
	maxLength := len(bc.Chain)

	for _, n := range bc.Neighbors {
		endpoint := fmt.Sprintf("http://%s/chain", n)
		resp, err := http.Get(endpoint)
		if err != nil {
			color.Red("错误：ResolveConflicts GET请求")
			return false
		} else {
			color.Green("正确：ResolveConflicts GET请求")
		}
		if resp.StatusCode == 200 {
			var bcResp Blockchain
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&bcResp)

			if err != nil {
				color.Red("错误：ResolveConflicts Decode")
				return false
			} else {
				color.Green("正确：ResolveConflicts Decode")
			}

			chain := bcResp.Chain
			color.Cyan("ResolveConflicts chain长度：%d", len(chain))
			if len(chain) > maxLength && bc.ValidChain(chain) {
				maxLength = len(chain)
				longestChain = chain
			}
		}
	}

	color.Cyan("ResolveConflicts longestChain长度：%d", len(longestChain))

	if len(longestChain) > 0 {
		bc.Chain = longestChain
		log.Printf("解决冲突：链已替换")
		return true
	}
	log.Printf("解决冲突：链未替换")
	return false
}
