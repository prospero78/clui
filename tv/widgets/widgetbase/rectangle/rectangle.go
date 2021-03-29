// Package rectangle -- графический примитив прямоугольника
package rectangle

// TRectangle -- операции с графическим примитивом прямоугольник"
type TRectangle struct {
	x, y  int
	width int
	hidth int
}

// New -- возвращает новый *TRectangle
func New() *TRectangle {
	return &TRectangle{}
}

// GetX -- возвращает координату X
func (sf *TRectangle) GetX() int {
	return sf.x
}

// SetX -- устанавливает координату X
func (sf *TRectangle) SetX(x int) {
	sf.x = x
}

// GetY -- возвращает координату Y
func (sf *TRectangle) GetY() int {
	return sf.y
}

// SetY -- устанавливает координату Y
func (sf *TRectangle) SetY(y int) {
	sf.y = y
}

// GetWidth -- возвращает ширну
func (sf *TRectangle) GetWidth() int {
	return sf.width
}

// SetWidth -- устанавливает ширину
func (sf *TRectangle) SetWidth(width int) {
	sf.width = width
}

// GetHidth -- возвращает высоту
func (sf *TRectangle) GetHidth() int {
	return sf.hidth
}

// SetHidth -- устанавливает высоту
func (sf *TRectangle) SetHidth(hidth int) {
	sf.hidth = hidth
}
