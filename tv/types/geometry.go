// Package types -- геометрические примитивы
package types

// IWidth -- интерфейс ширины
type IWidth interface {
	// Get -- возвращает хранимую ширину
	Get() AWidth
	// Set -- устанавливает хранимую ширину
	Set(AWidth) error
}

// IAutoWidth -- интерфейс к автоматической ширине
type IAutoWidth interface {
	// Is -- возвращает хранимое значение
	Is() AAutoWidth
	// Set -- взводит хранимое значение
	Set()
	// Reset -- сбрасывает хранимое значение
	Reset()
	// Change  -- изменяет состояние заданным образом
	Change(AAutoWidth)
}

// IAutoHeight -- интерфейс к автоматической высоте
type IAutoHeight interface {
	// Is -- возвращает хранимое значение
	Is() AAutoHeight
	// Set -- взводит хранимое значение
	Set()
	// Reset -- сбрасывает хранимое значение
	Reset()
	// Change  -- изменяет состояние заданным образом
	Change(AAutoHeight)
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
