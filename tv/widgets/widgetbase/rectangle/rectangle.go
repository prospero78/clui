// Package rectangle -- графический примитив прямоугольника
package rectangle

import (
	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/coord"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/size"
)

// TRectangle -- операции с графическим примитивом прямоугольник"
type TRectangle struct {
	coord types.ICoord
	size  types.ISize
}

// New -- возвращает новый *TRectangle
func New() *TRectangle {
	return &TRectangle{
		coord: coord.New(),
		size:  size.New(),
	}
}

// Coord -- возвращает координаты прямоугольника
func (sf *TRectangle) Coord() types.ICoord {
	return sf.coord
}

// Size -- возвращает размер прямоугольника
func (sf *TRectangle) Size() types.ISize {
	return sf.size
}
