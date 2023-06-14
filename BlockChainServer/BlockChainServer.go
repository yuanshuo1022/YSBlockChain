package main

import . "BlockChainFinalExam/BlockChain"

var cache map[string]*Blockchain = make(map[string]*Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}
