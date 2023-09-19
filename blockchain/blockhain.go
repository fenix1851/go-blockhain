package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

// Блокчейн - это структура данных, которая хранит в себе блоки.
type BlockChain struct {
	lastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	currentHash []byte
	Database    *badger.DB
}

// InitBlockChain - это функция, которая создает новый блокчейн.
// Она принимает адрес майнера, который получит вознаграждение за создание первого блока.
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte
	// Если блокчейн уже существует, то мы выходим из функции.
	if DBexists() {
		fmt.Println("Blockchain already exists.")
		runtime.Goexit()
	}
	// Если блокчейн не существует, то мы создаем его.
	// Для этого мы создаем новую базу данных и добавляем в нее первый блок.
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	// Открываем базу данных.
	db, err := badger.Open(opts)
	// Если произошла ошибка, то мы выводим ее и выходим из функции.
	Handle(err)
	// Записываем в базу данных первый блок.
	err = db.Update(func(txn *badger.Txn) error {
		// Создаем первый блок.
		// Для этого мы создаем транзакцию, которая будет храниться в первом блоке.
		cbtx := CoinbaseTx(address, genesisData)
		// Создаем первый блок.
		genesis := Genesis(cbtx)
		// Выводим сообщение о том, что генезис-блок был создан.
		fmt.Println("Genesis proved.")
		// Добавляем первый блок в базу данных.
		err = txn.Set(genesis.Hash, genesis.Serialize())
		// Если произошла ошибка, то мы выводим ее и выходим из функции.
		Handle(err)
		// Добавляем ключ "lh" и значение - хэш первого блока.
		err = txn.Set([]byte("lh"), genesis.Hash)
		// lastHash - это хэш последнего блока в блокчейне.
		lastHash = genesis.Hash
		// Возвращаем ошибку.
		return err
	})
	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain(address string) *BlockChain {
	if !DBexists() {
		fmt.Println("No existing blockchain found. Create one first.")
		runtime.Goexit()
	}

	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	db, err := badger.Open(opts)
	Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	Handle(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// AddBlock - это функция, которая добавляет новый блок в блокчейн.
// Для этого она получает последний блок в блокчейне и создает новый блок.
// После этого она добавляет новый блок в блокчейн.
func (chain *BlockChain) AddBlock(txs []*Transaction) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	Handle(err)

	newBlock := CreateBlock(txs, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.lastHash = newBlock.Hash
		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.lastHash, chain.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.currentHash)
		Handle(err)
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		Handle(err)
		block = Deserialize(encodedBlock)
		return err
	})
	Handle(err)
	iter.currentHash = block.PrevHash
	return block
}

// FindUnspentTransactions - это функция, которая находит все непотраченные транзакции.
// Unspent transactions - это транзакции, которые еще не были потрачены.
// Это значит, что они находятся в блокчейне, но еще не были использованы.
// т.е. они находятся в выходах транзакций.
// но не находятся во входах транзакций.

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	// Создаем слайс, который будет хранить все непотраченные транзакции.
	var unspentTXs []Transaction
	// Создаем слайс, который будет хранить все транзакции, которые были потрачены.
	spentTXOs := make(map[string][]int)
	// Создаем итератор для блокчейна.
	bcIterator := chain.Iterator()
	// Перебираем все блоки в блокчейне.
	for {
		// Получаем текущий блок.
		block := bcIterator.Next()
		// Перебираем все транзакции в блоке.
		for _, tx := range block.Transactions {
			// Получаем ID транзакции.
			txID := hex.EncodeToString(tx.ID)
			// Перебираем все выходы транзакции.
		Outputs:
			for outIdx, out := range tx.Outputs {
				// Если в слайсе spentTXOs есть транзакция с ID, который мы получили выше,
				// то мы перебираем все ее выходы.
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						// Если индекс выхода транзакции совпадает с индексом выхода транзакции,
						// которая была потрачена, то мы переходим к следующей транзакции.
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// Если выход транзакции не был потрачен, то мы добавляем его в слайс непотраченных транзакций.
				if out.CanBeUnlocked(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			// Если транзакция не является кухней, то мы переходим к следующей транзакции.
			if !tx.IsCoinbase() {
				// Перебираем все входы транзакции.
				for _, in := range tx.Inputes {
					// Если вход транзакции может быть разблокирован с помощью адреса,
					// то мы добавляем его в слайс потраченных транзакций.
					if in.CanUnlock(address) {
						// Получаем ID транзакции, которая была потрачена.
						// И добавляем ее в слайс потраченных транзакций.
						// Также мы добавляем индекс выхода транзакции, который был потрачен.
						// Это нужно для того, чтобы мы могли отслеживать, какие выходы транзакций были потрачены.
						// Так как одна транзакция может иметь несколько выходов.
						// И мы должны отслеживать, какие из них были потрачены.
						// И какие еще остались непотраченными.
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}
		// Если текущий блок является первым блоком в блокчейне, то мы выходим из цикла.
		if len(block.PrevHash) == 0 {
			break
		}
	}
	// Возвращаем слайс непотраченных транзакций.
	return unspentTXs
}

func (chain *BlockChain) FindUTXO(address string) []TXOutput {
	// Получаем все непотраченные транзакции.
	var UTXOs []TXOutput
	unspentTransactions := chain.FindUnspentTransactions(address)
	// Перебираем все непотраченные транзакции.
	for _, tx := range unspentTransactions {
		// Перебираем все выходы транзакции.
		for _, out := range tx.Outputs {
			// Если выход транзакции может быть разблокирован с помощью адреса,
			// то мы добавляем его в слайс UTXOs.
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	// Возвращаем слайс UTXOs.
	return UTXOs
}

// FindSpendableOutputs - это функция, которая находит все непотраченные выходы транзакций.
// Unspent outputs - это выходы транзакций, которые еще не были потрачены.
// Это значит, что они находятся в блокчейне, но еще не были использованы.
// т.е. они находятся в выходах транзакций.
// но не находятся во входах транзакций.
// Но в отличие от FindUnspentTransactions, эта функция возвращает не транзакции,
// а выходы транзакций, которые еще не были потрачены.
// Это нужно для того, чтобы мы могли создать новую транзакцию, которая будет потратить эти выходы транзакций.
// Так как в транзакции должны быть указаны входы и выходы.
// И входы должны быть указаны с помощью ID транзакции и индекса выхода транзакции.
// А выходы должны быть указаны с помощью адреса и суммы.

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	// Создаем слайс, который будет хранить все непотраченные выходы транзакций.
	unspentOutputs := make(map[string][]int)
	// Создаем переменную, которая будет хранить сумму всех непотраченных выходов транзакций.
	unspentTXs := chain.FindUnspentTransactions(address)
	accumulated := 0
	// Перебираем все непотраченные транзакции.
Work:
	for _, tx := range unspentTXs {
		// Получаем ID транзакции.
		txID := hex.EncodeToString(tx.ID)
		// Перебираем все выходы транзакции.
		for outIdx, out := range tx.Outputs {
			// Если выход транзакции может быть разблокирован с помощью адреса,
			// то мы добавляем его в слайс непотраченных выходов транзакций.
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	// Возвращаем сумму всех непотраченных выходов транзакций и слайс непотраченных выходов транзакций.
	return accumulated, unspentOutputs
}
