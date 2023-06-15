package main

import (
	"BlockChainFinalExam/BlockChain"
	"BlockChainFinalExam/utils"
	"BlockChainFinalExam/wallet"
	"flag"
	"fmt"
)

func main() {
	var cache = make(map[string]*BlockChain.Blockchain)
	bc, ok := cache["blockchain"]
	if !ok {
		//创建一个矿工钱包实例
		//minersWallet := wallet.NewWallet()
		minersWallet := wallet.LoadWallet(utils.MINNER_PRIVATE_KEY)
		// NewBlockchain与以前的方法不一样,增加了地址和端口2个参数,是为了区别不同的节点
		bc = BlockChain.NewBlockchain(minersWallet.BlockchainAddress(), 5000)
		cache["blockchain"] = bc
	}
	// 解析命令行参数
	command := flag.String("func", "", "要执行的命令")
	blockhash := flag.String("hash", "", "区块hash")
	number := flag.Uint64("number", 0, "区块号")
	address := flag.String("address", "", "钱包地址")
	trHash := flag.String("trhash", "", "交易哈希")
	// 解析命令行参数
	flag.Parse()
	// 根据命令执行相应的操作
	switch *command {
	case "BlockByHash": // ./YSBC -func=BlockByHash -hash=1f9440fa01ad0a69bcba6585e5db63ff71cbfa2b1ac15e05dc21d6d16f68b72b
		block, err := bc.FindBlockByHash(*blockhash)
		if err != nil {
			fmt.Println("找不到该区块，请检查是否有误:", err)
		} else {
			fmt.Println("找到区块:====")
			block.Print()
		}
	case "BlockByNumber": // ./YSBC -func=BlockByNumber -number=1
		block, err := bc.FindBlockByNumber(*number)
		if err != nil {
			fmt.Println("找不到区块:", err)
		} else {
			fmt.Println("找到区块:")
			block.Print()
		}
	case "TransactionByHash": // ./YSBC -func=TransactionByHash -trhash=9fd1dd9d77cb0604a2fabc85154fb21bda4f95eccc255b08a00209a845ea182d
		transaction := bc.FindTransactionByHash(*trHash)
		if transaction != nil {
			fmt.Println("找到交易:")
			transaction.Print()
		} else {
			fmt.Println("找不到交易")
		}
	case "AddrTransactions": // ./YSBC -func=AddrTransactions -address=9GawPas5KCPZb2EffdzhjAsss73VcNLrK76qjbugbN5q
		transactions := bc.FindUserTransactions(*address)
		fmt.Printf("%s:你的交易列表为:", *address)
		for _, tx := range transactions {
			tx.Print()
		}
	default:
		fmt.Println("未知的命令:输入YSBC -h获取帮助", *command)
	}
}
