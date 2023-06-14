package main

import (
	. "BlockChainFinalExam/BlockChain"
	"BlockChainFinalExam/wallet"
	"github.com/fatih/color"
	"log"
	"net/http"
	"strconv"
)

// GetBlockchain 获取区块链
func (bcs *BlockchainServer) GetBlockchain() *Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		//创建一个矿工钱包实例
		minersWallet := wallet.NewWallet()
		// NewBlockchain与以前的方法不一样,增加了地址和端口2个参数,是为了区别不同的节点
		bc = NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		color.Magenta("===矿工帐号信息====\n")
		color.Magenta("矿工private_key\n %v\n", minersWallet.PrivateKeyStr())
		color.Magenta("矿工publick_key\n %v\n", minersWallet.PublicKeyStr())
		color.Magenta("矿工blockchain_address\n %s\n", minersWallet.BlockchainAddress())
		color.Magenta("===============\n")
	}
	return bc
}
func (bcs *BlockchainServer) Run() {
	bcs.GetBlockchain().Run()

	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions) //GET 方式和  POST方式
	http.HandleFunc("/mine", bcs.Mine)
	http.HandleFunc("/mine/start", bcs.StartMine)
	http.HandleFunc("/amount", bcs.Amount)
	http.HandleFunc("/consensus", bcs.Consensus)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(bcs.Port())), nil))

}
