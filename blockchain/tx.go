package blockchain

// TXInput - это структура данных, которая хранит в себе ID - уникальный идентификатор транзакции, Out - номер транзакции, Sig - подпись.
type TXInput struct {
	ID []byte
	// Out - это номер транзакции, которая будет использоваться для проверки подписи.
	Out int
	// Sig - это подпись, которая будет использоваться для проверки подписи.
	Sig string
}

// TXOutput - это структура данных, которая хранит в себе Value - значение транзакции, PubKey - публичный ключ.

type TXOutput struct {
	Value int
	// PubKey - это публичный ключ, который будет использоваться для проверки подписи.
	// Публичный ключ - это ключ, который можно открыть и прочитать.
	PubKey string
}

// CanUnlock - это функция, которая проверяет, может ли транзакция разблокировать выход.
// unlockData - это данные, которые будут использоваться для проверки подписи.
func (in *TXInput) CanUnlock(unlockData string) bool {
	// Если подпись равна unlockData, то возвращаем true, иначе false.
	return in.Sig == unlockData
}

// CanBeUnlocked - это функция, которая проверяет, может ли транзакция разблокировать выход.
// unlockData - это данные, которые будут использоваться для проверки подписи.
func (out *TXOutput) CanBeUnlocked(unlockData string) bool {
	return out.PubKey == unlockData
}
