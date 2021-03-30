// Package coordy -- координата Y
package coordy

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TCoordY -- операции с координатой Y
type TCoordY struct {
	val   types.ACoordY
	block sync.RWMutex
}

// New -- возвращает новый ICoordY
func New() types.ICoordY {
	return &TCoordY{}
}

// Get -- возвращает координату Y
func (sf *TCoordY) Get() types.ACoordY {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает координату Y
func (sf *TCoordY) Set(val types.ACoordY) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}
