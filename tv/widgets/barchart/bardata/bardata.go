// Package bardata --
package bardata

import (
	"github.com/prospero78/goTV/tv/widgets/colorbg"
	"github.com/prospero78/goTV/tv/widgets/colorfg"
	"github.com/prospero78/goTV/tv/widgets/title"
	"github.com/prospero78/goTV/tv/widgets/valfloat64"
)

// TBarData -- показывает информацию в графике. Каждый
// столбик может быть настроено через настройки цвета и
// руны на полоске. Используйте ColorDefault для Fg и Bg,
// и 0 для Ch для отрисовки с помощью BarChart по умоланию.
type TBarData struct {
	*valfloat64.TValFloat64      // Значение столбика
	*title.TTitle                // Название столбика
	*colorfg.TColorFg            // Передний цвет
	*colorbg.TColorBg            // Фоновый цвет
	ch                      rune // Руна для чего-то
}

// New -- возвращает новый *TBarData
func New() *TBarData {
	return &TBarData{
		TValFloat64: valfloat64.New(),
		TTitle:      title.New(),
		TColorFg:    colorfg.New(),
		TColorBg:    colorbg.New(),
	}
}

// SetCh -- устанавливает руну (зачем??)
func (sf *TBarData) SetCh(ch rune) {
	sf.ch = ch
}

// GetCh -- возвращает руну (зачем??)
func (sf *TBarData) GetCh() rune {
	return sf.ch
}
