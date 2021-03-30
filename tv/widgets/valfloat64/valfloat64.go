// Package valfloat64 -- значение float64
package valfloat64

import "sync"

// TValFloat64 -- операции со значение float64
type TValFloat64 struct {
	val   float64
	block sync.RWMutex
}

// New -- возвращает новый *TValFloat64
func New() *TValFloat64 {
	return &TValFloat64{}
}

// Get -- возвращает значение
func (sf *TValFloat64) Get() float64 {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// SetVal -- устанавливает значение
func (sf *TValFloat64) Set(val float64) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}
