package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

const (
	subsidy = 10 // BTC btc단위
)

type Transaction struct { // 하나의 트랜잭션은 다수의 입출력이 있을 수 있으므로 배열로 표현
	ID   []byte
	Vin  []TXInput  // 입력값
	Vout []TXOutput // 출력값
}

//자본의 흐름
type TXInput struct {
	Txid      []byte // 트랜잭션의 ID
	Vout      int    // 트랜잭션이 가진 출력값의 인덱스
	ScriptSig string // 디지털 서명
}

//트랜잭션의 출력은 코인(value)을 의미
//어디로(ScriptPubKey)에 대해선 받는 사람의 공개키
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// 트랜잭션의 아이디 설정
// 트랜잭션을 직렬화 하고 해시해서 ID로 지정
func (tx *Transaction) SetID() {
	buf := new(bytes.Buffer)

	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(buf.Bytes())
	tx.ID = hash[:]
}

//코인베이스 트랜잭션인지 반환
// 거래 처리 시 코인베이스 트랜잭션에 입력 값 X -> 처리해야 되는 부분 때문에
func (tx *Transaction) IsCoinbase() bool {
	return bytes.Compare(tx.Vin[0].Txid, []byte{}) == 0 &&
		tx.Vin[0].Vout == -1 &&
		len(tx.Vin) == 1
}

// 입력의 잠금-해제 함수
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// 출력의 잠금-해제 함수
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func NewTransaction(vin []TXInput, vout []TXOutput) *Transaction {
	tx := &Transaction{nil, vin, vout}
	tx.SetID()

	return tx
}

// 블록을 채굴하면 채굴자에게 보상을 주기 위한 제일 첫 번째 트랜잭션
// 입력 X 채굴자에게 보상을 지급하기 위한 출력만 존재
func NewCoinbaseTx(data, to string) *Transaction {
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}

	return NewTransaction([]TXInput{txin}, []TXOutput{txout})
}

// 코인 전송
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERRPR: NOt enough funds")
	}

	// 입력 리스트 생성
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)

		log.Println(err)

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// 출력 리스트 생성
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
