// Package width -- длина объекта
package width

import (
	"fmt"
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TWidth -- операции с длиной объекта
type TWidth struct {
	val   types.AWidth
	block sync.RWMutex
}

// New -- возвращает новый IWight
func New() types.IWight {
	return &TWidth{}
}

// Get -- возвращает высоту объекта
func (sf *TWidth) Get() types.AWidth {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает высоту объекта
func (sf *TWidth) Set(val types.AWidth) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	if val < 0 {
		return fmt.Errorf("THeight.Set(): val(%v)<0", val)
	}
	sf.val = val
	return nil
}
