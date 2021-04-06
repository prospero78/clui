package tv

import (
	"sync"
	"sync/atomic"

	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/pos"
	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/width"
)

// TBaseControl is a base for all visible controls.
// Every new control must inherit it or implement
// the same set of methods
type TBaseControl struct {
	refID types.AUnicalID

	pos        types.IPos
	width      types.IWidth
	height     types.AHeight
	minW, minH int
	scale      int
	gapX, gapY int
	posPad     types.IPos
	fg, bg     term.Attribute
	fgActive   term.Attribute
	bgActive   term.Attribute
	align      Align
	parent     IControl
	inactive   bool
	modal      bool
	tabSkip    bool
	disabled   bool
	hidden     bool
	clipped    bool
	clipper    *rect
	pack       PackType
	children   []IControl
	mtx        sync.RWMutex
	onActive   func(active bool)
	style      string
	title      string
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
		refID:  nextRefId(),
		pos:    pos.New(),
		posPad: pos.New(),
		width:  width.New(),
	}
}

func (sf *TBaseControl) SetClipped(clipped bool) {
	sf.clipped = clipped
}

func (sf *TBaseControl) Clipped() bool {
	return sf.clipped
}

func (sf *TBaseControl) SetStyle(style string) {
	sf.style = style
}

func (sf *TBaseControl) Style() string {
	return sf.style
}

func (sf *TBaseControl) RefID() types.AUnicalID {
	return sf.refID
}

func (sf *TBaseControl) Title() string {
	return sf.title
}

func (sf *TBaseControl) SetTitle(title string) {
	sf.title = title
}

func (sf *TBaseControl) Size() (widht int, height int) {
	return int(sf.width.Get()), int(sf.height)
}

func (sf *TBaseControl) SetSize(width, height int) {
	if width < sf.minW {
		width = sf.minW
	}
	if height < sf.minH {
		height = sf.minH
	}

	if height != int(sf.height) || width != int(sf.width.Get()) {
		sf.height = types.AHeight(height)
		sf.width.Set(types.AWidth(width))
	}
}

func (sf *TBaseControl) Pos() types.IPos {
	return sf.pos
}

func (sf *TBaseControl) SetPos(x types.ACoordX, y types.ACoordY) {
	if sf.clipped && sf.clipper != nil {
		cx, cy, _, _ := sf.Clipper()
		px, py := sf.Paddings()

		distX := cx - sf.pos.GetX()
		distY := cy - sf.pos.GetY()

		sf.clipper.x = x + px
		sf.clipper.y = y + py

		sf.pos.Set((x-distX)+px, (y-distY)+py)
	} else {
		sf.pos.Set(x, y)
	}
}

func (sf *TBaseControl) applyConstraints() {
	ww, hh := sf.width, sf.height
	if int(ww.Get()) < sf.minW {
		ww.Set(types.AWidth(sf.minW))
	}
	if int(hh) < sf.minH {
		hh = types.AHeight(sf.minH)
	}
	if hh != sf.height || ww != sf.width {
		sf.SetSize(int(ww.Get()), int(hh))
	}
}

func (sf *TBaseControl) Constraints() (minw int, minh int) {
	return sf.minW, sf.minH
}

func (sf *TBaseControl) SetConstraints(minw, minh int) {
	sf.minW = minw
	sf.minH = minh
	sf.applyConstraints()
}

func (sf *TBaseControl) Active() bool {
	return !sf.inactive
}

func (sf *TBaseControl) SetActive(active bool) {
	sf.inactive = !active

	if sf.onActive != nil {
		sf.onActive(active)
	}
}

func (sf *TBaseControl) OnActive(fn func(active bool)) {
	sf.onActive = fn
}

func (sf *TBaseControl) TabStop() bool {
	return !sf.tabSkip
}

func (sf *TBaseControl) SetTabStop(tabstop bool) {
	sf.tabSkip = !tabstop
}

func (sf *TBaseControl) Enabled() bool {
	sf.mtx.RLock()
	defer sf.mtx.RUnlock()

	return !sf.disabled
}

func (c *TBaseControl) SetEnabled(enabled bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.disabled = !enabled
}

func (sf *TBaseControl) Visible() bool {
	sf.mtx.RLock()
	defer sf.mtx.RUnlock()

	return !sf.hidden
}

func (sf *TBaseControl) SetVisible(visible bool) {
	sf.mtx.Lock()
	defer sf.mtx.Unlock()

	if visible == !sf.hidden {
		return
	}

	sf.hidden = !visible
	if sf.parent == nil {
		return
	}

	p := sf.Parent()
	for p.Parent() != nil {
		p = p.Parent()
	}

	go func() {
		if FindFirstActiveControl(sf) != nil && !sf.inactive {
			PutEvent(Event{Type: EventKey, Key: term.KeyTab})
		}
		PutEvent(Event{Type: EventLayout, Target: p})
	}()
}

