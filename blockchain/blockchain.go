package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct { //block header
	PrevBlockHash []byte //이전 블록의 해시값
	Hash          []byte //해시값
	Timestamp     int64  //시간
	Data          []byte //데이터
}

type Blockchain struct { //블록들의 연결
	Blocks []*Block
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
	block := NewBlock(data, bc.Blocks[len(bc.Blocks)-1].Hash)
	bc.Blocks = append(bc.Blocks, block)
}
