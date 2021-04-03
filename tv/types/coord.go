// Package types -- описатель простейших координат
package types

// ICoordX -- интерфейс координаты X
type ICoordX interface {
	// Get -- возвращает хранимую координату
	Get() ACoordX
	// Set -- устанавливает хранимую координату
	Set(ACoordX)
}

// ICoordY -- интерфейс координаты Y
type ICoordY interface {
	// Get -- возвращает хранимую координату
	Get() ACoordY
	// Set -- устанавливает хранимую координату
	Set(ACoordY)
}

// IPos -- интерфейс к паре координат XY
type IPos interface {
	// X -- возвращает ICoordX
	X() ICoordX
	// GetX -- возвращает ACoordX
	GetX() ACoordX
	// Y -- возвращает ICoordY
	Y() ICoordY
	// GetY() -- возвращает ACoordY
	GetY() ACoordY
	// Get -- возвращает пару координат
	Get() (ACoordX, ACoordY)
	// Set -- устанавливает парукоординат
	Set(ACoordX, ACoordY)
}
