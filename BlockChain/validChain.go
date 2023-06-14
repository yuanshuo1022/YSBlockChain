package BlockChain

import (
	. "BlockChainFinalExam/Block"
	//. "BlockChainFinalExam/utils"
)

func (bc *Blockchain) ValidChain(chain []*Block) bool {
	previousBlock := chain[0] // 获取链中的第一个块作为初始前一个块
	currentIndex := 1         // 从索引1开始，因为索引0已经是初始前一个块
	for currentIndex < len(chain) {
		currentBlock := chain[currentIndex]
		// 检查当前块的previousHash是否与前一个块的哈希相等
		if currentBlock.GetParentHash() != previousBlock.GetHash() {
			return false
		}

		// 检查当前块的工作量证明是否有效
		if !bc.ValidProof(currentBlock.GetNonce(), currentBlock.GetParentHash(), currentBlock.GetTransactions(), MINING_DIFFICULT, currentBlock.GetNumber()) {
			return false
		}

		previousBlock = currentBlock // 更新前一个块为当前块
		currentIndex++               // 增加索引，继续下一个块的验证
	}

	return true
}
