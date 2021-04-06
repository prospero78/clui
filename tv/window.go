package tv

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/autoheight"
	"github.com/prospero78/goTV/tv/autowidth"
	"github.com/prospero78/goTV/tv/pos"
	"github.com/prospero78/goTV/tv/types"
)

// TWindow is an implementation of View managed by Composer.
type TWindow struct {
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
	autoWidth  types.IAutoWidth
	autoHeight types.IAutoHeight
}

type keyDownCb struct {
	data interface{}
	fn   func(evt Event, data interface{}) bool
}

func NewWindow(x types.ACoordX, y types.ACoordY,
	w, h int,
	title string,
	autoWidth types.AAutoWidth,
	autoHeight types.AAutoHeight) *TWindow {
	wnd := &TWindow{
		posOrig:    pos.New(),
		autoWidth:  autowidth.New(),
		autoHeight: autoheight.New(),
	}
	wnd.autoWidth.Change(autoWidth)
	wnd.autoHeight.Change(autoHeight)
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
	wnd.children = make([]IControl, 0)
	wnd.SetPaddings(1, 1)
	wnd.SetGaps(1, 0)
	wnd.SetScale(1)
	wnd.SetBorder(BorderAuto)

	return wnd
}

func (sf *TWindow) buttonCount() (left, right int) {
	if sf.buttons&ButtonClose == ButtonClose {
		right += 1
	}
	if sf.buttons&ButtonMaximize == ButtonMaximize {
		left += 1
	}
	if sf.buttons&ButtonBottom == ButtonBottom {
		left += 1
	}

	return left, right
}

func (sf *TWindow) drawFrame() {
	PushAttributes()
	defer PopAttributes()

	var bs BorderStyle
	switch {
	case sf.border == BorderAuto:
		if sf.inactive {
			bs = BorderThin
		} else {
			bs = BorderThick
		}
	case sf.border == BorderNone:
	default:
		bs = sf.border
	}

	DrawFrame(sf.pos.GetX(), sf.pos.GetY(), int(sf.width.Get()), int(sf.height), bs)
}

func (sf *TWindow) drawTitle() {
	PushAttributes()
	defer PopAttributes()

	lb, rb := sf.buttonCount()
	maxw := sf.width.Get() - 2
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

	fitTitle := sf.title
	rawText := UnColorizeText(fitTitle)
	if xs.Len(rawText) > int(maxw) {
		fitTitle = SliceColorized(fitTitle, 0, int(maxw)-3) + "..."
	}
	DrawText(sf.pos.GetX()+types.ACoordX(xshift), sf.pos.GetY(), fitTitle)
}

func (sf *TWindow) drawButtons() {
	lb, rb := sf.buttonCount()
	if lb+rb == 0 {
		return
	}

	PushAttributes()
	defer PopAttributes()

	chars := []rune(SysObject(ObjViewButtons))
	cMax, cBottom, cClose, cOpenB, cCloseB := chars[0], chars[1], chars[2], chars[3], chars[4]

	// draw close button (rb can be either 1 or 0)
	if rb != 0 {
		pos := sf.pos.GetX() + types.ACoordX(int(sf.width.Get())-rb-2)
		putCharUnsafe(pos, sf.pos.GetY(), cOpenB)
		putCharUnsafe(pos+1, sf.pos.GetY(), cClose)
		putCharUnsafe(pos+2, sf.pos.GetY(), cCloseB)
	}

	if lb > 0 {
		posX := sf.pos.GetX() + 1
		putCharUnsafe(posX, sf.pos.GetY(), cOpenB)
		posX += 1
		if sf.buttons&ButtonBottom == ButtonBottom {
			putCharUnsafe(posX, sf.pos.GetY(), cBottom)
			posX += 1
		}
		if sf.buttons&ButtonMaximize == ButtonMaximize {
			putCharUnsafe(posX, sf.pos.GetY(), cMax)
			posX += 1
		}
		putCharUnsafe(posX, sf.pos.GetY(), cCloseB)
	}
}

