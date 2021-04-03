package tv

import (
	"sync"
	"sync/atomic"

	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/pos"
	"github.com/prospero78/goTV/tv/types"
)

// TBaseControl is a base for all visible controls.
// Every new control must inherit it or implement
// the same set of methods
type TBaseControl struct {
	refID types.AUnicalID

	pos           types.IPos
	width, height int
	minW, minH    int
	scale         int
	gapX, gapY    int
	padX          types.ACoordX
	padY          types.ACoordY
	fg, bg        term.Attribute
	fgActive      term.Attribute
	bgActive      term.Attribute
	align         Align
	parent        Control
	inactive      bool
	modal         bool
	tabSkip       bool
	disabled      bool
	hidden        bool
	clipped       bool
	clipper       *rect
	pack          PackType
	children      []Control
	mtx           sync.RWMutex
	onActive      func(active bool)
	style         string
	title         string
}

var (
	globalRefId int64
)

func nextRefId() types.AUnicalID {
	id := atomic.AddInt64(&globalRefId, 1)
	return types.AUnicalID(id)
}

func NewBaseControl() TBaseControl {
	return TBaseControl{
		refID: nextRefId(),
		pos:   pos.New(),
	}
}

func (c *TBaseControl) SetClipped(clipped bool) {
	c.clipped = clipped
}

func (c *TBaseControl) Clipped() bool {
	return c.clipped
}

func (c *TBaseControl) SetStyle(style string) {
	c.style = style
}

func (c *TBaseControl) Style() string {
	return c.style
}

func (c *TBaseControl) RefID() types.AUnicalID {
	return c.refID
}

func (c *TBaseControl) Title() string {
	return c.title
}

func (c *TBaseControl) SetTitle(title string) {
	c.title = title
}

func (c *TBaseControl) Size() (widht int, height int) {
	return c.width, c.height
}

func (c *TBaseControl) SetSize(width, height int) {
	if width < c.minW {
		width = c.minW
	}
	if height < c.minH {
		height = c.minH
	}

	if height != c.height || width != c.width {
		c.height = height
		c.width = width
	}
}

func (c *TBaseControl) Pos() types.IPos {
	return c.pos
}

func (c *TBaseControl) SetPos(x types.ACoordX, y types.ACoordY) {
	if c.clipped && c.clipper != nil {
		cx, cy, _, _ := c.Clipper()
		px, py := c.Paddings()

		distX := cx - c.pos.GetX()
		distY := cy - c.pos.GetY()

		c.clipper.x = x + px
		c.clipper.y = y + py

		c.pos.Set((x-distX)+px, (y-distY)+py)
	} else {
		c.pos.Set(x, y)
	}
}

func (c *TBaseControl) applyConstraints() {
	ww, hh := c.width, c.height
	if ww < c.minW {
		ww = c.minW
	}
	if hh < c.minH {
		hh = c.minH
	}
	if hh != c.height || ww != c.width {
		c.SetSize(ww, hh)
	}
}

func (c *TBaseControl) Constraints() (minw int, minh int) {
	return c.minW, c.minH
}

func (c *TBaseControl) SetConstraints(minw, minh int) {
	c.minW = minw
	c.minH = minh
	c.applyConstraints()
}

func (c *TBaseControl) Active() bool {
	return !c.inactive
}

func (c *TBaseControl) SetActive(active bool) {
	c.inactive = !active

	if c.onActive != nil {
		c.onActive(active)
	}
}

func (c *TBaseControl) OnActive(fn func(active bool)) {
	c.onActive = fn
}

func (c *TBaseControl) TabStop() bool {
	return !c.tabSkip
}

func (c *TBaseControl) SetTabStop(tabstop bool) {
	c.tabSkip = !tabstop
}

func (c *TBaseControl) Enabled() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return !c.disabled
}

func (c *TBaseControl) SetEnabled(enabled bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.disabled = !enabled
}

func (c *TBaseControl) Visible() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return !c.hidden
}

func (c *TBaseControl) SetVisible(visible bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if visible == !c.hidden {
		return
	}

	c.hidden = !visible
	if c.parent == nil {
		return
	}

	p := c.Parent()
	for p.Parent() != nil {
		p = p.Parent()
	}

	go func() {
		if FindFirstActiveControl(c) != nil && !c.inactive {
			PutEvent(Event{Type: EventKey, Key: term.KeyTab})
		}
		PutEvent(Event{Type: EventLayout, Target: p})
	}()
}

func (c *TBaseControl) Parent() Control {
	return c.parent
}

func (c *TBaseControl) SetParent(parent Control) {
	if c.parent == nil {
		c.parent = parent
	}
}

func (c *TBaseControl) Modal() bool {
	return c.modal
}

func (c *TBaseControl) SetModal(modal bool) {
	c.modal = modal
}

