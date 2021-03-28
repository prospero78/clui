// Package types -- описатель констант
package types


const (// константы BorderStyle
	ABorderAuto ABorderStyle = iota - 1
	ABorderNone
	ABorderThin
	ABorderThick
)

const (
	// Fixed means 'never change size of the object when its parent resizes'
	Fixed int = 0
	// AutoSize is used only in constructors. It means that the constructor
	// should either calculate the size of an object, e.g. for Label it is its text
	// length, or use default intial values
	AutoSize int = -1
	// KeepSize is used as a placeholder when you want to change only one
	// value and keep other ones untouched. Used in SetSize and SetConstraints
	// methods only
	// Example: control.SetConstraint(10, KeepValue) changes only minimal width
	// of the control and do not change the current minimal control height
	KeepValue int = -1
)