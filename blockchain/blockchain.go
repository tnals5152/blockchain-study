package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const (
	lastHashKey  = "last"
	BlocksBucket = "blocks"
	dbFile       = "chain.db"
)

type Block struct { //block header
	PrevBlockHash []byte //이전 블록의 해시값
	Hash          []byte //해시값
	Timestamp     int64  //시간
	Data          []byte //데이터
	Nonce         int64  //임시값
}

type Blockchain struct { //블록들의 연결
	db   *bolt.DB
	last []byte
}

type blockchainIterator struct { //블록체인 내부 순회를 위한 반복자(구조체)
	db   *bolt.DB
	hash []byte
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

//block을 boltDB에 저장하기 위한 직렬화
func (b *Block) Serialize() []byte {
	var result bytes.Buffer // bytes.Buffer에 현재 Block을 인코드하고 바이트 배열을 반환

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func (bc *Blockchain) AddBlock(data string) {
	block := NewBlock(data, bc.last)

	db := bc.db

	if db == nil {
		log.Panic("db 세팅 먼저")
	}

	bc.db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(BlocksBucket))

		//block DB 저장
		err = bucket.Put(block.Hash, block.Serialize())

		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte(lastHashKey), block.Hash)

		if err != nil {
			log.Panic(err)
		}

		return
	})

	bc.last = block.Hash
}

func (bc *Blockchain) List() {
	bIter := NewBlockchainIterator(bc)

	for bIter.HasNext() {
		block := bIter.Next() // 다음 블록(이전 블록)

		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Println("pow: ", pow.Validate(block))

		fmt.Println("--------------------------------------------------")
	}
}

func (bIter *blockchainIterator) Next() (block *Block) { //반복해서 블록들을 읽어오기 위한 함수
	bIter.db.View(func(tx *bolt.Tx) error { //읽기 전용 (Update - 읽기, 쓰기, Batch - 배치(다중 업데이트))
		bucket := tx.Bucket([]byte(BlocksBucket))

		encodedBlock := bucket.Get(bIter.hash) //마지막 해시 값으로 block조회
		block = DeserializeBlock(encodedBlock)

		bIter.hash = block.PrevBlockHash // block의 이전 블록의 해시 값을 마지막 해시 값으로 저장

		return nil
	})

	return
}

func (bIter *blockchainIterator) HasNext() bool {
	return len(bIter.hash) != 0 // 다음 블록이 없으면 제네시스 블록의 이전 블록이 없으므로 len == 0
}

//byte 배열을 받아 Block으로 반환하는 역직렬화 함수
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

//새 블록 생성
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{prevBlockHash, []byte{}, time.Now().Unix(), []byte(data), 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

// 블록체인을 새로 생성한다
func NewBlockchain() *Blockchain { //genesis block : 블록체인의 가장 첫 블록

	var last []byte
	db, err := bolt.Open(dbFile, 0600, nil) //boltdb 파일 오픈

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(BlocksBucket)) //버킷 가져오기 tx.CreateBucketIfNotExists()

		if bucket != nil {
			//이미 저장된 블록체인이 있음 -> GetBlockchain()함수를 써야 됨
			log.Panic("이미 저장된 블록체인 있음")
			return
		}

		//새로운 블록체인 생성
		bucket, err = tx.CreateBucket([]byte(BlocksBucket))

		if err != nil {
			log.Panic(err)
		}

		//genesis block생성 = 초기 블록
		genesis := NewBlock("Genesis Block", []byte{})

		//genesis블록 db에 저장
		err = bucket.Put(genesis.Hash, genesis.Serialize())

		//last 키를 통해 마지막 블록의 해시 값 저장(처음 생성이니 처음이자 마지막 블록)
		err = bucket.Put([]byte(lastHashKey), genesis.Hash)

		if err != nil {
			log.Panic(err)
		}

		last = genesis.Hash
		return
	})

	return &Blockchain{db, last}
}

// 블록체인을 완전히 새로 생성하지 않고 기존에 있던 블록체인을 얻어올 경우에만 사용
// ex) 블록을 생성할 때와 출력할 때 사용
func GetBlockchain() *Blockchain {

	var last []byte
	db, err := bolt.Open(dbFile, 0600, nil) //boltdb 파일 오픈

	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(BlocksBucket)) //버킷 가져오기 tx.CreateBucketIfNotExists()

		if bucket != nil {
			//버킷이 존재하면 이미 블록체인이 있음
			//블록체인에서 마지막 블록의 해시 값 가져오기
			last = bucket.Get([]byte(lastHashKey))
			return
		}

		log.Panic("버킷 못 찾음")

		return
	})

	return &Blockchain{db, last}
}

func NewBlockchainIterator(bc *Blockchain) *blockchainIterator {
	return &blockchainIterator{bc.db, bc.last}
}
