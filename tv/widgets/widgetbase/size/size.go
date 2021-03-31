// Package size -- размеры объекта
package size

import (
	"fmt"

	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/size/height"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/size/width"
)

// TSize -- операции над размерами объекта
type TSize struct {
	width  types.IWidth
	height types.IHeight
}

// New -- возвращает новый *ISize
func New() types.ISize {
	return &TSize{
		width:  width.New(),
		height: height.New(),
	}
}

// Width -- возвращает ширину объекта
func (sf *TSize) Width() types.IWidth {
	return sf.width
}

// Height -- возвращает высоту объекта
func (sf *TSize) Height() types.IHeight {
	return sf.height
}

// Get -- возвращает высоту и ширину объекта
func (sf *TSize) Get() (width types.AWidth, height types.AHeight) {
	return sf.width.Get(), sf.height.Get()
}

// Set -- устанавливает высоту и ширину объекта
func (sf *TSize) Set(width types.AWidth, height types.AHeight) (err error) {
	if err = sf.width.Set(width); err != nil {
		return fmt.Errorf("TSize.Set(): in set width, err=%w", err)
	}
	if err = sf.height.Set(height); err != nil {
		return fmt.Errorf("TSize.Set(): in set height, err=%w", err)
	}
	return nil
}
