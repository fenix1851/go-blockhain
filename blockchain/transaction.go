package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

// Transaction - это структура данных, которая хранит в себе ID - уникальный идентификатор транзакции, Inputes - входящие транзакции, Outputs - исходящие транзакции.
type Transaction struct {
	ID      []byte
	Inputes []TXInput
	Outputs []TXOutput
}

// CoinbaseTx - это функция, которая создает новую транзакцию.
func CoinbaseTx(to, data string) *Transaction {
	// Если данные пустые, то мы присваиваем им строку "Reward to 'to'".
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// Создаем новую транзакцию.
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{100, to}
	// Создаем новую транзакцию.
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.ID = tx.Hash()
	// Возвращаем транзакцию.
	return &tx
}

//	Hash - это функция, которая хеширует транзакцию.
func (tx *Transaction) Hash() []byte {
	// encoded - это буфер, в который мы будем записывать байты.
	var encoded bytes.Buffer
	// hash - это хэш транзакции.
	var hash [32]byte
	// encode - это кодировщик, который будет кодировать транзакцию в байты.
	encode := gob.NewEncoder(&encoded)
	// Encode - это функция, которая кодирует транзакцию в байты.
	err := encode.Encode(tx)
	// Если произошла ошибка, то мы выводим ее в консоль.
	Handle(err)
	// Вычисляем хэш транзакции.
	hash = sha256.Sum256(encoded.Bytes())
	// Возвращаем хэш транзакции.
	return hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	// Если длина входящих транзакций равна 1, то это Coinbase транзакция.
	return len(tx.Inputes) == 1 && len(tx.Inputes[0].ID) == 0 && tx.Inputes[0].Out == -1
}

// NewTransaction - это функция, которая создает новую транзакцию.
func NewTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	// Создаем новую транзакцию.
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		fmt.Println("Not enough funds")
		return nil
	}

	// Создаем входящие транзакции.
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)
		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Создаем исходящие транзакции.
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	// Создаем новую транзакцию.
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	// Возвращаем транзакцию.
	return &tx
}
