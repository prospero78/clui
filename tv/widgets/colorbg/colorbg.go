// Package colorbg -- фон заднего плана
package colorbg

import (
	"sync"

	term "github.com/nsf/termbox-go"
)

// TColorBg -- операции с передним фоном
type TColorBg struct {
	val   term.Attribute
	block sync.RWMutex
}

// New -- возвращает новый *TColorBg
func New() *TColorBg {
	return &TColorBg{}
}

// Set -- устанавливает цвет заднего фона
func (sf *TColorBg) Set(color term.Attribute) {
	sf.block.Lock()
	sf.val = color
	sf.block.Unlock()
}

// GetFg -- возвращает цвет заднего фона
func (sf *TColorBg) Get() term.Attribute {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}
