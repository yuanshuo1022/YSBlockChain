package main

import (
	"BlockChainFinalExam/AmountResponse"
	"BlockChainFinalExam/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
)

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Call WalletAmount  METHOD:%s\n", req.Method)
	switch req.Method {
	case http.MethodPost:

		var data map[string]interface{}
		// 解析JSON数据

		err := json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			http.Error(w, "无法解析JSON数据", http.StatusBadRequest)
			return
		}

		// 获取JSON字段的值
		blockchainAddress := data["blockchain_address"].(string)
		color.Blue("请求查询账户%s的余额", blockchainAddress)

		// 构建请求数据
		requestData := struct {
			BlockchainAddress string `json:"blockchain_address"`
		}{
			BlockchainAddress: blockchainAddress,
		}

		// 将请求数据编码为JSON
		jsonData, err := json.Marshal(requestData)
		if err != nil {
			fmt.Printf("编码JSON时发生错误:%v", err)
			return
		}

		bcsResp, _ := http.Post(ws.Gateway()+"/amount", "application/json", bytes.NewBuffer(jsonData))

		//返回给客户端
		w.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar AmountResponse.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}

			resp_message := struct {
				Message string `json:"message"`
				Amount  uint64 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			}
			m, _ := json.Marshal(resp_message)
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}
