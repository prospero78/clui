// Package wigdetbase -- базовый виджет длявсех виджетов
package widgetbase

import (
	"sync"

	mTerm "github.com/nsf/termbox-go"
	"github.com/prospero78/goTV/tv/types"
	"github.com/prospero78/goTV/tv/widgets/widgetbase/widgetid"
)

// TWidgetBase -- основа для всех виджетов.
// Каждый новый виджет должен включать его или реализовывать
// его методы
type TWidgetBase struct {
	widgetID      types.AWidgetID
	x, y          int
	width, height int // Высота и ширина виджета
	minW, minH    int // Минимальные размеры виджеты
	scale         int
	padX, padY    int
	gapX, gapY    int
	isTabSkip     bool // Пропускать при нажатии клавиши TAB
	isDisabled    bool // Признак отключенности окна
	isHidden      bool // Признак скрытости элемента
	isInactive    bool
	isModal       bool // Признак модальности окна
	isClipped     bool // Признак обрезки виджета
	fg, bg        mTerm.Attribute
	fgActive      mTerm.Attribute
	bgActive      mTerm.Attribute
	align         types.AAlign
	parent        types.IWidget
	pack          types.APackDirect
	children      []types.IWidget
	block         sync.RWMutex
	onActive      func(active bool)
	clipper       *rect
	title         string // Заголовок окна
	style         string // Стиль окна
}

// New -- возвращает новый TBaseControl
func New() TWidgetBase {
	return TWidgetBase{
		widgetID: widgetid.GetWidgetID().NextID(),
	}
}

// SetClipped -- устанавливает признак обрезки виджета
func (sf *TWidgetBase) SetClipped(isClipped bool) {
	sf.isClipped = isClipped
}

func (c *TWidgetBase) Clipped() bool {
	return c.isClipped
}

func (c *TWidgetBase) SetStyle(style string) {
	c.style = style
}

func (c *TWidgetBase) Style() string {
	return c.style
}

func (c *TWidgetBase) RefID() types.AWidgetID {
	return c.widgetID
}

func (c *TWidgetBase) Title() string {
	return c.title
}

func (c *TWidgetBase) SetTitle(title string) {
	c.title = title
}

func (c *TWidgetBase) Size() (widht int, height int) {
	return c.width, c.height
}

