package main

import (
	"blockchain/blockchain"
	"fmt"
)

//프로토타입

func main() {
	bc := blockchain.NewBlockchain() //새로운 블록체인 생성

	//블록 추가
	bc.AddBlock("Send 1")
	bc.AddBlock("Send 2")

	for _, block := range bc.Blocks {
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Println()
	}
}
