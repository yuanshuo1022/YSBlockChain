package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
)

func init() {
	color.Green("==============")
	color.Red("====启动区块链节点=====")
	color.Green("==============")

	log.SetPrefix("Blockchain: ")
}

func main() {

	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()
	fmt.Printf("port::%v \n", *port)
	app := NewBlockchainServer(uint16(*port))
	app.Run()

}
