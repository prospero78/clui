// Package types -- графические примитивы
package types

// ICoordX -- интерфейс к координате X
type ICoordX interface {
	// Set -- устанавливает координату X
	Set(ACoordX)
	// Get -- возвращает координату X
	Get() ACoordX
}

// ICoordY -- интерфейс к координате Y
type ICoordY interface {
	// Set -- устанавливает координату Y
	Set(ACoordY)
	// Get -- возвращает координату Y
	Get() ACoordY
}

// ICoord -- интерфейс к типу координат
type ICoord interface {
	// GetX -- возвращает ICoordX
	GetX() ICoordX
	// GetY -- возвращает ICoordY
	GetY() ICoordY
	// Get -- возвращает обе координаты X, Y
	Get() (ACoordX, ACoordY)
	// Set -- устанавливает обе координаты X, Y
	Set(ACoordX, ACoordY)
}

// IHeight -- интерфейс к высоте объекта
type IHeight interface {
	// Set -- устанавливает высоту объекта
	Set(AHeight) error
	// Get -- возвращает высоту объекта
	Get() AHeight
}

// IWidth -- интерфейс к длине объекта
type IWidth interface {
	// Set -- устанавливает длину объекта
	Set(AWidth) error
	// Get -- возвращает длину объекта
	Get() AWidth
}

// ISize -- интерфейс к размеру объекта
type ISize interface {
	// Height -- возвращает объект-высоту
	Height() IHeight
	// Width -- возвращает объект-длину
	Width() IWidth
	// Get -- возвращает ширину и высоту объекта
	Get() (width AWidth, height AHeight)
	// Set -- устанавливает высоту и ширину объекта
	Set(width AWidth, height AHeight) error
}

// IRectangle -- интерфейс к прямоугольнику
type IRectangle interface {
	// Coord -- возвращает координаты
	Coord() ICoord
	// Size -- возвращает размер
	Size() ISize
}
