// Package size -- размеры объекта
package size

import "fmt"

// TSize -- операции над размерами объекта
type TSize struct {
	width  int
	height int
}

// New -- возвращает новый *TSize
func New() *TSize {
	return &TSize{}
}

// GetWidth -- возвращает ширину объекта
func (sf *TSize) GetWidth() int {
	return sf.width
}

// SetWidth -- устанавливает ширину объекта
func (sf *TSize) SetWidth(width int) (err error) {
	if width < 0 {
		return fmt.Errorf("TSize.SetWidth(): width(%v)<0", width)
	}
	sf.width = width
}

// GetHeight -- возвращает высоту объекта
func (sf *TSize) GetHeight() int {
	return sf.height
}

// SetHeight -- устанавливает высоту объекта
func (sf *TSize) SetHeight(height int) (err error) {
	if height < 0 {
		return fmt.Errorf("TSize.SetHeight(): height(%v)<0", height)
	}
	sf.height = height
}

// GetPos -- возвращает позицию точки
func (sf *TSize) GetPos() (x, y int) {
	return sf.x, sf.y
}
