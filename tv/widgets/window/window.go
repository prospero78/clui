package window

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/widgets/event"
	"github.com/prospero78/goTV/tv/widgets/widgetbase"
)

// TWindow is an implementation of View managed by Composer.
type TWindow struct {
	widgetbase.TWidgetBase

	buttons   tv.ViewButton
	maximized bool
	// maximization support
	origWidth  int
	origHeight int
	origX      int
	origY      int
	hidden     bool
	immovable  bool
	fixedSize  bool
	border     types.ABorderStyle

	onClose        func(event.TEvent) bool
	onScreenResize func(event.TEvent)

	onKeyDown *keyDownCb
}

type keyDownCb struct {
	data interface{}
	fn   func(evt event.TEvent, data interface{}) bool
}

func CreateWindow(x, y, w, h int, title string) *TWindow {
	wnd := new(TWindow)
	wnd.TWidgetBase = widgetbase.New()

	if w == cons.AutoSize || w < 1 || w > 1000 {
		w = 10
	}
	if h == cons.AutoSize || h < 1 || h > 1000 {
		w = 5
	}

	wnd.SetConstraints(w, h)
	wnd.SetSize(w, h)
	wnd.SetPos(x, y)
	wnd.SetTitle(title)
	wnd.buttons = cons.ButtonClose | cons.ButtonBottom | cons.ButtonMaximize
	wnd.children = make([]types.IWidget, 0)
	wnd.SetPaddings(1, 1)
	wnd.SetGaps(1, 0)
	wnd.SetScale(1)
	wnd.SetBorder(cons.BorderAuto)

	return wnd
}

func (wnd *TWindow) buttonCount() (left, right int) {
	if wnd.buttons&cons.ButtonClose == cons.ButtonClose {
		right += 1
	}
	if wnd.buttons&cons.ButtonMaximize == cons.ButtonMaximize {
		left += 1
	}
	if wnd.buttons&cons.ButtonBottom == cons.ButtonBottom {
		left += 1
	}

	return left, right
}

func (wnd *TWindow) drawFrame() {
	PushAttributes()
	defer PopAttributes()

	var bs cons.BorderStyle
	switch {
	case wnd.border == cons.BorderAuto:
		if wnd.IsInactive() {
			bs = cons.BorderThin
		} else {
			bs = cons.BorderThick
		}
	case wnd.border == cons.BorderNone:
	default:
		bs = wnd.border
	}

	DrawFrame(wnd.x, wnd.y, wnd.width, wnd.height, bs)
}

func (wnd *TWindow) drawTitle() {
	PushAttributes()
	defer PopAttributes()

	lb, rb := wnd.buttonCount()
	maxw := wnd.width - 2
	xshift := 1
	if lb > 0 {
		lbSize := lb + 2 + 1
		maxw -= lbSize
		xshift += lbSize
	}
	if rb > 0 {
		maxw -= rb + 2
	}
	if maxw < 3 {
		return
	}

	fitTitle := wnd.title
	rawText := UnColorizeText(fitTitle)
	if xs.Len(rawText) > maxw {
		fitTitle = SliceColorized(fitTitle, 0, maxw-3) + "..."
	}
	DrawText(wnd.x+xshift, wnd.y, fitTitle)
}

func (wnd *TWindow) drawButtons() {
	lb, rb := wnd.buttonCount()
	if lb+rb == 0 {
		return
	}

	PushAttributes()
	defer PopAttributes()

	chars := []rune(SysObject(cons.ObjViewButtons))
	cMax, cBottom, cClose, cOpenB, cCloseB := chars[0], chars[1], chars[2], chars[3], chars[4]

	// draw close button (rb can be either 1 or 0)
	if rb != 0 {
		pos := wnd.x + wnd.width - rb - 2
		putCharUnsafe(pos, wnd.y, cOpenB)
		putCharUnsafe(pos+1, wnd.y, cClose)
		putCharUnsafe(pos+2, wnd.y, cCloseB)
	}

	if lb > 0 {
		pos := wnd.x + 1
		putCharUnsafe(pos, wnd.y, cOpenB)
		pos += 1
		if wnd.buttons&cons.ButtonBottom == cons.ButtonBottom {
			putCharUnsafe(pos, wnd.y, cBottom)
			pos += 1
		}
		if wnd.buttons&cons.ButtonMaximize == cons.ButtonMaximize {
			putCharUnsafe(pos, wnd.y, cMax)
			pos += 1
		}
		putCharUnsafe(pos, wnd.y, cCloseB)
	}
}

