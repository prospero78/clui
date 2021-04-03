// Package types -- геометрические примитивы
package types

// IWidth -- интерфейс ширины
type IWidth interface {
	// Get -- возвращает хранимую ширину
	Get() AWidth
	// Set -- устанавливает хранимую ширину
	Set(AWidth) error
}

// IAutoSize -- интерфейс к автоматической подстройке размера
type IAutoSize interface {
	// Is -- возвращает хранимое значение
	Is() AAutoSize
	// Set -- взводит хранимое значение
	Set()
	// Reset -- сбрасывает хранимое значение
	Reset()
	// Change  -- изменяет состояние заданным образом
	Change(AAutoSize)
}