func (c *TBaseControl) Paddings() (px types.ACoordX, py types.ACoordY) {
	return c.padX, c.padY
}

func (c *TBaseControl) SetPaddings(px types.ACoordX, py types.ACoordY) {
	if px >= 0 {
		c.padX = px
	}
	if py >= 0 {
		c.padY = py
	}
}

func (c *TBaseControl) Gaps() (dx int, dy int) {
	return c.gapX, c.gapY
}

func (c *TBaseControl) SetGaps(dx, dy int) {
	if dx >= 0 {
		c.gapX = dx
	}
	if dy >= 0 {
		c.gapY = dy
	}
}

func (c *TBaseControl) Pack() PackType {
	return c.pack
}

func (c *TBaseControl) SetPack(pack PackType) {
	c.pack = pack
}

func (c *TBaseControl) Scale() int {
	return c.scale
}

func (c *TBaseControl) SetScale(scale int) {
	if scale >= 0 {
		c.scale = scale
	}
}

func (c *TBaseControl) Align() Align {
	return c.align
}

func (c *TBaseControl) SetAlign(align Align) {
	c.align = align
}

func (c *TBaseControl) TextColor() term.Attribute {
	return c.fg
}

func (c *TBaseControl) SetTextColor(clr term.Attribute) {
	c.fg = clr
}

func (c *TBaseControl) BackColor() term.Attribute {
	return c.bg
}

func (c *TBaseControl) SetBackColor(clr term.Attribute) {
	c.bg = clr
}

func (c *TBaseControl) childCount() int {
	cnt := 0
	for _, child := range c.children {
		if child.Visible() {
			cnt++
		}
	}

	return cnt
}

func (c *TBaseControl) ResizeChildren() {
	children := c.childCount()
	if children == 0 {
		return
	}

	fullWidth := c.width - int(2*c.padX)
	fullHeight := c.height - int(2*c.padY)
	if c.pack == Horizontal {
		fullWidth -= (children - 1) * c.gapX
	} else {
		fullHeight -= (children - 1) * c.gapY
	}

	totalSc := c.ChildrenScale()
	minWidth := 0
	minHeight := 0
	for _, child := range c.children {
		if !child.Visible() {
			continue
		}

		cw, ch := child.MinimalSize()
		if c.pack == Horizontal {
			minWidth += cw
		} else {
			minHeight += ch
		}
	}

	aStep := 0
	diff := fullWidth - minWidth
	if c.pack == Vertical {
		diff = fullHeight - minHeight
	}
	if totalSc > 0 {
		aStep = int(float32(diff) / float32(totalSc))
	}

	for _, ctrl := range c.children {
		if !ctrl.Visible() {
			continue
		}

		tw, th := ctrl.MinimalSize()
		sc := ctrl.Scale()
		d := ctrl.Scale() * aStep
		if c.pack == Horizontal {
			if sc != 0 {
				if sc == totalSc {
					tw += diff
					d = diff
				} else {
					tw += d
				}
			}
			th = fullHeight
		} else {
			if sc != 0 {
				if sc == totalSc {
					th += diff
					d = diff
				} else {
					th += d
				}
			}
			tw = fullWidth
		}
		diff -= d
		totalSc -= sc

		ctrl.SetSize(tw, th)
		ctrl.ResizeChildren()
	}
}

func (c *TBaseControl) AddChild(control Control) {
	if c.children == nil {
		c.children = make([]Control, 1)
		c.children[0] = control
	} else {
		if c.ChildExists(control) {
			panic("Double adding a child")
		}

		c.children = append(c.children, control)
	}

	var ctrl Control
	var mainCtrl Control
	ctrl = c
	for ctrl != nil {
		ww, hh := ctrl.MinimalSize()
		cw, ch := ctrl.Size()
		if ww > cw || hh > ch {
			if ww > cw {
				cw = ww
			}
			if hh > ch {
				ch = hh
			}
			ctrl.SetConstraints(cw, ch)
		}

		if ctrl.Parent() == nil {
			mainCtrl = ctrl
		}
		ctrl = ctrl.Parent()
	}

	if mainCtrl != nil {
		mainCtrl.ResizeChildren()
		mainCtrl.PlaceChildren()
	}

	if c.clipped && c.clipper == nil {
		c.setClipper()
	}
}

func (c *TBaseControl) Children() []Control {
	child := make([]Control, len(c.children))
	copy(child, c.children)
	return child
}

func (c *TBaseControl) ChildExists(control Control) bool {
	if len(c.children) == 0 {
		return false
	}

	for _, ctrl := range c.children {
		if ctrl == control {
			return true
		}
	}

	return false
}

func (c *TBaseControl) ChildrenScale() int {
	if c.childCount() == 0 {
		return c.scale
	}

	total := 0
	for _, ctrl := range c.children {
		if ctrl.Visible() {
			total += ctrl.Scale()
		}
	}

	return total
}

