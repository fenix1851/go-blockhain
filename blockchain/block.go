package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Блок - это структура данных, которая хранит в себе хэш предыдущего блока, хэш текущего блока и данные.
type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

// createBlock - это функция, которая создает новый блок.
// Для этого она создает новый блок и присваивает ему данные и хэш предыдущего блока.
func CreateBlock(txs []*Transaction, PrevHash []byte) *Block {
	block := &Block{[]byte{}, txs, PrevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis - это функция, которая создает первый блок в блокчейне.
// Она принимает транзакцию, которая будет записана в блок.
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize - это функция, которая сериализует блок.
// Для этого она кодирует блок в байты.
func (b *Block) Serialize() []byte {
	// result - это буфер, в который мы будем записывать байты.
	var result bytes.Buffer
	// encoder - это кодировщик, который будет кодировать блок в байты.
	encoder := gob.NewEncoder(&result)
	// Encode - это функция, которая кодирует блок в байты.
	err := encoder.Encode(b)
	// Если произошла ошибка, то мы останавливаем программу.
	Handle(err)
	// Возвращаем байты.
	return result.Bytes()
}

// DeserializeBlock - это функция, которая десериализует блок.
// Для этого она декодирует байты в блок.
func Deserialize(d []byte) *Block {
	// block - это блок, который мы будем декодировать.
	var block Block
	// decoder - это декодировщик, который будет декодировать байты в блок.
	decoder := gob.NewDecoder(bytes.NewReader(d))
	// Decode - это функция, которая декодирует байты в блок.
	err := decoder.Decode(&block)
	// Если произошла ошибка, то мы останавливаем программу.
	Handle(err)
	// Возвращаем блок.
	return &block
}

func Handle(err error) {
	if err != nil {
		panic(err)
	}
}

// HashTransactions - это функция, которая хэширует транзакции в блоке.
// Для этого она создает слайс байтов, в который записывает хэши транзакций.
func (b *Block) HashTransactions() []byte {
	// txHashes - это слайс байтов, в который мы будем записывать хэши транзакций.
	var txHashes [][]byte
	// txHash - это хэш транзакций.
	var txHash [32]byte
	// Проходимся по транзакциям в блоке.
	for _, tx := range b.Transactions {
		// Добавляем хэш транзакции в слайс байтов.
		txHashes = append(txHashes, tx.ID)
	}
	// Склеиваем слайс байтов в один слайс байтов.
	// Вычисляем хэш слайса байтов.
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	// Возвращаем хэш транзакций.
	return txHash[:]
}
