// Package autosize -- потокобезопасный признак автоматического размера
package autosize

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TAutoSize -- операции с потокобезопасным признаком автоматического размера
type TAutoSize struct {
	val   types.AAutoSize
	block sync.RWMutex
}

// New -- возвращает новый TAutoSize
func New() types.IAutoSize {
	return &TAutoSize{}
}

// Change -- имзеняет хранимый признак авторазмера
func (sf *TAutoSize) Change(val types.AAutoSize) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}

// Set -- устанавливает хранимый признак авторазмера
func (sf *TAutoSize) Set() {
	sf.block.Lock()
	sf.val = true
	sf.block.Unlock()
}

// Reset -- сбрасывает хранимый признак авторазмера
func (sf *TAutoSize) Reset() {
	sf.block.Lock()
	sf.val = false
	sf.block.Unlock()
}

// Is -- возвращает хранимое значение
func (sf *TAutoSize) Is() types.AAutoSize {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}
