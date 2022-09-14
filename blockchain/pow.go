package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
)

const targetBits = 24 //poW 합의 알고리즘의 타겟 비트(원래는 변경 가능한 값)

type ProofOfWork struct { //작업 증명을 담당할 구조체 선언
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) //시프트 연산을 통해 2의 uint(256-targetBits) 배의 값으로 변경
	// 맞춰야하는 값으로 인식을 하고??
	// 이러한 값보다 작은 값이 들어오면 트랜잭션이 성공적으로 검증이 되는 걸 의미??
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	maxNumber := ^uint(0) //비트 반전을 통해 제일 큰 값을 반환
	nonce := 0

	fmt.Printf("블록 마이닝 시작  %s\n", pow.block.Data)

	for uint(nonce) < maxNumber {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 { // hashInt가 pow.target보다 작을 경우
			break
		}
		nonce++
	}

	fmt.Println("마이닝 성공")
	return nonce, hash[:]
}

func (pow *ProofOfWork) prepareData(nonce int) []byte { //블록의 값들을 활용해서 병합하는 역할
	return bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	},
		[]byte{})
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)

	return buff.Bytes()
}
