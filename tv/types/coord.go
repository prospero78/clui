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
