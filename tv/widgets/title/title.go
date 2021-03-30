// Package title -- заголовок элемента
package title

import "sync"

// TTitle -- заголовок элемента
type TTitle struct {
	val string
	block sync.RWMutex
}

// New -- возвращает новый *TTitle
func New() *TTitle {
	return &TTitle{}
}

// Get -- возвращает хранимый заголовок
func (sf *TTitle) Get() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// Set -- устанавливает хранимый заголовок
func (sf *TTitle) Set(val string) {
	sf.block.Lock()
	sf.val = val
	sf.block.Unlock()
}