// Draw repaints the control on the screen
func (wnd *TWindow) Draw() {
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()
	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(wnd.fg, wnd.Style(), cons.ColorViewText), RealColor(wnd.bg, wnd.Style(), cons.ColorViewBack)
	SetBackColor(bg)

	FillRect(wnd.x, wnd.y, wnd.width, wnd.height, ' ')

	wnd.DrawChildren()

	SetBackColor(bg)
	SetTextColor(fg)

	wnd.drawFrame()
	wnd.drawTitle()
	wnd.drawButtons()
}

// HitTest returns type of a Window region at a given screen coordinates. The
// method is used to detect if a mouse cursor on a window border or outside,
// which window icon is under cursor etc
func (c *TWindow) HitTest(x, y int) types.AHitResult {
	if x > c.x && x < c.x+c.width-1 &&
		y > c.y && y < c.y+c.height-1 {
		return cons.HitInside
	}

	hResult := cons.HitOutside

	switch {
	case x == c.x && y == c.y:
		hResult = cons.HitTopLeft
	case x == c.x+c.width-1 && y == c.y:
		hResult = cons.HitTopRight
	case x == c.x && y == c.y+c.height-1:
		hResult = cons.HitBottomLeft
	case x == c.x+c.width-1 && y == c.y+c.height-1:
		hResult = cons.HitBottomRight
	case x == c.x && y > c.y && y < c.y+c.height-1:
		hResult = cons.HitLeft
	case x == c.x+c.width-1 && y > c.y && y < c.y+c.height-1:
		hResult = cons.HitRight
	case y == c.y && x > c.x && x < c.x+c.width-1:
		lb, rb := c.buttonCount()
		fromL, fromR := lb, rb
		if lb > 0 {
			fromL += 2
		}
		if rb > 0 {
			fromR += 2
		}
		if x > c.x+fromL && x < c.x+c.width-fromR {
			hResult = cons.HitTop
		} else {
			hResult = cons.HitTop
			if c.buttons&cons.ButtonClose == cons.ButtonClose && rb != 0 && x == c.x+c.width-2 {
				hResult = cons.HitButtonClose
			} else if lb != 0 && x > c.x+1 && x < c.x+2+lb {
				dx := x - c.x - 2
				hitRes := []cons.HitResult{cons.HitTop, cons.HitTop}
				pos := 0
				if c.buttons&cons.ButtonBottom == cons.ButtonBottom {
					hitRes[pos] = cons.HitButtonBottom
					pos += 1
				}
				if c.buttons&cons.ButtonMaximize == cons.ButtonMaximize {
					hitRes[pos] = cons.HitButtonMaximize
					pos += 1
				}
				if dx < len(hitRes) {
					hResult = hitRes[dx]
				}
			}
		}
	case y == c.y+c.height-1 && x > c.x && x < c.x+c.width-1:
		hResult = cons.HitBottom

	}

	if hResult != cons.HitOutside {
		if c.immovable && hResult == cons.HitTop {
			hResult = cons.HitInside
		}
		if c.fixedSize &&
			(hResult == cons.HitBottom || hResult == cons.HitLeft ||
				hResult == cons.HitRight || hResult == cons.HitTopLeft ||
				hResult == cons.HitTopRight || hResult == cons.HitBottomRight ||
				hResult == cons.HitBottomLeft || hResult == cons.HitButtonMaximize) {
			hResult = cons.HitInside
		}
	}

	return hResult
}

