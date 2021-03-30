// Package types -- алиасы типов
package types

type ( // Алисы типов для контроля данных
	// ACoordX -- координата Х
	ACoordX int
	// ACoordY -- координата Y
	ACoordY int

	// AHeight --высота объекта
	AHeight int
	// AWidth -- ширина объекта
	AWidth int

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

	// ABorderStyle -- встроенный виджет: auto, none, thin, and thick
	ABorderStyle int

	// ASortOrder -- путь сортировки записей в TableView
	ASortOrder int
	ADragType  int

	// AVisible -- признак видимости виджета
	AVisible bool
)