func (sf *TBaseControl) Parent() IControl {
	return sf.parent
}

func (sf *TBaseControl) SetParent(parent IControl) {
	if sf.parent == nil {
		sf.parent = parent
	}
}

func (sf *TBaseControl) Modal() bool {
	return sf.modal
}

func (sf *TBaseControl) SetModal(modal bool) {
	sf.modal = modal
}

func (sf *TBaseControl) Paddings() (px types.ACoordX, py types.ACoordY) {
	return sf.posPad.Get()
}

func (sf *TBaseControl) SetPaddings(px types.ACoordX, py types.ACoordY) {
	if px >= 0 {
		sf.posPad.X().Set(px)
	}
	if py >= 0 {
		sf.posPad.Y().Set(py)
	}
}

func (sf *TBaseControl) Gaps() (dx int, dy int) {
	return sf.gapX, sf.gapY
}

func (sf *TBaseControl) SetGaps(dx, dy int) {
	if dx >= 0 {
		sf.gapX = dx
	}
	if dy >= 0 {
		sf.gapY = dy
	}
}

func (sf *TBaseControl) Pack() PackType {
	return sf.pack
}

func (sf *TBaseControl) SetPack(pack PackType) {
	sf.pack = pack
}

func (sf *TBaseControl) Scale() int {
	return sf.scale
}

func (sf *TBaseControl) SetScale(scale int) {
	if scale >= 0 {
		sf.scale = scale
	}
}

func (sf *TBaseControl) Align() Align {
	return sf.align
}

func (sf *TBaseControl) SetAlign(align Align) {
	sf.align = align
}

func (sf *TBaseControl) TextColor() term.Attribute {
	return sf.fg
}

func (sf *TBaseControl) SetTextColor(clr term.Attribute) {
	sf.fg = clr
}

func (sf *TBaseControl) BackColor() term.Attribute {
	return sf.bg
}

func (sf *TBaseControl) SetBackColor(clr term.Attribute) {
	sf.bg = clr
}

func (sf *TBaseControl) childCount() int {
	cnt := 0
	for _, child := range sf.children {
		if child.Visible() {
			cnt++
		}
	}

	return cnt
}

