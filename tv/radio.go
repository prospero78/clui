package tv

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/widgets/event"
	"github.com/prospero78/goTV/tv/widgets/widgetbase"
)

/*
Radio button control. Unite a few radios in one radio group to
make a user select one of available choices.
*/
type Radio struct {
	widgetbase.TWidgetBase
	group    *RadioGroup
	selected bool

	onChange func(bool)
}

/*
CreateRadio creates a new radio button.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width - is minimal width of the control.
title - radio title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateRadio(parent types.IWidget, width int, title string, scale int) *Radio {
	c := new(Radio)
	c.TWidgetBase = widgetbase.New()

	if width == cons.AutoSize {
		width = xs.Len(title) + 4
	}

	c.parent = parent

	c.SetSize(width, 1) // TODO: only one line heigth is supported at that moment
	c.SetConstraints(width, 1)
	c.SetTitle(title)
	c.SetTabStop(true)
	c.SetScale(scale)

	c.onChange = nil

	if parent != nil {
		parent.AddChild(c)
	}

	return c
}

// Draw repaints the control on its View surface
func (c *Radio) Draw() {
	if c.isHidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	x, y := c.Pos()
	w, h := c.Size()

	fg, bg := RealColor(c.fg, c.Style(), cons.ColorControlText), RealColor(c.bg, c.Style(), cons.ColorControlBack)
	if !c.Enabled() {
		fg, bg = RealColor(c.fg, c.Style(), cons.ColorControlDisabledText), RealColor(c.bg, c.Style(), cons.ColorControlDisabledBack)
	} else if c.Active() {
		fg, bg = RealColor(c.fg, c.Style(), cons.ColorControlActiveText), RealColor(c.bg, c.Style(), cons.ColorControlActiveBack)
	}

	parts := []rune(SysObject(cons.ObjRadio))
	cOpen, cClose, cEmpty, cCheck := parts[0], parts[1], parts[2], parts[3]
	cState := cEmpty
	if c.selected {
		cState = cCheck
	}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, h, ' ')
	if w < 3 {
		return
	}

	PutChar(x, y, cOpen)
	PutChar(x+2, y, cClose)
	PutChar(x+1, y, cState)

	if w < 5 {
		return
	}

	shift, text := AlignColorizedText(c.title, w-4, c.align)
	DrawText(x+4+shift, y, text)
}

// ProcessEvent processes all events come from the control parent. If a control
// processes an event it should return true. If the method returns false it means
// that the control do not want or cannot process the event and the caller sends
// the event to the control parent.
// The control processes only space button and mouse clicks to make control selected.
// Deselecting control is not possible: one has to click another radio of the radio
// group to deselect this button
func (c *Radio) ProcessEvent(_event event.TEvent) bool {
	if (!c.Active() && _event.Type == cons.EventKey) || !c.Enabled() {
		return false
	}

	if (_event.Type == cons.EventKey && _event.Key == term.KeySpace) || _event.Type == cons.EventClick {
		if c.group == nil {
			c.SetSelected(true)
		} else {
			c.group.SelectItem(c)
		}
		return true
	}

	return false
}

// SetSelected makes the button selected. One should not use
// the method directly, it is for RadioGroup control
func (c *Radio) SetSelected(val bool) {
	c.selected = val

	if c.onChange != nil {
		go c.onChange(val)
	}
}

// Selected returns if the radio is selected
func (c *Radio) Selected() bool {
	return c.selected
}

// SetGroup sets the radio group to which the radio belongs
func (c *Radio) SetGroup(group *RadioGroup) {
	c.group = group
}

// OnChange sets the callback that is called whenever the state
// of the Radio is changed. Argument of callback is the current
func (c *Radio) OnChange(fn func(bool)) {
	c.block.Lock()
	defer c.block.Unlock()

	c.onChange = fn
}
