package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

const (
	addressChecksumLen = 4
	version            = byte(0x00)
)

// Wallet - это структура данных, которая хранит в себе
// PrivateKey - приватный ключ,
// PublicKey - публичный ключ.
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// Address - это функция, которая возвращает адрес кошелька.
func (w Wallet) Address() []byte {
	// pubKeyHash - это хэш публичного ключа.
	pubKeyHash := publicKeyHash(w.PublicKey)
	// versionedPayload - это версия публичного ключа.
	versionedPayload := append([]byte{version}, pubKeyHash...)
	// checksum - это контрольная сумма.
	checksum := Checksum(versionedPayload)
	// fullPayload - это полная версия публичного ключа.
	fullPayload := append(versionedPayload, checksum...)
	// address - это адрес кошелька.
	address := Base58Encode(fullPayload)
	// Выводим адрес кошелька в консоль.

	// fmt.Printf("Address: %s \n", address)
	// fmt.Printf("Public key: %x \n", w.PublicKey)
	// fmt.Printf("Public key hash: %x \n", pubKeyHash)

	// Возвращаем адрес кошелька.
	return address
}

// NewKeyPair - это функция, которая создает новую пару ключей.
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	// curve - это кривая, которая будет использоваться для создания ключей.
	curve := elliptic.P256()
	// priv - это приватный ключ.
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	// Если произошла ошибка, то мы выводим ее в консоль.
	Handle(err)
	// pub - это публичный ключ.
	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	// Возвращаем приватный ключ и публичный ключ.
	return *priv, pub
}

// publicKeyHash - это функция, которая хеширует публичный ключ.
func publicKeyHash(publicKey []byte) []byte {
	// hash - это хэш публичного ключа.
	hash := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(hash[:])
	// Если произошла ошибка, то мы выводим ее в консоль.
	Handle(err)
	// Возвращаем хэш публичного ключа.
	return hasher.Sum(nil)
}

func Checksum(payload []byte) []byte {
	// firstHash - это первый хэш.
	firstHash := sha256.Sum256(payload)
	// secondHash - это второй хэш.
	secondHash := sha256.Sum256(firstHash[:])
	// Возвращаем второй хэш.
	return secondHash[:addressChecksumLen]
}

// NewWallet - это функция, которая создает новый кошелек.
func NewWallet() *Wallet {
	// Создаем новый кошелек.
	private, public := NewKeyPair()
	// Возвращаем кошелек.
	return &Wallet{private, public}
}
