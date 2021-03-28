// Package types -- алиасы типов
package types

type ( // Алисы типов для контроля данных
	// AWidgetID -- уникальный ID виджета
	AWidgetID int64

	// ATitle -- заголовок виджета
	ATitle string

	// APackDirect -- направление упаковки виджетов (Vertical или Horizontal)
	APackDirect int

	// AAlign -- направление текста left, right and center
	AAlign int

	// AHitResult -- использование мышки в событии клика
	AHitResult int

	// AEventType -- тип события, генерируемое виджетом
	AEventType int
)
