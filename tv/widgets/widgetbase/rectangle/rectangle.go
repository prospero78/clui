// Package rectangle -- графический примитив прямоугольника
package rectangle

import (
	"github.com/prospero78/goTV/tv/widgets/widgetbase/coord"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/size"
)

// TRectangle -- операции с графическим примитивом прямоугольник"
type TRectangle struct {
	*coord.TCoord
	*size.TSize
}

// New -- возвращает новый *TRectangle
func New() *TRectangle {
	return &TRectangle{
		TCoord: coord.New(),
		TSize:  size.New(),
	}
}
