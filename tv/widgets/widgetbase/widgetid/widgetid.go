// Package widgetid -- сквозной ID виджета
package widgetid

import (
	"sync"

	"github.com/prospero78/goTV/tv/types"
)

// TWidgetID -- операции с ID виджета
type TWidgetID struct {
	isMake bool
	val    types.AWidgetID
	block  sync.RWMutex
}

var (
	widgetID TWidgetID // глобальный объект счётчика ID виджетов
)

// GetWidgetID -- возвращает указатель на глобыльный счётчик виджетов
func GetWidgetID() *TWidgetID {
	if widgetID.isMake {
		return &widgetID
	}
	widgetID = TWidgetID{}
	return &widgetID
}

// NextID -- возвращает новый ID виджет
func (sf *TWidgetID) NextID() types.AWidgetID {
	sf.block.Lock()
	defer func() {
		sf.val++
		sf.block.Unlock()
	}()
	return sf.val
}

func init() {
	_ = GetWidgetID()
}
