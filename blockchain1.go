package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

//프로토타입

type Block struct { //block header
	PrevBlockHash []byte //이전 블록의 해시값
	Hash          []byte //해시값
	Timestamp     int64  //시간
	Data          []byte //데이터
}

type Blockchain struct { //블록들의 연결
	blocks []*Block
}

func (b *Block) SetHash() {
	header := bytes.Join([][]byte{
		b.PrevBlockHash,
		b.Data,
		[]byte(strconv.FormatInt(b.Timestamp, 16)), //정수를 16진수로 변환
	}, []byte{})
	hash := sha256.Sum256(header)
	b.Hash = hash[:]
}

//새 블록 생성
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data)}
	block.SetHash()

	return block
}

func NewBlockchain() *Blockchain { //genesis block : 블록체인의 가장 첫 블록
	return &Blockchain{[]*Block{
		NewBlock("Genesis Block", []byte{}),
	}}
}

func (bc *Blockchain) AddBlock(data string) {
	block := NewBlock(data, bc.blocks[len(bc.blocks)-1].Hash)
	bc.blocks = append(bc.blocks, block)
}

func main() {
	bc := NewBlockchain() //새로운 블록체인 생성

	//블록 추가
	bc.AddBlock("Send 1")
	bc.AddBlock("Send 2")

	for _, block := range bc.blocks {
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Println()
	}
}
