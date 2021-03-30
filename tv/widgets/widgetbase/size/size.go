// Package size -- размеры объекта
package size

import (
	"fmt"

	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/coord/height"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/coord/width"
)

// TSize -- операции над размерами объекта
type TSize struct {
	width  types.IWight
	height types.IHeight
}

// New -- возвращает новый *ISize
func New() types.ISize {
	return &TSize{
		width:  width.New(),
		height: height.New,
	}
}

// GetWidth -- возвращает ширину объекта
func (sf *TSize) GetWidth() types.IWight {
	return sf.width
}

// GetHeight -- возвращает высоту объекта
func (sf *TSize) GetHeight() types.IHeight {
	return sf.height
}

// GetSize -- возвращает высоту и ширину объекта
func (sf *TSize) GetSize() (width types.AWidth, height types.AHeight) {
	return sf.width.Get(), sf.height.Get()
}

// SetSize -- устанавливает высоту и ширину объекта
func (sf *TSize) SetSize(width types.AWidth, height types.AHeight) (err error) {
	if err = sf.width.Set(width); err != nil {
		return fmt.Errorf("TSize.SetSize(): in set width, err=%w", err)
	}
	if err = sf.height.Set(height); err != nil {
		return fmt.Errorf("TSize.SetSize(): in set height, err=%w", err)
	}
	return nil
}
