// Package coordx -- координата Х
package coordx

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TCoordX -- операции с координатой X
type TCoordX struct {
	val   types.ACoordX
	block sync.RWMutex
}

// New -- возвращает новый *ICoordX
func New() types.ICoordX {
	return &TCoordX{}
}

// Get -- возвращает координату X
func (sf *TCoordX) Get() types.ACoordX {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает координату X
func (sf *TCoordX) Set(val types.ACoordX) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}