func (c *TWidgetBase) SetSize(width, height int) {
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

func (c *TWidgetBase) Pos() (x int, y int) {
	return c.x, c.y
}

func (c *TWidgetBase) SetPos(x, y int) {
	if c.isClipped && c.clipper != nil {
		cx, cy, _, _ := c.Clipper()
		px, py := c.Paddings()

		distX := cx - c.x
		distY := cy - c.y

		c.clipper.x = x + px
		c.clipper.y = y + py

		c.x = (x - distX) + px
		c.y = (y - distY) + py
	} else {
		c.x = x
		c.y = y
	}
}

func (c *TWidgetBase) applyConstraints() {
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

func (c *TWidgetBase) Constraints() (minw int, minh int) {
	return c.minW, c.minH
}

func (c *TWidgetBase) SetConstraints(minw, minh int) {
	c.minW = minw
	c.minH = minh
	c.applyConstraints()
}

func (c *TWidgetBase) Active() bool {
	return !c.isInactive
}

func (c *TWidgetBase) SetActive(active bool) {
	c.isInactive = !active

	if c.onActive != nil {
		c.onActive(active)
	}
}

func (c *TWidgetBase) OnActive(fn func(active bool)) {
	c.onActive = fn
}

func (c *TWidgetBase) TabStop() bool {
	return !c.isTabSkip
}

func (c *TWidgetBase) SetTabStop(tabstop bool) {
	c.isTabSkip = !tabstop
}

func (c *TWidgetBase) Enabled() bool {
	c.block.RLock()
	defer c.block.RUnlock()

	return !c.isDisabled
}

func (c *TWidgetBase) SetEnabled(enabled bool) {
	c.block.Lock()
	defer c.block.Unlock()

	c.isDisabled = !enabled
}

func (c *TWidgetBase) Visible() bool {
	c.block.RLock()
	defer c.block.RUnlock()

	return !c.isHidden
}

func (c *TWidgetBase) SetVisible(visible bool) {
	c.block.Lock()
	defer c.block.Unlock()

	if visible == !c.isHidden {
		return
	}

	c.isHidden = !visible
	if c.parent == nil {
		return
	}

	p := c.Parent()
	for p.Parent() != nil {
		p = p.Parent()
	}

	go func() {
		if FindFirstActiveControl(c) != nil && !c.isInactive {
			PutEvent(Event{Type: EventKey, Key: mTerm.KeyTab})
		}
		PutEvent(Event{Type: EventLayout, Target: p})
	}()
}

func (c *TWidgetBase) Parent() types.IWidget {
	return c.parent
}

func (c *TWidgetBase) SetParent(parent types.IWidget) {
	if c.parent == nil {
		c.parent = parent
	}
}

func (c *TWidgetBase) Modal() bool {
	return c.isModal
}

func (c *TWidgetBase) SetModal(modal bool) {
	c.isModal = modal
}

func (c *TWidgetBase) Paddings() (px int, py int) {
	return c.padX, c.padY
}

func (c *TWidgetBase) SetPaddings(px, py int) {
	if px >= 0 {
		c.padX = px
	}
	if py >= 0 {
		c.padY = py
	}
}

func (c *TWidgetBase) Gaps() (dx int, dy int) {
	return c.gapX, c.gapY
}

func (c *TWidgetBase) SetGaps(dx, dy int) {
	if dx >= 0 {
		c.gapX = dx
	}
	if dy >= 0 {
		c.gapY = dy
	}
}

func (c *TWidgetBase) Pack() PackType {
	return c.pack
}

func (c *TWidgetBase) SetPack(pack PackType) {
	c.pack = pack
}

func (c *TWidgetBase) Scale() int {
	return c.scale
}

func (c *TWidgetBase) SetScale(scale int) {
	if scale >= 0 {
		c.scale = scale
	}
}

func (c *TWidgetBase) Align() Align {
	return c.align
}

func (c *TWidgetBase) SetAlign(align Align) {
	c.align = align
}

func (c *TWidgetBase) TextColor() mTerm.Attribute {
	return c.fg
}

func (c *TWidgetBase) SetTextColor(clr mTerm.Attribute) {
	c.fg = clr
}

func (c *TWidgetBase) BackColor() mTerm.Attribute {
	return c.bg
}

func (c *TWidgetBase) SetBackColor(clr mTerm.Attribute) {
	c.bg = clr
}

func (c *TWidgetBase) childCount() int {
	cnt := 0
	for _, child := range c.children {
		if child.Visible() {
			cnt++
		}
	}

	return cnt
}

func (c *TWidgetBase) ResizeChildren() {
	children := c.childCount()
	if children == 0 {
		return
	}

	fullWidth := c.width - 2*c.padX
	fullHeight := c.height - 2*c.padY
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

func (c *TWidgetBase) AddChild(control types.IWidget) {
	if c.children == nil {
		c.children = make([]types.IWidget, 1)
		c.children[0] = control
	} else {
		if c.ChildExists(control) {
			panic("Double adding a child")
		}

		c.children = append(c.children, control)
	}

	var ctrl types.IWidget
	var mainCtrl types.IWidget
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

	if c.isClipped && c.clipper == nil {
		c.setClipper()
	}
}

func (c *TWidgetBase) Children() []types.IWidget {
	child := make([]types.IWidget, len(c.children))
	copy(child, c.children)
	return child
}

func (c *TWidgetBase) ChildExists(control types.IWidget) bool {
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

func (c *TWidgetBase) ChildrenScale() int {
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

func (c *TWidgetBase) MinimalSize() (w int, h int) {
	children := c.childCount()
	if children == 0 {
		return c.minW, c.minH
	}

	totalX := 2 * c.padX
	totalY := 2 * c.padY

	if c.pack == Vertical {
		totalY += (children - 1) * c.gapY
	} else {
		totalX += (children - 1) * c.gapX
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
			totalY += hh
			if ww+2*c.padX > totalX {
				totalX = ww + 2*c.padX
			}
		} else {
			totalX += ww
			if hh+2*c.padY > totalY {
				totalY = hh + 2*c.padY
			}
		}
	}

	if totalX < c.minW {
		totalX = c.minW
	}
	if totalY < c.minH {
		totalY = c.minH
	}

	return totalX, totalY
}

func (c *TWidgetBase) Draw() {
	panic("BaseControl Draw Called")
}

func (c *TWidgetBase) DrawChildren() {
	if c.isHidden {
		return
	}

	PushClip()
	defer PopClip()

	cp := ClippedParent(c)
	var cTarget types.IWidget

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

func (c *TWidgetBase) Clipper() (int, int, int, int) {
	clipped := ClippedParent(c)

	if clipped == nil || (c.isClipped && c.clipper != nil) {
		return c.clipper.x, c.clipper.y, c.clipper.w, c.clipper.h
	}

	return CalcClipper(c)
}

func (c *TWidgetBase) setClipper() {
	x, y, w, h := CalcClipper(c)
	c.clipper = &rect{x: x, y: y, w: w, h: h}
}

func (c *TWidgetBase) HitTest(x, y int) HitResult {
	if x > c.x && x < c.x+c.width-1 &&
		y > c.y && y < c.y+c.height-1 {
		return HitInside
	}

	if (x == c.x || x == c.x+c.width-1) &&
		y >= c.y && y < c.y+c.height {
		return HitBorder
	}

	if (y == c.y || y == c.y+c.height-1) &&
		x >= c.x && x < c.x+c.width {
		return HitBorder
	}

	return HitOutside
}

func (c *TWidgetBase) ProcessEvent(ev Event) bool {
	return SendEventToChild(c, ev)
}

func (c *TWidgetBase) PlaceChildren() {
	children := c.childCount()
	if c.children == nil || children == 0 {
		return
	}

	xx, yy := c.x+c.padX, c.y+c.padY
	for _, ctrl := range c.children {
		if !ctrl.Visible() {
			continue
		}

		ctrl.SetPos(xx, yy)
		ww, hh := ctrl.Size()
		if c.pack == Vertical {
			yy += c.gapY + hh
		} else {
			xx += c.gapX + ww
		}

		ctrl.PlaceChildren()
	}
}

// ActiveColors return the attributes for the controls when it
// is active: text and background colors
func (c *TWidgetBase) ActiveColors() (mTerm.Attribute, mTerm.Attribute) {
	return c.fgActive, c.bgActive
}

// SetActiveTextColor changes text color of the active control
func (c *TWidgetBase) SetActiveTextColor(clr mTerm.Attribute) {
	c.fgActive = clr
}

// SetActiveBackColor changes background color of the active control
func (c *TWidgetBase) SetActiveBackColor(clr mTerm.Attribute) {
	c.bgActive = clr
}

func (c *TWidgetBase) removeChild(control types.IWidget) {
	children := []types.IWidget{}

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
func (c *TWidgetBase) Destroy() {
	c.parent.removeChild(c)
	c.parent.SetConstraints(0, 0)
}
