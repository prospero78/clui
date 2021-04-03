package tv

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"github.com/prospero78/goTV/tv/autosize"
	"github.com/prospero78/goTV/tv/pos"
	"github.com/prospero78/goTV/tv/types"
)

// Window is an implementation of View managed by Composer.
type Window struct {
	TBaseControl
	origWidth  int
	origHeight int
	posOrig    types.IPos
	buttons    ViewButton
	maximized  bool
	// maximization support

	hidden    bool
	immovable bool
	fixedSize bool
	border    BorderStyle

	onClose        func(Event) bool
	onScreenResize func(Event)

	onKeyDown  *keyDownCb
	autoWidth  types.IAutoSize
	autoHeight types.IAutoSize
}

type keyDownCb struct {
	data interface{}
	fn   func(evt Event, data interface{}) bool
}

func NewWindow(x types.ACoordX, y types.ACoordY,
	w, h int,
	title string,
	autoWidth, autoHight types.AAutoSize) *Window {
	wnd := &Window{
		posOrig:    pos.New(),
		autoWidth:  autosize.New(),
		autoHeight: autosize.New(),
	}
	wnd.autoWidth.Change(autoWidth)
	wnd.autoHeight.Change(autoHight)
	wnd.TBaseControl = NewBaseControl()

	if wnd.autoWidth.Is() || w < 1 || w > 1000 {
		w = 10
	}
	if wnd.autoHeight.Is() || h < 1 || h > 1000 {
		w = 5
	}

	wnd.SetConstraints(w, h)
	wnd.SetSize(w, h)
	wnd.SetPos(x, y)
	wnd.SetTitle(title)
	wnd.buttons = ButtonClose | ButtonBottom | ButtonMaximize
	wnd.children = make([]Control, 0)
	wnd.SetPaddings(1, 1)
	wnd.SetGaps(1, 0)
	wnd.SetScale(1)
	wnd.SetBorder(BorderAuto)

	return wnd
}

func (wnd *Window) buttonCount() (left, right int) {
	if wnd.buttons&ButtonClose == ButtonClose {
		right += 1
	}
	if wnd.buttons&ButtonMaximize == ButtonMaximize {
		left += 1
	}
	if wnd.buttons&ButtonBottom == ButtonBottom {
		left += 1
	}

	return left, right
}

func (wnd *Window) drawFrame() {
	PushAttributes()
	defer PopAttributes()

	var bs BorderStyle
	switch {
	case wnd.border == BorderAuto:
		if wnd.inactive {
			bs = BorderThin
		} else {
			bs = BorderThick
		}
	case wnd.border == BorderNone:
	default:
		bs = wnd.border
	}

	DrawFrame(wnd.pos.GetX(), wnd.pos.GetY(), int(wnd.width), wnd.height, bs)
}

func (wnd *Window) drawTitle() {
	PushAttributes()
	defer PopAttributes()

	lb, rb := wnd.buttonCount()
	maxw := wnd.width - 2
	xshift := 1
	if lb > 0 {
		lbSize := lb + 2 + 1
		maxw -= types.AWidth(lbSize)
		xshift += lbSize
	}
	if rb > 0 {
		maxw -= types.AWidth(rb) + 2
	}
	if maxw < 3 {
		return
	}

	fitTitle := wnd.title
	rawText := UnColorizeText(fitTitle)
	if xs.Len(rawText) > int(maxw) {
		fitTitle = SliceColorized(fitTitle, 0, int(maxw)-3) + "..."
	}
	DrawText(wnd.pos.GetX()+types.ACoordX(xshift), wnd.pos.GetY(), fitTitle)
}

func (wnd *Window) drawButtons() {
	lb, rb := wnd.buttonCount()
	if lb+rb == 0 {
		return
	}

	PushAttributes()
	defer PopAttributes()

	chars := []rune(SysObject(ObjViewButtons))
	cMax, cBottom, cClose, cOpenB, cCloseB := chars[0], chars[1], chars[2], chars[3], chars[4]

	// draw close button (rb can be either 1 or 0)
	if rb != 0 {
		pos := wnd.pos.GetX() + types.ACoordX(int(wnd.width)-rb-2)
		putCharUnsafe(pos, wnd.pos.GetY(), cOpenB)
		putCharUnsafe(pos+1, wnd.pos.GetY(), cClose)
		putCharUnsafe(pos+2, wnd.pos.GetY(), cCloseB)
	}

	if lb > 0 {
		posX := wnd.pos.GetX() + 1
		putCharUnsafe(posX, wnd.pos.GetY(), cOpenB)
		posX += 1
		if wnd.buttons&ButtonBottom == ButtonBottom {
			putCharUnsafe(posX, wnd.pos.GetY(), cBottom)
			posX += 1
		}
		if wnd.buttons&ButtonMaximize == ButtonMaximize {
			putCharUnsafe(posX, wnd.pos.GetY(), cMax)
			posX += 1
		}
		putCharUnsafe(posX, wnd.pos.GetY(), cCloseB)
	}
}

