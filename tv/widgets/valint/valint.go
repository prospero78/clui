// Package valint -- значение int
package valint

import "sync"

// TValInt -- операции со значение int
type TValInt struct {
	val   int
	block sync.RWMutex
}

// New -- возвращает новый *TValInt
func New() *TValInt {
	return &TValInt{}
}

// Get -- возвращает значение
func (sf *TValInt) Get() int {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// SetVal -- устанавливает значение
func (sf *TValInt) Set(val int) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}
