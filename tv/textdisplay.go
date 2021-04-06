package tv

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/autoheight"
	"github.com/prospero78/goTV/tv/autowidth"
	"github.com/prospero78/goTV/tv/types"
)

type TextDisplay struct {
	TBaseControl
	// colorized bool
	topLine   int
	lineCount int

	onDrawLine        func(int) string
	onPositionChanged func(int, int)

	autoHeight types.IAutoHeight
	autoWidth  types.IAutoWidth
}

// TextReader is deprecated due to its confusing name. Use TextDisplay instead.
// In next major library version TextReader will be removed
type TextReader = TextDisplay

func CreateTextReader(parent IControl, width, height int, scale int) *TextDisplay {
	return CreateTextDisplay(parent, width, height, scale)
}

func CreateTextDisplay(parent IControl, width, height int, scale int) *TextDisplay {
	l := &TextDisplay{
		TBaseControl: NewBaseControl(),
		autoHeight:   autoheight.New(),
		autoWidth:    autowidth.New(),
	}

	if height == 0 {
		height = 10
		l.autoHeight.Set()
	}
	if width == 0 {
		width = 20
		l.autoWidth.Set()
	}

	l.SetSize(width, height)
	l.SetConstraints(width, height)
	l.parent = parent

	l.SetScale(scale)

	if parent != nil {
		parent.AddChild(l)
	}

	l.onDrawLine = nil
	l.onPositionChanged = nil

	return l
}

func (l *TextDisplay) drawText() {
	if l.onDrawLine == nil {
		return
	}

	PushAttributes()
	defer PopAttributes()

	bg, fg := RealColor(l.bg, l.Style(), ColorEditBack), RealColor(l.fg, l.Style(), ColorEditText)
	if l.Active() {
		bg, fg = RealColor(l.bg, l.Style(), ColorEditActiveBack), RealColor(l.fg, l.Style(), ColorEditActiveText)
	}
	SetTextColor(fg)
	SetBackColor(bg)

	ind := 0
	for ind < int(l.height) {
		var str string
		if ind+l.topLine < l.lineCount {
			str = l.onDrawLine(ind + l.topLine)
		} else {
			if ind+l.topLine == l.lineCount+5 {
				str = xs.Center("--- THE END ---", int(l.width), " ")
			} else {
				str = ""
			}
		}

		if str != "" {
			str = SliceColorized(str, 0, int(l.width))
			DrawText(l.pos.GetX(), l.pos.GetY()+types.ACoordY(ind), str)
		}

		ind++
	}
}

// Draw repaints the control on its View surface
func (l *TextDisplay) Draw() {
	if l.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	x, y := l.pos.Get()
	w, h := l.Size()

	bg, fg := RealColor(l.bg, l.Style(), ColorEditBack), RealColor(l.fg, l.Style(), ColorEditText)
	if l.Active() {
		bg, fg = RealColor(l.bg, l.Style(), ColorEditActiveBack), RealColor(l.fg, l.Style(), ColorEditActiveText)
	}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, h, ' ')
	l.drawText()
}

func (l *TextDisplay) home() {
	if l.topLine != 0 {
		l.topLine = 0

		if l.onPositionChanged != nil {
			l.onPositionChanged(l.topLine, l.lineCount)
		}
	}
}

func (l *TextDisplay) End() {
	if l.lineCount > 0 && l.topLine != l.lineCount-1 {
		l.topLine = l.lineCount - 1

		if l.onPositionChanged != nil {
			l.onPositionChanged(l.topLine, l.lineCount)
		}
	}
}

func (l *TextDisplay) moveUp(count int) {
	if l.topLine != 0 {
		l.topLine -= count
		if l.topLine < 0 {
			l.topLine = 0
		}

		if l.onPositionChanged != nil {
			l.onPositionChanged(l.topLine, l.lineCount)
		}
	}
}

func (l *TextDisplay) moveDown(count int) {
	if l.lineCount > 0 && l.topLine != l.lineCount-1 {
		l.topLine += count
		if l.topLine > l.lineCount-1 {
			l.topLine = l.lineCount - 1
		}

		if l.onPositionChanged != nil {
			l.onPositionChanged(l.topLine, l.lineCount)
		}
	}
}

func (l *TextDisplay) processMouseClick(ev Event) bool {
	if ev.Key != term.MouseLeft {
		return false
	}

	dy := ev.Y - l.pos.GetY()
	ww := l.height

	if int(dy) < int(l.height)/2 {
		l.moveUp(int(ww) - 1)
	} else {
		l.moveDown(int(ww) - 1)
	}

	return true
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (l *TextDisplay) ProcessEvent(event Event) bool {
	if !l.Active() || !l.Enabled() {
		return false
	}

	switch event.Type {
	case EventKey:
		switch event.Key {
		case term.KeyHome:
			l.home()
			return true
		case term.KeyEnd:
			l.End()
			return true
		case term.KeyArrowUp:
			l.moveUp(1)
			return true
		case term.KeyArrowDown:
			l.moveDown(1)
			return true
		case term.KeyPgup:
			l.moveUp(int(l.height) - 1)
			return true
		case term.KeyPgdn, term.KeySpace:
			l.moveDown(int(l.height) - 1)
			return true
		}

		switch event.Ch {
		case 'k', 'K':
			l.moveUp(1)
			return true
		case 'j', 'J':
			l.moveDown(1)
			return true
		case 'u', 'U':
			l.moveUp(int(l.height) - 1)
			return true
		case 'd', 'D':
			l.moveDown(int(l.height) - 1)
			return true
		default:
			return false
		}
	case EventMouse:
		return l.processMouseClick(event)
	}

	return false
}

// OnDrawLine is called every time the reader is going to display a line
// the argument of the function is the line number to display
func (l *TextDisplay) OnDrawLine(fn func(int) string) {
	l.onDrawLine = fn
}

// OnPositionChanged is called every time the reader changes the top line or
// the total number of lines is changed
// Callback gets two numbers: the current top line, and the total number of
// lines. Top line number starts from 0.
func (l *TextDisplay) OnPositionChanged(fn func(int, int)) {
	l.onPositionChanged = fn
}

func (l *TextDisplay) LineCount() int {
	return l.lineCount
}

func (l *TextDisplay) SetLineCount(lineNo int) {
	if l.topLine == lineNo {
		return
	}

	if lineNo < l.topLine-1 {
		l.topLine = lineNo - 1
	}
	l.lineCount = lineNo

	if l.onPositionChanged != nil {
		l.onPositionChanged(l.topLine, l.lineCount)
	}
}

func (l *TextDisplay) TopLine() int {
	return l.topLine
}

func (l *TextDisplay) SetTopLine(top int) {
	if top < l.lineCount {
		l.topLine = top

		if l.onPositionChanged != nil {
			l.onPositionChanged(l.topLine, l.lineCount)
		}
	}
}
