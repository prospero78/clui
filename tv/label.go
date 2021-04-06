package tv

import (
	xs "github.com/huandu/xstrings"

	"github.com/prospero78/goTV/tv/autoheight"
	"github.com/prospero78/goTV/tv/autowidth"
	"github.com/prospero78/goTV/tv/types"
)

/*
Label is a decorative control that can display text in horizontal
or vertical direction. Other available text features are alignment
and multi-line ability. Text can be single- or multi-colored with
tags inside the text. Multi-colored strings have limited support
of alignment feature: if text is longer than Label width the text
is always left aligned
*/
type Label struct {
	TBaseControl
	direction   Direction
	multiline   bool
	textDisplay Align
	autoWidth   types.IAutoWidth
	autoHeight  types.IAutoHeight
}

/*
NewLabel creates a new label.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
w and h - are minimal size of the control.
title - is Label title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateLabel(parent IControl, w, h int, title string, scale int) *Label {
	c := &Label{
		TBaseControl: NewBaseControl(),
		autoWidth:    autowidth.New(),
		autoHeight:   autoheight.New(),
	}
	if w == 0 {
		c.autoWidth.Set()
		w = xs.Len(title)
	}
	if h == 0 {
		c.autoHeight.Set()
		h = 1
	}

	c.parent = parent

	c.SetTitle(title)
	c.SetSize(w, h)
	c.SetConstraints(w, h)
	c.SetScale(scale)
	c.tabSkip = true
	c.textDisplay = AlignLeft

	if parent != nil {
		parent.AddChild(c)
	}

	return c
}

// Direction returns direction of text output: vertical or horizontal
func (l *Label) Direction() Direction {
	return l.direction
}

// SetDirection sets the text output direction
func (l *Label) SetDirection(dir Direction) {
	l.direction = dir
}

func (l *Label) Draw() {
	if l.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(l.fg, l.Style(), ColorText), RealColor(l.bg, l.Style(), ColorBack)
	if !l.Enabled() {
		fg = RealColor(l.fg, l.Style(), ColorDisabledText)
	}

	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(l.pos.GetX(), l.pos.GetY(), int(l.width.Get()), int(l.height.Get()), ' ')

	if l.title == "" {
		return
	}

	if l.multiline {
		parser := NewColorParser(l.title, fg, bg)
		elem := parser.NextElement()
		xx, yy := l.pos.Get()
		for elem.Type != ElemEndOfText {
			if xx >= l.pos.GetX()+types.ACoordX(l.width.Get()) || yy >= l.pos.GetY()+types.ACoordY(l.height.Get()) {
				break
			}

			if elem.Type == ElemLineBreak {
				xx = l.pos.GetX()
				yy++
			} else if elem.Type == ElemPrintable {
				SetTextColor(elem.Fg)
				SetBackColor(elem.Bg)
				putCharUnsafe(xx, yy, elem.Ch)

				if l.direction == Horizontal {
					xx++
					if xx >= l.pos.GetX()+types.ACoordX(l.width.Get()) {
						xx = l.pos.GetX()
						yy++
					}
				} else {
					yy++
					if yy >= l.pos.GetY()+types.ACoordY(l.height.Get()) {
						yy = l.pos.GetY()
						xx++
					}
				}
			}

			elem = parser.NextElement()
		}
	} else {
		if l.direction == Horizontal {
			shift, str := AlignColorizedText(l.title, int(l.width.Get()), l.align)
			if str != l.title && l.align != l.textDisplay {
				shift, str = AlignColorizedText(l.title, int(l.width.Get()), l.textDisplay)
			}
			DrawText(l.pos.GetX()+types.ACoordX(shift), l.pos.GetY(), str)
		} else {
			shift, str := AlignColorizedText(l.title, int(l.height.Get()), l.align)
			if str != l.title && l.align != l.textDisplay {
				shift, str = AlignColorizedText(l.title, int(l.width.Get()), l.textDisplay)
			}
			DrawTextVertical(l.pos.GetX(), l.pos.GetY()+types.ACoordY(shift), str)
		}
	}
}

// Multiline returns if text is displayed on several lines if the
// label title is longer than label width or title contains
// line breaks
func (l *Label) Multiline() bool {
	return l.multiline
}

// SetMultiline sets if the label should output text as one line
// or automatically display it in several lines
func (l *Label) SetMultiline(multi bool) {
	l.multiline = multi
}

// TextDisplay returns which part of the lable title is displayed in case of
// title is longer than the label:
// - AlignLeft - the head of the title is shown
// - AlignRight - the tail of the title is shown
// The property is used only by single line Label
func (l *Label) TextDisplay() Align {
	return l.textDisplay
}

// SetTextDisplay sets which part of the title is displayed in case of the title
// is longer than the lable. Only AlignLeft and AlignRigth are valid values
// for the property. Any other value does is skipped and does not affect
// displaying the title
func (l *Label) SetTextDisplay(align Align) {
	if align != AlignLeft && align != AlignRight {
		return
	}

	l.textDisplay = align
}
