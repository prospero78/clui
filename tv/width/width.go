// Package width -- потокобезопасная ширина
package width

import (
	"fmt"
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TWidth -- потокобезопасные операции с шириной
type TWidth struct {
	val   types.AWidth
	block sync.RWMutex
}

// New -- возвращает новый IWidth
func New() types.IWidth {
	return &TWidth{}
}

// Get -- возвращает хранимую ширину
func (sf *TWidth) Get() types.AWidth {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает хранимую ширину
func (sf *TWidth) Set(val types.AWidth) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	if val < 0 {
		return fmt.Errorf("TWidth.Set(): val(%v)<0", val)
	}
	sf.val = val
	return nil
}
