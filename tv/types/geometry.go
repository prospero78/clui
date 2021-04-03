// Package types -- геометрические примитивы
package types

// IWidth -- интерфейс ширины
type IWidth interface {
	// Get -- возвращает хранимую ширину
	Get() AWidth
	// Set -- устанавливает хранимую ширину
	Set(AWidth) error
}
