// Package pos -- позиция из X и Y
package pos

import (
	"github.com/prospero78/goTV/tv/coordx"
	"github.com/prospero78/goTV/tv/coordy"
	"github.com/prospero78/goTV/tv/types"
)

type TPos struct {
	x types.ICoordX
	y types.ICoordY
}

// New -- возвращает новый IPos
func New() types.IPos {
	return &TPos{
		x: coordx.New(),
		y: coordy.New(),
	}
}

// Y -- возвращает хранимый Y
func (sf *TPos) Y() types.ICoordY {
	return sf.y
}

// GetY -- возвращает хранимый ACoordY
func (sf *TPos) GetY() types.ACoordY {
	return sf.y.Get()
}

// X -- возвращает хранимый X
func (sf *TPos) X() types.ICoordX {
	return sf.x
}

// GetX -- возвращает хранимый ACoordX
func (sf *TPos) GetX() types.ACoordX {
	return sf.x.Get()
}

// Set -- устанавливает хранимую пару
func (sf *TPos) Set(x types.ACoordX, y types.ACoordY) {
	sf.x.Set(x)
	sf.y.Set(y)
}

// Get -- возвращает хранимую пару
func (sf *TPos) Get() (types.ACoordX, types.ACoordY) {
	return sf.x.Get(), sf.y.Get()
}