func (sf *TBaseControl) ResizeChildren() {
	children := sf.childCount()
	if children == 0 {
		return
	}

	fullWidth := int(sf.width.Get()) - int(2*sf.posPad.GetX())
	fullHeight := int(sf.height) - int(2*sf.posPad.Y().Get())
	if sf.pack == Horizontal {
		fullWidth -= (children - 1) * sf.gapX
	} else {
		fullHeight -= (children - 1) * sf.gapY
	}

	totalSc := sf.ChildrenScale()
	minWidth := 0
	minHeight := 0
	for _, child := range sf.children {
		if !child.Visible() {
			continue
		}

		cw, ch := child.MinimalSize()
		if sf.pack == Horizontal {
			minWidth += cw
		} else {
			minHeight += ch
		}
	}

	aStep := 0
	diff := fullWidth - minWidth
	if sf.pack == Vertical {
		diff = fullHeight - minHeight
	}
	if totalSc > 0 {
		aStep = int(float32(diff) / float32(totalSc))
	}

	for _, ctrl := range sf.children {
		if !ctrl.Visible() {
			continue
		}

		tw, th := ctrl.MinimalSize()
		sc := ctrl.Scale()
		d := ctrl.Scale() * aStep
		if sf.pack == Horizontal {
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

func (sf *TBaseControl) AddChild(control IControl) {
	if sf.children == nil {
		sf.children = make([]IControl, 1)
		sf.children[0] = control
	} else {
		if sf.ChildExists(control) {
			panic("Double adding a child")
		}

		sf.children = append(sf.children, control)
	}

	var ctrl IControl
	var mainCtrl IControl
	ctrl = sf
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

	if sf.clipped && sf.clipper == nil {
		sf.setClipper()
	}
}

func (sf *TBaseControl) Children() []IControl {
	child := make([]IControl, len(sf.children))
	copy(child, sf.children)
	return child
}

func (sf *TBaseControl) ChildExists(control IControl) bool {
	if len(sf.children) == 0 {
		return false
	}

	for _, ctrl := range sf.children {
		if ctrl == control {
			return true
		}
	}

	return false
}

func (sf *TBaseControl) ChildrenScale() int {
	if sf.childCount() == 0 {
		return sf.scale
	}

	total := 0
	for _, ctrl := range sf.children {
		if ctrl.Visible() {
			total += ctrl.Scale()
		}
	}

	return total
}

func (sf *TBaseControl) MinimalSize() (w int, h int) {
	children := sf.childCount()
	if children == 0 {
		return sf.minW, sf.minH
	}

	totalX := 2 * sf.posPad.GetX()
	totalY := 2 * sf.posPad.GetY()

	if sf.pack == Vertical {
		totalY += types.ACoordY((children - 1) * sf.gapY)
	} else {
		totalX += types.ACoordX((children - 1) * sf.gapX)
	}

	for _, ctrl := range sf.children {
		if ctrl.Clipped() {
			continue
		}

		if !ctrl.Visible() {
			continue
		}
		ww, hh := ctrl.MinimalSize()
		if sf.pack == Vertical {
			totalY += types.ACoordY(hh)
			if types.ACoordX(ww)+2*sf.posPad.GetX() > totalX {
				totalX = types.ACoordX(ww) + 2*sf.posPad.GetX()
			}
		} else {
			totalX += types.ACoordX(ww)
			if types.ACoordY(hh)+2*sf.posPad.GetY() > totalY {
				totalY = types.ACoordY(hh) + 2*sf.posPad.GetY()
			}
		}
	}

	if totalX < types.ACoordX(sf.minW) {
		totalX = types.ACoordX(sf.minW)
	}
	if totalY < types.ACoordY(sf.minH) {
		totalY = types.ACoordY(sf.minH)
	}

	return int(totalX), int(totalY)
}

func (sf *TBaseControl) Draw() {
	panic("BaseControl Draw Called")
}

func (sf *TBaseControl) DrawChildren() {
	if sf.hidden {
		return
	}

	PushClip()
	defer PopClip()

	cp := ClippedParent(sf)
	var cTarget IControl

	cTarget = sf
	if cp != nil {
		cTarget = cp
	}

	x, y, w, h := cTarget.Clipper()
	SetClipRect(x, y, w, h)

	for _, child := range sf.children {
		child.Draw()
	}
}

func (sf *TBaseControl) Clipper() (types.ACoordX, types.ACoordY, int, int) {
	clipped := ClippedParent(sf)

	if clipped == nil || (sf.clipped && sf.clipper != nil) {
		return sf.clipper.x, sf.clipper.y, sf.clipper.w, sf.clipper.h
	}

	return CalcClipper(sf)
}

func (sf *TBaseControl) setClipper() {
	x, y, w, h := CalcClipper(sf)
	sf.clipper = &rect{x: x, y: y, w: w, h: h}
}

func (sf *TBaseControl) HitTest(x types.ACoordX, y types.ACoordY) HitResult {
	if x > sf.pos.GetX() && x < sf.pos.GetX()+types.ACoordX(sf.width.Get()-1) &&
		y > sf.pos.GetY() && y < sf.pos.GetY()+types.ACoordY(sf.height-1) {
		return HitInside
	}

	if (x == sf.pos.GetX() || x == sf.pos.GetX()+types.ACoordX(sf.width.Get()-1)) &&
		y >= sf.pos.GetY() && y < sf.pos.GetY()+types.ACoordY(sf.height) {
		return HitBorder
	}

	if (y == sf.pos.GetY() || y == sf.pos.GetY()+types.ACoordY(sf.height-1)) &&
		x >= sf.pos.GetX() && x < sf.pos.GetX()+types.ACoordX(sf.width.Get()) {
		return HitBorder
	}

	return HitOutside
}

func (sf *TBaseControl) ProcessEvent(ev Event) bool {
	return SendEventToChild(sf, ev)
}

func (sf *TBaseControl) PlaceChildren() {
	children := sf.childCount()
	if sf.children == nil || children == 0 {
		return
	}

	xx, yy := sf.pos.GetX()+sf.posPad.GetX(), sf.pos.GetY()+sf.posPad.GetY()
	for _, ctrl := range sf.children {
		if !ctrl.Visible() {
			continue
		}

		ctrl.SetPos(xx, yy)
		ww, hh := ctrl.Size()
		if sf.pack == Vertical {
			yy += types.ACoordY(sf.gapY + hh)
		} else {
			xx += types.ACoordX(sf.gapX + ww)
		}

		ctrl.PlaceChildren()
	}
}

// ActiveColors return the attributes for the controls when it
// is active: text and background colors
func (sf *TBaseControl) ActiveColors() (term.Attribute, term.Attribute) {
	return sf.fgActive, sf.bgActive
}

// SetActiveTextColor changes text color of the active control
func (sf *TBaseControl) SetActiveTextColor(clr term.Attribute) {
	sf.fgActive = clr
}

// SetActiveBackColor changes background color of the active control
func (sf *TBaseControl) SetActiveBackColor(clr term.Attribute) {
	sf.bgActive = clr
}

func (sf *TBaseControl) removeChild(control IControl) {
	children := []IControl{}

	for _, child := range sf.children {
		if child.RefID() == control.RefID() {
			continue
		}

		children = append(children, child)
	}
	sf.children = nil

	for _, child := range children {
		sf.AddChild(child)
	}
}

// Destroy removes an object from its parental chain
func (sf *TBaseControl) Destroy() {
	sf.parent.removeChild(sf)
	sf.parent.SetConstraints(0, 0)
}

// End -- заглушка для интерфейса
func (sf *TBaseControl) End() {
}
