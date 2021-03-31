// Package coord -- тип координат
package coord

import (
	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/coord/coordx"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/coord/coordy"
)

// TCoord -- операции с координатами
type TCoord struct {
	x types.ICoordX
	y types.ICoordY
}

// New -- возвращает новый *TCoord
func New() types.ICoord {
	return &TCoord{
		x: coordx.New(),
		y: coordy.New(),
	}
}

// GetX -- возвращает  X
func (sf *TCoord) GetX() types.ICoordX {
	return sf.x
}

// GetY -- возвращает  Y
func (sf *TCoord) GetY() types.ICoordY {
	return sf.y
}

// Get -- возвращает позицию точки
func (sf *TCoord) Get() (types.ACoordX, types.ACoordY) {
	return sf.x.Get(), sf.y.Get()
}

// Set -- устанавливает обе координаты
func (sf *TCoord) Set(x types.ACoordX, y types.ACoordY) {
	sf.x.Set(x)
	sf.y.Set(y)
}