// Draw repaints the control on the screen
func (sf *TWindow) Draw() {
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()
	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(sf.fg, sf.Style(), ColorViewText), RealColor(sf.bg, sf.Style(), ColorViewBack)
	SetBackColor(bg)

	FillRect(sf.pos.GetX(), sf.pos.GetY(), int(sf.width.Get()), int(sf.height), ' ')

	sf.DrawChildren()

	SetBackColor(bg)
	SetTextColor(fg)

	sf.drawFrame()
	sf.drawTitle()
	sf.drawButtons()
}

// HitTest returns type of a Window region at a given screen coordinates. The
// method is used to detect if a mouse cursor on a window border or outside,
// which window icon is under cursor etc
func (sf *TWindow) HitTest(x types.ACoordX, y types.ACoordY) HitResult {
	if x > sf.pos.GetX() && x < sf.pos.GetX()+types.ACoordX(sf.width.Get()-1) &&
		y > sf.pos.GetY() && y < sf.pos.GetY()+types.ACoordY(sf.height-1) {
		return HitInside
	}

	hResult := HitOutside

	switch {
	case x == sf.pos.GetX() && y == sf.pos.GetY():
		hResult = HitTopLeft
	case x == sf.pos.GetX()+types.ACoordX(sf.width.Get()-1) && y == sf.pos.GetY():
		hResult = HitTopRight
	case x == sf.pos.GetX() && y == sf.pos.GetY()+types.ACoordY(sf.height-1):
		hResult = HitBottomLeft
	case x == sf.pos.GetX()+types.ACoordX(sf.width.Get()-1) && y == sf.pos.GetY()+types.ACoordY(sf.height-1):
		hResult = HitBottomRight
	case x == sf.pos.GetX() && y > sf.pos.GetY() && y < sf.pos.GetY()+types.ACoordY(sf.height-1):
		hResult = HitLeft
	case x == sf.pos.GetX()+types.ACoordX(sf.width.Get()-1) && y > sf.pos.GetY() && y < sf.pos.GetY()+types.ACoordY(sf.height-1):
		hResult = HitRight
	case y == sf.pos.GetY() && x > sf.pos.GetX() && x < sf.pos.GetX()+types.ACoordX(sf.width.Get()-1):
		lb, rb := sf.buttonCount()
		fromL, fromR := lb, rb
		if lb > 0 {
			fromL += 2
		}
		if rb > 0 {
			fromR += 2
		}
		if x > sf.pos.GetX()+types.ACoordX(fromL) && x < sf.pos.GetX()+types.ACoordX(int(sf.width.Get())-fromR) {
			hResult = HitTop
		} else {
			hResult = HitTop
			if sf.buttons&ButtonClose == ButtonClose && rb != 0 && x == sf.pos.GetX()+types.ACoordX(sf.width.Get()-2) {
				hResult = HitButtonClose
			} else if lb != 0 && x > sf.pos.GetX()+1 && x < sf.pos.GetX()+types.ACoordX(2+lb) {
				dx := x - sf.pos.GetX() - 2
				hitRes := []HitResult{HitTop, HitTop}
				pos := 0
				if sf.buttons&ButtonBottom == ButtonBottom {
					hitRes[pos] = HitButtonBottom
					pos += 1
				}
				if sf.buttons&ButtonMaximize == ButtonMaximize {
					hitRes[pos] = HitButtonMaximize
					// pos += 1
				}
				if int(dx) < len(hitRes) {
					hResult = hitRes[dx]
				}
			}
		}
	case y == sf.pos.GetY()+types.ACoordY(sf.height-1) &&
		x > sf.pos.GetX() &&
		x < sf.pos.GetX()+types.ACoordX(sf.width.Get()-1):
		hResult = HitBottom
	}

	if hResult != HitOutside {
		if sf.immovable && hResult == HitTop {
			hResult = HitInside
		}
		if sf.fixedSize &&
			(hResult == HitBottom || hResult == HitLeft ||
				hResult == HitRight || hResult == HitTopLeft ||
				hResult == HitTopRight || hResult == HitBottomRight ||
				hResult == HitBottomLeft || hResult == HitButtonMaximize) {
			hResult = HitInside
		}
	}

	return hResult
}

