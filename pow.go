package pow

import (
	"math/big"

	"./blockchain/block"
)

const targetBits = 24 //poW 합의 알고리즘의 타겟 비트(원래는 변경 가능한 값)

type ProofOwWork struct { //작업 증명을 담당할 구조체 선언
	block  *block.Block
	target *big.Int
}

func NewProofOwWork(b *block.Block) *ProofOwWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) //시프트 연산을 통해 2의 uint(256-targetBits) 배의 값으로 변경
	// 맞춰야하는 값으로 인식을 하고??
	// 이러한 값보다 작은 값이 들어오면 트랜잭션이 성공적으로 검증이 되는 걸 의미??
	return &ProofOwWork{b, target}
}

// func (pow *ProofOwWork)
