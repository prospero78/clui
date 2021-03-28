// Package widgetvisible -- потокобезопасный признак видимости виджета
package widgetvisible

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TWidgetVisible -- потокобезопасные операции с видимостью виджета
type TWidgetVisible struct {
	val   types.AVisible
	block sync.RWMutex
}

// New -- возвращает новый TWidgetVisible
func New() TWidgetVisible {
	return TWidgetVisible{
		val: true,
	}
}

// Get -- возвращает хранимое значение видимости
func (sf *TWidgetVisible) Get() types.AVisible {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает признак видимости
func (sf *TWidgetVisible) Set() {
	sf.block.Lock()
	sf.val = true
	sf.block.Unlock()
}


// Reset -- сбрасывает признак видимости
func (sf *TWidgetVisible) Reset() {
	sf.block.Lock()
	sf.val = false
	sf.block.Unlock()
}