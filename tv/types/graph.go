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
	// GetPos -- возвращает обе координаты X, Y
	GetPos() (ACoordX, ACoordY)
	// SetPos -- устанавливает обе координаты X, Y
	SetPos(ACoordX, ACoordY)
}

// IHeight -- интерфейс к высоте объекта
type IHeight interface {
	// Set -- устанавливает высоту объекта
	Set(AHeight) error
	// Get -- возвращает высоту объекта
	Get() AHeight
}

// IWight -- интерфейс к длине объекта
type IWight interface {
	// Set -- устанавливает длину объекта
	Set(AWidth) error
	// Get -- возвращает длину объекта
	Get() AWidth
}

// ISize -- интерфейс к размеру объекта
type ISize interface {
	// GetHeight -- возвращает объект-высоту
	GetHeight() IHeight
	// GetWidth -- возвращает объект-длину
	GetWidth() IWight
	// GetSize -- возвращает ширину и высоту объекта
	GetSize() (width, height int)
	// SetWidth -- устанавливает ширину объекта
	SetWidth(int) error
	// SetSize -- устанавливает высоту и ширину объекта
	SetSize(width, height int) error
}

// IRectangle -- интерфейс к прямоугольнику
type IRectangle interface {
	ICoord
	ISize
}
