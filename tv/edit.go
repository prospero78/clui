package tv

import (
	"strings"

	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/types"
)

// OnChange sets the callback that is called when EditField content is changed
func (e *TEditField) OnChange(fn func(Event)) {
	e.onChange = fn
}

// OnKeyPress sets the callback that is called when a user presses a Key while
// the controls is active. If a handler processes the key it should return
// true. If handler returns false it means that the default handler will
// process the key
func (e *TEditField) OnKeyPress(fn func(term.Key, rune) bool) {
	e.onKeyPress = fn
}

// SetTitle changes the EditField content and emits OnChage eventif the new value does not equal to old one
func (e *TEditField) SetTitle(title string) {
	e.setTitleInternal(title)
	e.offset = 0
	e.End()
}

func (e *TEditField) setTitleInternal(title string) {
	if e.title != title {
		e.title = title

		if e.onChange != nil {
			ev := Event{Msg: title}
			e.onChange(ev)
		}
	}

	if title == "" {
		e.cursorPos = types.ACoordX(xs.Len(title))
	}
}

// Repaint draws the control on its View surface
func (e *TEditField) Draw() {
	if e.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	x, y := e.Pos().Get()
	w, _ := e.Size()

	parts := []rune(SysObject(ObjEdit))
	chLeft, chRight := string(parts[0]), string(parts[1])
	chStar := "*"
	if len(parts) > 3 {
		chStar = string(parts[3])
	}

	var textOut string
	curOff := 0
	if e.offset == 0 && xs.Len(e.title) < int(e.width.Get()) {
		if e.showStars {
			textOut = strings.Repeat(chStar, xs.Len(e.title))
		} else {
			textOut = e.title
		}
	} else {
		fromIdx := types.ACoordX(0)
		toIdx := types.ACoordX(0)
		if e.offset == 0 {
			toIdx = types.ACoordX(e.width.Get() - 1)
			if e.showStars {
				textOut = strings.Repeat(chStar, int(toIdx)) + chRight
			} else {
				textOut = xs.Slice(e.title, 0, int(toIdx)) + chRight
			}
			curOff = -e.offset
		} else {
			curOff = 1 - e.offset
			fromIdx = types.ACoordX(e.offset)
			if int(e.width.Get()-1) <= xs.Len(e.title)-e.offset {
				toIdx = types.ACoordX(e.offset + int(e.width.Get()) - 2)
				if e.showStars {
					textOut = chLeft + strings.Repeat(chStar, int(toIdx-fromIdx)) + chRight
				} else {
					textOut = chLeft + xs.Slice(e.title, int(fromIdx), int(toIdx)) + chRight
				}
			} else {
				if e.showStars {
					textOut = chLeft + strings.Repeat(chStar, xs.Len(e.title)-int(fromIdx))
				} else {
					textOut = chLeft + xs.Slice(e.title, int(fromIdx), -1)
				}
			}
		}
	}

	fg, bg := RealColor(e.fg, e.Style(), ColorEditText), RealColor(e.bg, e.Style(), ColorEditBack)
	if !e.Enabled() {
		fg, bg = RealColor(e.fg, e.Style(), ColorDisabledText), RealColor(e.fg, e.Style(), ColorDisabledBack)
	} else if e.Active() {
		fg, bg = RealColor(e.fg, e.Style(), ColorEditActiveText), RealColor(e.bg, e.Style(), ColorEditActiveBack)
	}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, 1, ' ')
	DrawRawText(x, y, textOut)
	if e.Active() {
		SetCursorPos(e.cursorPos+e.pos.GetX()+types.ACoordX(curOff), e.pos.GetY())
	}
}

func (e *TEditField) InsertRune(ch rune) {
	if e.readonly {
		return
	}

	if e.maxWidth > 0 && xs.Len(e.title) >= e.maxWidth {
		return
	}

	idx := e.cursorPos

	switch {
	case idx == 0:
		e.setTitleInternal(string(ch) + e.title)
	case int(idx) >= xs.Len(e.title):
		e.setTitleInternal(e.title + string(ch))
	default:
		e.setTitleInternal(xs.Slice(e.title, 0, int(idx)) + string(ch) + xs.Slice(e.title, int(idx), -1))
	}

	e.cursorPos++

	if int(e.cursorPos) >= int(e.width.Get()) {
		if e.offset == 0 {
			e.offset = 2
		} else {
			e.offset++
		}
	}
}