func (sf *TWindow) ProcessEvent(ev Event) bool {
	switch ev.Type {
	case EventMove:
		sf.PlaceChildren()
	case EventResize:
		sf.ResizeChildren()
		sf.PlaceChildren()
	case EventClose:
		if sf.onClose != nil {
			if !sf.onClose(ev) {
				return false
			}
		}
		return true
	case EventKey:
		if ev.Key == term.KeyTab || ev.Key == term.KeyArrowUp || ev.Key == term.KeyArrowDown {
			if SendEventToChild(sf, ev) {
				return true
			}

			aC := ActiveControl(sf)
			nC := NextControl(sf, aC, ev.Key != term.KeyArrowUp)

			var clipped IControl

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
		if SendEventToChild(sf, ev) {
			return true
		}
		if sf.onKeyDown != nil {
			return sf.onKeyDown.fn(ev, sf.onKeyDown.data)
		}
		return false
	default:
		if ev.Type == EventMouse && ev.Key == term.MouseLeft {
			DeactivateControls(sf)
		}
		return SendEventToChild(sf, ev)
	}

	return false
}

// OnClose sets the callback that is called when the Window is about to destroy
func (sf *TWindow) OnClose(fn func(Event) bool) {
	sf.onClose = fn
}

// OnKeyDown sets the callback that is called when a user presses a key
// while the Window is active
func (sf *TWindow) OnKeyDown(fn func(Event, interface{}) bool, data interface{}) {
	if fn == nil {
		sf.onKeyDown = nil
	} else {
		sf.onKeyDown = &keyDownCb{data: data, fn: fn}
	}
}

// OnScreenResize sets the callback that is called when size of terminal changes
func (sf *TWindow) OnScreenResize(fn func(Event)) {
	sf.onScreenResize = fn
}

// Border returns the default window border
func (sf *TWindow) Border() BorderStyle {
	return sf.border
}

// SetBorder changes the default window border
func (sf *TWindow) SetBorder(border BorderStyle) {
	sf.border = border
}

// SetMaximized opens the view to full screen or restores its
// previous size
func (sf *TWindow) SetMaximized(maximize bool) {
	if maximize == sf.maximized {
		return
	}

	if maximize {
		x, y := sf.pos.Get()
		sf.posOrig.X().Set(x)
		sf.posOrig.Y().Set(y)
		sf.origWidth, sf.origHeight = sf.Size()
		sf.maximized = true
		sf.SetPos(0, 0)
		width, height := ScreenSize()
		sf.SetSize(width, height)
	} else {
		sf.maximized = false
		sf.SetPos(sf.posOrig.GetX(), sf.posOrig.GetY())
		sf.SetSize(sf.origWidth, sf.origHeight)
	}
	sf.ResizeChildren()
	sf.PlaceChildren()
}

// Maximized returns if the view is in full screen mode
func (sf *TWindow) Maximized() bool {
	return sf.maximized
}

// Visible returns if the window must be drawn on the screen
func (sf *TWindow) Visible() bool {
	return !sf.hidden
}

// SetVisible allows to temporarily remove the window from screen
// and show it later without reconstruction
func (sf *TWindow) SetVisible(visible bool) {
	if sf.hidden != visible {
		return
	}

	sf.hidden = !visible
	if sf.hidden {
		sf.SetModal(false)
		if WindowManager().topWindow() == sf {
			WindowManager().moveActiveWindowToBottom()
		}
	} else {
		WindowManager().activateWindow(sf)
	}
}

// Movable returns if the Window can be moved with mouse or keyboard
func (sf *TWindow) Movable() bool {
	return !sf.immovable
}

// Sizable returns if size of the Window can be changed with mouse or keyboard
func (sf *TWindow) Sizable() bool {
	return !sf.fixedSize
}

// SetMovable turns on and off ability to change Window position with mouse
// or keyboard
func (sf *TWindow) SetMovable(movable bool) {
	sf.immovable = !movable
}

// SetSizable turns on and off ability to change Window size with mouse
// or keyboard
func (sf *TWindow) SetSizable(sizable bool) {
	sf.fixedSize = !sizable
}

// TitleButtons returns a set of buttons shown in the Window title bar
func (sf *TWindow) TitleButtons() ViewButton {
	return sf.buttons
}

// SetTitleButtons sets the title bar buttons available for a user
func (sf *TWindow) SetTitleButtons(buttons ViewButton) {
	sf.buttons = buttons
}
