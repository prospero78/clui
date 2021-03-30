// Package coord -- тип координат
package coord

import "github.com/prospero78/goTV/tv/types"

// TCoord -- операции с координатами
type TCoord struct {
	x types.ICoordX
	y types.ICoordY
}

// New -- возвращает новый *TCoord
func New() *TCoord {
	return &TCoord{}
}

// GetX -- возвращает  X
func (sf *TCoord) GetX() types.ICoordX {
	return sf.x
}

// GetY -- возвращает  Y
func (sf *TCoord) GetY() types.ICoordY {
	return sf.y
}

// GetPos -- возвращает позицию точки
func (sf *TCoord) GetPos() (types.ACoordX, types.ACoordY) {
	return sf.x.Get(), sf.y.Get()
}

// SetPos -- устанавливает обе координаты
func (sf *TCoord) SetPos(x types.ACoordX, y types.ACoordY) {
	sf.x.Set(x)
	sf.y.Set(y)
}