func (e *TEditField) Backspace() {
	if e.title == "" || e.cursorPos == 0 || e.readonly {
		return
	}

	length := xs.Len(e.title)

	switch {
	case int(e.cursorPos) >= length:
		e.cursorPos--
		e.setTitleInternal(xs.Slice(e.title, 0, length-1))
	case e.cursorPos == 1:
		e.cursorPos = 0
		e.setTitleInternal(xs.Slice(e.title, 1, -1))
		e.offset = 0
	default:
		e.cursorPos--
		e.setTitleInternal(xs.Slice(e.title, 0, int(e.cursorPos)) + xs.Slice(e.title, int(e.cursorPos+1), -1))
	}

	if length-1 < int(e.width.Get()) {
		e.offset = 0
	}
}

func (e *TEditField) Del() {
	length := xs.Len(e.title)

	if e.title == "" || int(e.cursorPos) == length || e.readonly {
		return
	}

	if int(e.cursorPos) == length-1 {
		e.setTitleInternal(xs.Slice(e.title, 0, length-1))
	} else {
		e.setTitleInternal(xs.Slice(e.title, 0, int(e.cursorPos)) + xs.Slice(e.title, int(e.cursorPos+1), -1))
	}

	if length-1 < int(e.width.Get()) {
		e.offset = 0
	}
}

func (e *TEditField) CharLeft() {
	if e.cursorPos == 0 || e.title == "" {
		return
	}

	if int(e.cursorPos) == e.offset {
		e.offset--
	}

	e.cursorPos--
}

func (e *TEditField) CharRight() {
	length := xs.Len(e.title)
	if int(e.cursorPos) == length || e.title == "" {
		return
	}

	e.cursorPos++
	if int(e.cursorPos) != length && int(e.cursorPos) >= e.offset+int(e.width.Get())-2 {
		e.offset++
	}
}

func (e *TEditField) Home() {
	e.offset = 0
	e.cursorPos = 0
}

func (e *TEditField) End() {
	length := xs.Len(e.title)
	e.cursorPos = types.ACoordX(length)

	if length < int(e.width.Get()) {
		return
	}

	e.offset = length - (int(e.width.Get()) - 2)
}

// Clear empties the EditField and emits OnChange event
func (e *TEditField) Clear() {
	e.Home()
	e.setTitleInternal("")
}

// SetMaxWidth sets the maximum lenght of the EditField text. If the current text is longer it is truncated
func (e *TEditField) SetMaxWidth(w int) {
	e.maxWidth = w
	if w > 0 && xs.Len(e.title) > w {
		e.title = xs.Slice(e.title, 0, w)
		e.End()
	}
}

// MaxWidth returns the current maximum text length. Zero means no limit
func (e *TEditField) MaxWidth() int {
	return e.maxWidth
}

// SetSize changes control size. Constant DoNotChange can be
// used as placeholder to indicate that the control attrubute
// should be unchanged.
// Method does nothing if new size is less than minimal size
// EditField height cannot be changed - it equals 1 always
func (e *TEditField) SetSize(width, height int) {
	if width != KeepValue && (width > 1000 || width < e.minW) {
		return
	}
	if height != KeepValue && (height > 200 || height < e.minH) {
		return
	}

	if width != KeepValue {
		e.width.Set(types.AWidth(width))
	}

	e.height.Set(1)
}

// PasswordMode returns whether password mode is enabled for the control
func (e *TEditField) PasswordMode() bool {
	return e.showStars
}

// SetPasswordMode changes the way an EditField displays it content.
// If PasswordMode is false then the EditField works as regular text entry
// control. If PasswordMode is true then the EditField shows its content hidden
// with star characters ('*' by default)
func (e *TEditField) SetPasswordMode(pass bool) {
	e.showStars = pass
}
