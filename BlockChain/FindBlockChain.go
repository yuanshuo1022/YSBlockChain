package BlockChain

import (
	. "BlockChainFinalExam/Block"
	. "BlockChainFinalExam/Transactions"
	"errors"
	"fmt"
	"math/big"
)

// FindBlockByHash  通过hash查找区块
func (bc *Blockchain) FindBlockByHash(hash string) (*Block, error) {
	for _, b := range bc.Chain {
		if fmt.Sprintf("%x", b.GetHash()) == hash {
			return b, nil
		}
	}
	return nil, errors.New("区块hash对应的区块不存在")
}

// FindBlockByNumber  根据区块号获取区块
func (bc *Blockchain) FindBlockByNumber(number uint64) (*Block, error) {
	if int(number) < len(bc.Chain) { //判断区块号大小
		for _, b := range bc.Chain { //遍历blockchain
			if b.GetNumber().Cmp(big.NewInt(int64(number))) == 0 { //判断是否相等
				return b, nil
			}
		}
	}
	return nil, errors.New("区块号对应的区块不存在")
}

// FindTransactionByHash 通过交易哈希查询该交易
func (bc *Blockchain) FindTransactionByHash(hash string) *Transaction {
	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if hash == fmt.Sprintf("%x", tx.Hash) {
				return tx
			}
		}
	}
	return nil
}

// FindUserTransactions 通过钱包地址Address返回该用户所有交易
func (bc *Blockchain) FindUserTransactions(address string) []Transaction {
	transactions := make([]Transaction, 0)
	// 遍历区块链的每个区块
	for _, block := range bc.Chain {
		// 遍历每个区块中的交易
		for _, tx := range block.Transactions {
			// 检查交易的发送方地址和接收方地址是否与目标地址匹配
			if address == tx.ReceiveAddress || address == tx.SenderAddress {
				// 将符合条件的交易添加到结果集中
				transactions = append(transactions, Transaction{
					SenderAddress:  tx.SenderAddress,
					ReceiveAddress: tx.ReceiveAddress,
					Value:          tx.Value,
					Hash:           tx.Hash,
				})
			}
		}
	}

	return transactions
}