func (c *TWindow) ProcessEvent(ev event.TEvent) bool {
	switch ev.Type {
	case cons.EventMove:
		c.PlaceChildren()
	case cons.EventResize:
		c.ResizeChildren()
		c.PlaceChildren()
	case cons.EventClose:
		if c.onClose != nil {
			if !c.onClose(ev) {
				return false
			}
		}
		return true
	case cons.EventKey:
		if ev.Key == term.KeyTab || ev.Key == term.KeyArrowUp || ev.Key == term.KeyArrowDown {
			if SendEventToChild(c, ev) {
				return true
			}

			aC := ActiveControl(c)
			nC := NextControl(c, aC, ev.Key != term.KeyArrowUp)

			var clipped types.IWidget

			if aC != nil && aC.Clipped() {
				clipped = aC
			} else if nC != nil {
				clipped = ClippedParent(nC)
			}

			if clipped != nil {
				dir := 1
				if ev.Key != term.KeyArrowUp {
					dir = -1
				}

				clipped.ProcessEvent(event.TEvent{Type: cons.EventActivateChild, Target: nC, X: dir})
			}

			if nC != aC {
				if aC != nil {
					aC.SetActive(false)
					aC.ProcessEvent(event.TEvent{Type: cons.EventActivate, X: 0})
				}
				if nC != nil {
					nC.SetActive(true)
					nC.ProcessEvent(event.TEvent{Type: cons.EventActivate, X: 1})
				}
			}
			return true
		}
		if SendEventToChild(c, ev) {
			return true
		}
		if c.onKeyDown != nil {
			return c.onKeyDown.fn(ev, c.onKeyDown.data)
		}
		return false
	default:
		if ev.Type == cons.EventMouse && ev.Key == term.MouseLeft {
			DeactivateControls(c)
		}
		return SendEventToChild(c, ev)
	}

	return false
}

// OnClose sets the callback that is called when the Window is about to destroy
func (w *TWindow) OnClose(fn func(event.TEvent) bool) {
	w.onClose = fn
}

// OnKeyDown sets the callback that is called when a user presses a key
// while the Window is active
func (w *TWindow) OnKeyDown(fn func(event.TEvent, interface{}) bool, data interface{}) {
	if fn == nil {
		w.onKeyDown = nil
	} else {
		w.onKeyDown = &keyDownCb{data: data, fn: fn}
	}
}

// OnScreenResize sets the callback that is called when size of terminal changes
func (w *TWindow) OnScreenResize(fn func(event.TEvent)) {
	w.onScreenResize = fn
}

// Border returns the default window border
func (w *TWindow) Border() cons.BorderStyle {
	return w.border
}

// SetBorder changes the default window border
func (w *TWindow) SetBorder(border cons.BorderStyle) {
	w.border = border
}

// SetMaximized opens the view to full screen or restores its
// previous size
func (w *TWindow) SetMaximized(maximize bool) {
	if maximize == w.maximized {
		return
	}

	if maximize {
		w.origX, w.origY = w.Pos()
		w.origWidth, w.origHeight = w.Size()
		w.maximized = true
		w.SetPos(0, 0)
		width, height := ScreenSize()
		w.SetSize(width, height)
	} else {
		w.maximized = false
		w.SetPos(w.origX, w.origY)
		w.SetSize(w.origWidth, w.origHeight)
	}
	w.ResizeChildren()
	w.PlaceChildren()
}

// Maximized returns if the view is in full screen mode
func (w *TWindow) Maximized() bool {
	return w.maximized
}

// Visible returns if the window must be drawn on the screen
func (w *TWindow) Visible() bool {
	return !w.hidden
}

// SetVisible allows to temporarily remove the window from screen
// and show it later without reconstruction
func (w *TWindow) SetVisible(visible bool) {
	if w.hidden != visible {
		return
	}

	w.hidden = !visible
	if w.hidden {
		w.SetModal(false)
		if WindowManager().topWindow() == w {
			WindowManager().moveActiveWindowToBottom()
		}
	} else {
		WindowManager().activateWindow(w)
	}
}

// Movable returns if the Window can be moved with mouse or keyboard
func (w *TWindow) Movable() bool {
	return !w.immovable
}

// Sizable returns if size of the Window can be changed with mouse or keyboard
func (w *TWindow) Sizable() bool {
	return !w.fixedSize
}

// SetMovable turns on and off ability to change Window position with mouse
// or keyboard
func (w *TWindow) SetMovable(movable bool) {
	w.immovable = !movable
}

// SetSizable turns on and off ability to change Window size with mouse
// or keyboard
func (w *TWindow) SetSizable(sizable bool) {
	w.fixedSize = !sizable
}

// TitleButtons returns a set of buttons shown in the Window title bar
func (w *TWindow) TitleButtons() ViewButton {
	return w.buttons
}

// SetTitleButtons sets the title bar buttons available for a user
func (w *TWindow) SetTitleButtons(buttons ViewButton) {
	w.buttons = buttons
}
