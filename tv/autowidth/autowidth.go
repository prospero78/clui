// Package autowidth -- потокобезопасный признак автоматической ширины
package autowidth

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TAutoWidth -- операции с потокобезопасным признаком автоматической ширины
type TAutoWidth struct {
	val   types.AAutoWidth
	block sync.RWMutex
}

// New -- возвращает новый IAutoWidth
func New() types.IAutoWidth {
	return &TAutoWidth{}
}

// Change -- имзеняет хранимый признак авторазмера
func (sf *TAutoWidth) Change(val types.AAutoWidth) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}

// Set -- устанавливает хранимый признак авторазмера
func (sf *TAutoWidth) Set() {
	sf.block.Lock()
	sf.val = true
	sf.block.Unlock()
}

// Reset -- сбрасывает хранимый признак авторазмера
func (sf *TAutoWidth) Reset() {
	sf.block.Lock()
	sf.val = false
	sf.block.Unlock()
}

// Is -- возвращает хранимое значение
func (sf *TAutoWidth) Is() types.AAutoWidth {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}
