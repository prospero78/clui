// Package coordy -- потокобезопасная координата Y
package coordy

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TCoordY -- потокобезопасные операции с координатой Y
type TCoordY struct {
	val   types.ACoordY
	block sync.RWMutex
}

// New -- возвращает новый ICoordY
func New() types.ICoordY {
	return &TCoordY{}
}

// Get -- возвращает хранимую координату
func (sf *TCoordY) Get() types.ACoordY {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает хранимую координату
func (sf *TCoordY) Set(val types.ACoordY) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}
