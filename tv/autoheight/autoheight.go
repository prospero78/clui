// Package autoheight -- потокобезопасный признак автоматической высоты
package autoheight

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TAutoHeight -- операции с потокобезопасным признаком автоматической высоты
type TAutoHeight struct {
	val   types.AAutoHight
	block sync.RWMutex
}

// New -- возвращает новый IAutoHight
func New() types.IAutoHight {
	return &TAutoHeight{}
}

// Change -- имзеняет хранимый признак авторазмера
func (sf *TAutoHeight) Change(val types.AAutoHight) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}

// Set -- устанавливает хранимый признак авторазмера
func (sf *TAutoHeight) Set() {
	sf.block.Lock()
	sf.val = true
	sf.block.Unlock()
}

// Reset -- сбрасывает хранимый признак авторазмера
func (sf *TAutoHeight) Reset() {
	sf.block.Lock()
	sf.val = false
	sf.block.Unlock()
}

// Is -- возвращает хранимое значение
func (sf *TAutoHeight) Is() types.AAutoHight {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}
