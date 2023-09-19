package blockchain

// ProofOfWork - это механизм, который позволяет нам доказать, что мы выполнили некоторую работу.

// Мы возьмём дату из блока

// Поставим nonce в 0

// Мы будем увеличивать nonce, пока хэш не будет начинаться с 4 нулей

// Когда мы найдём nonce, который даст нам хэш, который начинается с 4 нулей, мы сможем сказать, что мы выполнили некоторую работу

// и мы сможем добавить блок в блокчейн

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Difficulty - это сложность нашего майнинга
// Чем больше сложность, тем больше нулей должно быть в начале хэша
const Difficulty = 18

// ProofOfWork - это структура данных, которая хранит в себе блок и целевое значение
// Целевое значение - это число, которое должно быть больше, чем хэш блока
// Чем больше сложность, тем больше нулей должно быть в начале хэша
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof - это функция, которая создает новый ProofOfWork
// Для этого она создает новый ProofOfWork и присваивает ему блок и целевое значение
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

// Validate - это функция, которая проверяет, является ли хэш блока действительным
// Для этого она сравнивает хэш блока с целевым значением
// Если хэш блока меньше целевого значения, то хэш блока действительный
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// ToHex - это функция, которая преобразует число в байты
// Это нужно для того, чтобы мы могли добавить nonce в хэш
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// InitData - это функция, которая объединяет данные в один слайс байтов
// Это нужно для того, чтобы мы могли добавить nonce в хэш
// nonce - это число, которое мы будем увеличивать, пока хэш не будет начинаться с 4 нулей
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(Difficulty)),
			ToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run - это функция, которая выполняет майнинг
// Мы возьмём дату из блока
// Поставим nonce в 0
// Мы будем увеличивать nonce, пока хэш не будет начинаться с 4 нулей
// Когда мы найдём nonce, который даст нам хэш, который начинается с 4 нулей, мы сможем сказать, что мы выполнили некоторую работу
// и мы сможем добавить блок в блокчейн
func (pow *ProofOfWork) Run() (int, []byte) {
	// Создаём переменную, которая будет хранить хэш в виде big.Int
	var intHash big.Int
	// Создаём переменную, которая будет хранить хэш в виде слайса байтов
	var hash [32]byte
	// Создаём переменную, которая будет хранить nonce
	nonce := 0

	// Будем увеличивать nonce, пока хэш не будет начинаться с 4 нулей
	for nonce < math.MaxInt64 {
		// Создаём слайс байтов, который будет хранить данные
		data := pow.InitData(nonce)
		// Вычисляем хэш
		hash = sha256.Sum256(data)
		// Выводим хэш в консоль
		fmt.Printf("\r%x", hash)
		// Преобразуем хэш в big.Int
		intHash.SetBytes(hash[:])
		// Сравниваем хэш с целевым значением
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	// Выводим пустую строку в консоль
	fmt.Println()
	// Возвращаем nonce и хэш
	return nonce, hash[:]

}
