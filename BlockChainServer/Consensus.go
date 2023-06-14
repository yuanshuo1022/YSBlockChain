package main

import (
	"BlockChainFinalExam/utils"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
)

// Consensus 处理共识请求

func (bcs *BlockchainServer) Consensus(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		{
			color.Cyan("####################Consensus###############")
			//调用GetBlockchain方法获取区块链实例
			bc := bcs.GetBlockchain()
			//并调用其ResolveConflicts方法执行共识算法，解决潜在的分叉情况,并返回是否替换当前区块链的结果
			replaced := bc.ResolveConflicts()
			color.Red("[共识]Consensus replaced: %v\n", replaced)

			w.Header().Add("Content-Type", "application/json")
			if replaced {
				io.WriteString(w, string(utils.JsonStatus("success")))
			} else {
				io.WriteString(w, string(utils.JsonStatus("fail")))
			}
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}