// Draw repaints the control on the screen
func (wnd *Window) Draw() {
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()
	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(wnd.fg, wnd.Style(), ColorViewText), RealColor(wnd.bg, wnd.Style(), ColorViewBack)
	SetBackColor(bg)

	FillRect(wnd.pos.GetX(), wnd.pos.GetY(), int(wnd.width), wnd.height, ' ')

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
func (c *Window) HitTest(x types.ACoordX, y types.ACoordY) HitResult {
	if x > c.pos.GetX() && x < c.pos.GetX()+types.ACoordX(c.width-1) &&
		y > c.pos.GetY() && y < c.pos.GetY()+types.ACoordY(c.height-1) {
		return HitInside
	}

	hResult := HitOutside

	switch {
	case x == c.pos.GetX() && y == c.pos.GetY():
		hResult = HitTopLeft
	case x == c.pos.GetX()+types.ACoordX(c.width-1) && y == c.pos.GetY():
		hResult = HitTopRight
	case x == c.pos.GetX() && y == c.pos.GetY()+types.ACoordY(c.height-1):
		hResult = HitBottomLeft
	case x == c.pos.GetX()+types.ACoordX(c.width-1) && y == c.pos.GetY()+types.ACoordY(c.height-1):
		hResult = HitBottomRight
	case x == c.pos.GetX() && y > c.pos.GetY() && y < c.pos.GetY()+types.ACoordY(c.height-1):
		hResult = HitLeft
	case x == c.pos.GetX()+types.ACoordX(c.width-1) && y > c.pos.GetY() && y < c.pos.GetY()+types.ACoordY(c.height-1):
		hResult = HitRight
	case y == c.pos.GetY() && x > c.pos.GetX() && x < c.pos.GetX()+types.ACoordX(c.width-1):
		lb, rb := c.buttonCount()
		fromL, fromR := lb, rb
		if lb > 0 {
			fromL += 2
		}
		if rb > 0 {
			fromR += 2
		}
		if x > c.pos.GetX()+types.ACoordX(fromL) && x < c.pos.GetX()+types.ACoordX(int(c.width)-fromR) {
			hResult = HitTop
		} else {
			hResult = HitTop
			if c.buttons&ButtonClose == ButtonClose && rb != 0 && x == c.pos.GetX()+types.ACoordX(c.width-2) {
				hResult = HitButtonClose
			} else if lb != 0 && x > c.pos.GetX()+1 && x < c.pos.GetX()+types.ACoordX(2+lb) {
				dx := x - c.pos.GetX() - 2
				hitRes := []HitResult{HitTop, HitTop}
				pos := 0
				if c.buttons&ButtonBottom == ButtonBottom {
					hitRes[pos] = HitButtonBottom
					pos += 1
				}
				if c.buttons&ButtonMaximize == ButtonMaximize {
					hitRes[pos] = HitButtonMaximize
					// pos += 1
				}
				if int(dx) < len(hitRes) {
					hResult = hitRes[dx]
				}
			}
		}
	case y == c.pos.GetY()+types.ACoordY(c.height-1) && x > c.pos.GetX() && x < c.pos.GetX()+types.ACoordX(c.width-1):
		hResult = HitBottom
	}

	if hResult != HitOutside {
		if c.immovable && hResult == HitTop {
			hResult = HitInside
		}
		if c.fixedSize &&
			(hResult == HitBottom || hResult == HitLeft ||
				hResult == HitRight || hResult == HitTopLeft ||
				hResult == HitTopRight || hResult == HitBottomRight ||
				hResult == HitBottomLeft || hResult == HitButtonMaximize) {
			hResult = HitInside
		}
	}

	return hResult
}

func (c *Window) ProcessEvent(ev Event) bool {
	switch ev.Type {
	case EventMove:
		c.PlaceChildren()
	case EventResize:
		c.ResizeChildren()
		c.PlaceChildren()
	case EventClose:
		if c.onClose != nil {
			if !c.onClose(ev) {
				return false
			}
		}
		return true
	case EventKey:
		if ev.Key == term.KeyTab || ev.Key == term.KeyArrowUp || ev.Key == term.KeyArrowDown {
			if SendEventToChild(c, ev) {
				return true
			}

			aC := ActiveControl(c)
			nC := NextControl(c, aC, ev.Key != term.KeyArrowUp)

			var clipped Control

			if aC != nil && aC.Clipped() {
				clipped = aC
			} else if nC != nil {
				clipped = ClippedParent(nC)
			}

			if clipped != nil {
				dir := types.ACoordX(1)
				if ev.Key != term.KeyArrowUp {
					dir = -1
				}

				clipped.ProcessEvent(Event{Type: EventActivateChild, Target: nC, X: dir})
			}

			if nC != aC {
				if aC != nil {
					aC.SetActive(false)
					aC.ProcessEvent(Event{Type: EventActivate, X: 0})
				}
				if nC != nil {
					nC.SetActive(true)
					nC.ProcessEvent(Event{Type: EventActivate, X: 1})
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
		if ev.Type == EventMouse && ev.Key == term.MouseLeft {
			DeactivateControls(c)
		}
		return SendEventToChild(c, ev)
	}

	return false
}

// OnClose sets the callback that is called when the Window is about to destroy
func (w *Window) OnClose(fn func(Event) bool) {
	w.onClose = fn
}

// OnKeyDown sets the callback that is called when a user presses a key
// while the Window is active
func (w *Window) OnKeyDown(fn func(Event, interface{}) bool, data interface{}) {
	if fn == nil {
		w.onKeyDown = nil
	} else {
		w.onKeyDown = &keyDownCb{data: data, fn: fn}
	}
}

// OnScreenResize sets the callback that is called when size of terminal changes
func (w *Window) OnScreenResize(fn func(Event)) {
	w.onScreenResize = fn
}

// Border returns the default window border
func (w *Window) Border() BorderStyle {
	return w.border
}

// SetBorder changes the default window border
func (w *Window) SetBorder(border BorderStyle) {
	w.border = border
}

// SetMaximized opens the view to full screen or restores its
// previous size
func (w *Window) SetMaximized(maximize bool) {
	if maximize == w.maximized {
		return
	}

	if maximize {
		x, y := w.pos.Get()
		w.posOrig.X().Set(x)
		w.posOrig.Y().Set(y)
		w.origWidth, w.origHeight = w.Size()
		w.maximized = true
		w.SetPos(0, 0)
		width, height := ScreenSize()
		w.SetSize(width, height)
	} else {
		w.maximized = false
		w.SetPos(w.posOrig.GetX(), w.posOrig.GetY())
		w.SetSize(w.origWidth, w.origHeight)
	}
	w.ResizeChildren()
	w.PlaceChildren()
}

// Maximized returns if the view is in full screen mode
func (w *Window) Maximized() bool {
	return w.maximized
}

// Visible returns if the window must be drawn on the screen
func (w *Window) Visible() bool {
	return !w.hidden
}

// SetVisible allows to temporarily remove the window from screen
// and show it later without reconstruction
func (w *Window) SetVisible(visible bool) {
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
func (w *Window) Movable() bool {
	return !w.immovable
}

// Sizable returns if size of the Window can be changed with mouse or keyboard
func (w *Window) Sizable() bool {
	return !w.fixedSize
}

// SetMovable turns on and off ability to change Window position with mouse
// or keyboard
func (w *Window) SetMovable(movable bool) {
	w.immovable = !movable
}

// SetSizable turns on and off ability to change Window size with mouse
// or keyboard
func (w *Window) SetSizable(sizable bool) {
	w.fixedSize = !sizable
}

// TitleButtons returns a set of buttons shown in the Window title bar
func (w *Window) TitleButtons() ViewButton {
	return w.buttons
}

// SetTitleButtons sets the title bar buttons available for a user
func (w *Window) SetTitleButtons(buttons ViewButton) {
	w.buttons = buttons
}
