package main

import (
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
)

// Amount 处理余额请求
func (bcs *BlockchainServer) Amount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		{
			var data map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&data)
			if err != nil {
				http.Error(w, "无法解析JSON数据", http.StatusBadRequest)
				return
			}
			//从data中获取blockchain_address字段的值，并转换为字符串类型
			blockchainAddress, ok := data["blockchain_address"].(string)
			if !ok {
				http.Error(w, "无效的区块链地址", http.StatusBadRequest)
				return
			}

			color.Green("查询账户：%s 余额请求", blockchainAddress)
			//获取余额
			amount := bcs.GetBlockchain().CalculateTotalAmount(blockchainAddress)

			response := struct {
				Amount uint64 `json:"amount"`
			}{
				Amount: amount,
			}

			m, _ := json.Marshal(response)

			w.Header().Add("Content-Type", "application/json")
			io.WriteString(w, string(m))
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}
