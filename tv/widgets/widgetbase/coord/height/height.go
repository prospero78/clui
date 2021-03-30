// Package height -- высота объекта
package height

import (
	"fmt"
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// THeight -- операции с высотой объекта
type THeight struct {
	val   types.AHeight
	block sync.RWMutex
}

// New -- возвращает новый IHeight
func New() types.IHeight {
	return &THeight{}
}

// Get -- возвращает высоту объекта
func (sf *THeight) Get() types.AHeight {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает высоту объекта
func (sf *THeight) Set(val types.AHeight) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	if val < 0 {
		return fmt.Errorf("THeight.Set(): val(%v)<0", val)
	}
	sf.val = val
	return nil
}