func (c *TBaseControl) MinimalSize() (w int, h int) {
	children := c.childCount()
	if children == 0 {
		return c.minW, c.minH
	}

	totalX := 2 * c.padX
	totalY := 2 * c.padY

	if c.pack == Vertical {
		totalY += types.ACoordY((children - 1) * c.gapY)
	} else {
		totalX += types.ACoordX((children - 1) * c.gapX)
	}

	for _, ctrl := range c.children {
		if ctrl.Clipped() {
			continue
		}

		if !ctrl.Visible() {
			continue
		}
		ww, hh := ctrl.MinimalSize()
		if c.pack == Vertical {
			totalY += types.ACoordY(hh)
			if types.ACoordX(ww)+2*c.padX > totalX {
				totalX = types.ACoordX(ww) + 2*c.padX
			}
		} else {
			totalX += types.ACoordX(ww)
			if types.ACoordY(hh)+2*c.padY > totalY {
				totalY = types.ACoordY(hh) + 2*c.padY
			}
		}
	}

	if totalX < types.ACoordX(c.minW) {
		totalX = types.ACoordX(c.minW)
	}
	if totalY < types.ACoordY(c.minH) {
		totalY = types.ACoordY(c.minH)
	}

	return int(totalX), int(totalY)
}

func (c *TBaseControl) Draw() {
	panic("BaseControl Draw Called")
}

func (c *TBaseControl) DrawChildren() {
	if c.hidden {
		return
	}

	PushClip()
	defer PopClip()

	cp := ClippedParent(c)
	var cTarget Control

	cTarget = c
	if cp != nil {
		cTarget = cp
	}

	x, y, w, h := cTarget.Clipper()
	SetClipRect(x, y, w, h)

	for _, child := range c.children {
		child.Draw()
	}
}

func (c *TBaseControl) Clipper() (types.ACoordX, types.ACoordY, int, int) {
	clipped := ClippedParent(c)

	if clipped == nil || (c.clipped && c.clipper != nil) {
		return c.clipper.x, c.clipper.y, c.clipper.w, c.clipper.h
	}

	return CalcClipper(c)
}

func (c *TBaseControl) setClipper() {
	x, y, w, h := CalcClipper(c)
	c.clipper = &rect{x: x, y: y, w: w, h: h}
}

func (c *TBaseControl) HitTest(x types.ACoordX, y types.ACoordY) HitResult {
	if x > c.pos.GetX() && x < c.pos.GetX()+types.ACoordX(c.width-1) &&
		y > c.pos.GetY() && y < c.pos.GetY()+types.ACoordY(c.height-1) {
		return HitInside
	}

	if (x == c.pos.GetX() || x == c.pos.GetX()+types.ACoordX(c.width-1)) &&
		y >= c.pos.GetY() && y < c.pos.GetY()+types.ACoordY(c.height) {
		return HitBorder
	}

	if (y == c.pos.GetY() || y == c.pos.GetY()+types.ACoordY(c.height-1)) &&
		x >= c.pos.GetX() && x < c.pos.GetX()+types.ACoordX(c.width) {
		return HitBorder
	}

	return HitOutside
}

func (c *TBaseControl) ProcessEvent(ev Event) bool {
	return SendEventToChild(c, ev)
}

func (c *TBaseControl) PlaceChildren() {
	children := c.childCount()
	if c.children == nil || children == 0 {
		return
	}

	xx, yy := c.pos.GetX()+c.padX, c.pos.GetY()+c.padY
	for _, ctrl := range c.children {
		if !ctrl.Visible() {
			continue
		}

		ctrl.SetPos(xx, yy)
		ww, hh := ctrl.Size()
		if c.pack == Vertical {
			yy += types.ACoordY(c.gapY + hh)
		} else {
			xx += types.ACoordX(c.gapX + ww)
		}

		ctrl.PlaceChildren()
	}
}

// ActiveColors return the attributes for the controls when it
// is active: text and background colors
func (c *TBaseControl) ActiveColors() (term.Attribute, term.Attribute) {
	return c.fgActive, c.bgActive
}

// SetActiveTextColor changes text color of the active control
func (c *TBaseControl) SetActiveTextColor(clr term.Attribute) {
	c.fgActive = clr
}

// SetActiveBackColor changes background color of the active control
func (c *TBaseControl) SetActiveBackColor(clr term.Attribute) {
	c.bgActive = clr
}

func (c *TBaseControl) removeChild(control Control) {
	children := []Control{}

	for _, child := range c.children {
		if child.RefID() == control.RefID() {
			continue
		}

		children = append(children, child)
	}
	c.children = nil

	for _, child := range children {
		c.AddChild(child)
	}
}

// Destroy removes an object from its parental chain
func (c *TBaseControl) Destroy() {
	c.parent.removeChild(c)
	c.parent.SetConstraints(0, 0)
}
