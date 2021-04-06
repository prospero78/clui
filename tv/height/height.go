// Package height -- потокобезопасная высота
package height

import (
	"fmt"
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// THeight -- потокобезопасные операции с высотой
type THeight struct {
	val   types.AHeight
	block sync.RWMutex
}

// New -- возвращает новый IHeight
func New() types.IHeight {
	return &THeight{}
}

// Get -- возвращает хранимую ширину
func (sf *THeight) Get() types.AHeight {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает хранимую ширину
func (sf *THeight) Set(val types.AHeight) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	if val < 0 {
		return fmt.Errorf("TWidth.Set(): val(%v)<0", val)
	}
	sf.val = val
	return nil
}
