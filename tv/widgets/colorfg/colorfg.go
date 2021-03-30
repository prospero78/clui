// Package colorfg -- фон переднего плана
package colorfg

import (
	"sync"

	term "github.com/nsf/termbox-go"
)

// TColorFg -- операции с передним фоном
type TColorFg struct {
	val   term.Attribute
	block sync.RWMutex
}

// New -- возвращает новый *TColorFg
func New() *TColorFg {
	return &TColorFg{}
}

// Set -- устанавливает цвет переднего фона
func (sf *TColorFg) Set(color term.Attribute) {
	sf.block.Lock()
	sf.val = color
	sf.block.Unlock()
}

// Get -- возвращает цвет переднего фона
func (sf *TColorFg) Get() term.Attribute {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}
