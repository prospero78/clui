// Package coord -- тип координат
package coord

// TCoord -- операции с координатами
type TCoord struct {
	x int
	y int
}

// New -- возвращает новый *TCoord
func New() *TCoord {
	return &TCoord{
		x: x,
		y: y,
	}
}

// GetX -- возвращает координату X
func (sf *TCoord) GetX() int {
	return sf.x
}

// SetX -- устанавливает координату X
func (sf *TCoord) SetX(x int) {
	sf.x = x
}

// GetY -- возвращает координату Y
func (sf *TCoord) GetY() int {
	return sf.y
}

// SetY -- устанавливает координату Y
func (sf *TCoord) SetY(y int) {
	sf.y = y
}

// GetPos -- возвращает позицию точки
func (sf *TCoord) GetPos() (x, y int) {
	return sf.x, sf.y
}
