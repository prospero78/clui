// Package bardatacell -- ячейка столбика для графика
package bardatacell

import (
	"github.com/prospero78/goTV/tv/widgets/colorbg"
	"github.com/prospero78/goTV/tv/widgets/colorfg"
	"github.com/prospero78/goTV/tv/widgets/title"
	"github.com/prospero78/goTV/tv/widgets/valfloat64"
	"github.com/prospero78/goTV/tv/widgets/valint"
)

// TBarDataCell для ипспользования обратного вызова
// для пользователя для отрисовки с помощью
// настроенных цветов и руны
type TBarDataCell struct {
	*valfloat64.TValFloat64 // Значение столбика текущей отрисовки
	*title.TTitle           // Заголовок столбика

	id *valint.TValInt // Порядовый номер столбика

	barMax   *valfloat64.TValFloat64 // Максимальное значение стобика
	totalMax *valfloat64.TValFloat64 // Самое высокое значение столбика

	// Дефолтные атрибуты и руна для отрисовки на столбике
	*colorfg.TColorFg
	*colorbg.TColorBg
	Ch rune
}

// New -- возвращает новый *TBarDataCell
func New() *TBarDataCell {
	return &TBarDataCell{
		TTitle:      title.New(),
		TColorFg:    colorfg.New(),
		TColorBg:    colorbg.New(),
		TValFloat64: valfloat64.New(),
		barMax:      valfloat64.New(),
		totalMax:    valfloat64.New(),
		id:          valint.New(),
	}
}
