package tv

import (
	"sync/atomic"
	"time"

	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/autoheight"
	"github.com/prospero78/goTV/tv/autowidth"
	"github.com/prospero78/goTV/tv/types"
)

/*
Button is a simpe push button control. Every time a user clicks a Button, it
emits OnClick event. Event has only one valid field Sender.
Button can be clicked with mouse or using space on keyboard while the Button is active.
*/
type Button struct {
	TBaseControl
	shadowColor term.Attribute
	pressed     int32
	shadowType  ButtonShadow
	onClick     func(Event)
	autoWidth   types.IAutoWidth
	autoHeight  types.IAutoHeight
}

/*
NewButton creates a new Button.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
title - button title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateButton(parent IControl,
	width, height int,
	title string,
	scale int,
	autoWidth types.AAutoWidth,
	autoHight types.AAutoHeight) *Button {
	b := &Button{
		TBaseControl: NewBaseControl(),
		autoWidth:    autowidth.New(),
		autoHeight:   autoheight.New(),
	}

	b.parent = parent
	b.align = AlignCenter
	b.autoWidth.Change(autoWidth)
	b.autoHeight.Change(autoHight)
	if b.autoHeight.Is() {
		height = 4
	}
	if b.autoWidth.Is() {
		width = xs.Len(title) + 2 + 1
	}

	if height < 4 {
		height = 4
	}
	if width < 6 {
		width = 6
	}

	b.SetTitle(title)
	b.SetSize(width, height)
	b.SetConstraints(width, height)
	b.SetScale(scale)

	if parent != nil {
		parent.AddChild(b)
	}

	return b
}

// Repaint draws the control on its View surface
func (b *Button) Draw() {
	if b.hidden {
		return
	}

	b.mtx.RLock()
	defer b.mtx.RUnlock()
	PushAttributes()
	defer PopAttributes()

	x, y := b.pos.Get()
	w, h := b.Size()

	fg, bg := b.fg, b.bg
	shadow := RealColor(b.shadowColor, b.Style(), ColorButtonShadow)
	switch {
	case b.disabled:
		fg, bg = RealColor(fg, b.Style(), ColorButtonDisabledText), RealColor(bg, b.Style(), ColorButtonDisabledBack)
	case b.Active():
		fg, bg = RealColor(b.fgActive, b.Style(), ColorButtonActiveText), RealColor(b.bgActive, b.Style(), ColorButtonActiveBack)
	default:
		fg, bg = RealColor(fg, b.Style(), ColorButtonText), RealColor(bg, b.Style(), ColorButtonBack)
	}

	dy := (h - 1) / 2
	SetTextColor(fg)
	shift, text := AlignColorizedText(b.title, w-1, b.align)
	if b.isPressed() == 0 {
		switch b.shadowType {
		case ShadowFull:
			SetBackColor(shadow)
			FillRect(x+1, y+types.ACoordY(h-1), w-1, 1, ' ')
			FillRect(x+types.ACoordX(w-1), y+1, 1, h-1, ' ')
		case ShadowHalf:
			parts := []rune(SysObject(ObjButton))
			var bottomCh, rightCh rune
			if len(parts) < 2 {
				bottomCh, rightCh = '▀', '█'
			} else {
				bottomCh, rightCh = parts[0], parts[1]
			}
			SetTextColor(shadow)
			FillRect(x+1, y+types.ACoordY(h-1), w-1, 1, bottomCh)
			FillRect(x+types.ACoordX(w-1), y+1, 1, h-2, rightCh)
		}
		SetTextColor(fg)
		SetBackColor(bg)
		FillRect(x, y, w-1, h-1, ' ')
		DrawText(x+types.ACoordX(shift), y+types.ACoordY(dy), text)
	} else {
		SetBackColor(bg)
		FillRect(x+1, y+1, w-1, h-1, ' ')
		DrawText(x+types.ACoordX(1+shift), y+types.ACoordY(1+dy), b.title)
	}
}

func (b *Button) isPressed() int32 {
	return atomic.LoadInt32(&b.pressed)
}

func (b *Button) setPressed(pressed int32) {
	atomic.StoreInt32(&b.pressed, pressed)
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (b *Button) ProcessEvent(event Event) bool {
	if !b.Enabled() {
		return false
	}

	if event.Type == EventKey {
		if event.Key == term.KeySpace && b.isPressed() == 0 {
			b.setPressed(1)
			ev := Event{Type: EventRedraw}

			go func() {
				PutEvent(ev)
				time.Sleep(100 * time.Millisecond)
				b.setPressed(0)
				PutEvent(ev)
			}()

			if b.onClick != nil {
				b.onClick(event)
			}
			return true
		} else if event.Key == term.KeyEsc && b.isPressed() != 0 {
			b.setPressed(0)
			ReleaseEvents()
			return true
		}
	} else if event.Type == EventMouse {
		if event.Key == term.MouseLeft {
			b.setPressed(1)
			GrabEvents(b)
			return true
		} else if event.Key == term.MouseRelease && b.isPressed() != 0 {
			ReleaseEvents()
			if event.X >= b.pos.GetX() &&
				event.Y >= b.pos.GetY() &&
				event.X < b.pos.GetX()+types.ACoordX(b.width.Get()) &&
				event.Y < b.pos.GetY()+types.ACoordY(b.height.Get()) {
				if b.onClick != nil {
					b.onClick(event)
				}
			}
			b.setPressed(0)
			return true
		}
	}

	return false
}

// OnClick sets the callback that is called when one clicks button
// with mouse or pressing space on keyboard while the button is active
func (b *Button) OnClick(fn func(Event)) {
	b.onClick = fn
}

// ShadowType returns type of a show the button drops
func (b *Button) ShadowType() ButtonShadow {
	return b.shadowType
}

// SetShadowType changes the shadow the button drops
func (b *Button) SetShadowType(sh ButtonShadow) {
	b.mtx.Lock()
	b.shadowType = sh
	b.mtx.Unlock()
}
