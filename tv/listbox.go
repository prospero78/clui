package tv

import (
	"strings"

	term "github.com/nsf/termbox-go"

	"github.com/prospero78/goTV/tv/autoheight"
	"github.com/prospero78/goTV/tv/autowidth"
	"github.com/prospero78/goTV/tv/types"
)

/*
ListBox is control to display a list of items and allow to user to select any of them.
Content is scrollable with arrow keys or by clicking up and bottom buttons
on the scroll(now content is scrollable with mouse dragging only on Windows).

ListBox calls onSelectItem item function after a user changes currently
selected item with mouse or using keyboard. Event structure has 2 fields filled:
Y - selected item number in list(-1 if nothing is selected),
Msg - text of the selected item.
*/
type ListBox struct {
	TBaseControl
	// own listbox members
	items         []string
	currSelection int
	topLine       int
	buttonPos     int

	onSelectItem func(Event)
	onKeyPress   func(term.Key) bool

	autoWidth  types.IAutoWidth
	autoHeight types.IAutoHeight
}

/*
CreateListBox creates a new frame.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and height - are minimal size of the control.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateListBox(parent IControl, width, height int, scale int) *ListBox {
	l := &ListBox{
		TBaseControl: NewBaseControl(),
		autoWidth:    autowidth.New(),
		autoHeight:   autoheight.New(),
	}

	if height == 0 {
		height = 3
		l.autoHeight.Set()
	}
	if width == 0 {
		width = 5
		l.autoWidth.Set()
	}

	l.SetSize(width, height)
	l.SetConstraints(width, height)
	l.currSelection = -1
	l.items = make([]string, 0)
	l.topLine = 0
	l.parent = parent
	l.buttonPos = -1

	l.SetTabStop(true)
	l.SetScale(scale)

	l.onSelectItem = nil

	if parent != nil {
		parent.AddChild(l)
	}

	return l
}

func (l *ListBox) drawScroll() {
	PushAttributes()
	defer PopAttributes()

	pos := ThumbPosition(l.currSelection, len(l.items), int(l.height.Get()))
	l.buttonPos = pos

	DrawScrollBar(l.pos.GetX()+types.ACoordX(l.width.Get()-1), l.pos.GetY(), 1, int(l.height.Get()), pos)
}

func (l *ListBox) drawItems() {
	PushAttributes()
	defer PopAttributes()

	maxCurr := len(l.items) - 1
	curr := l.topLine
	dy := types.ACoordY(0)
	maxDy := types.ACoordY(l.height.Get() - 1)
	maxWidth := l.width.Get() - 1

	fg, bg := RealColor(l.fg, l.Style(), ColorEditText), RealColor(l.bg, l.Style(), ColorEditBack)
	if l.Active() {
		fg, bg = RealColor(l.fg, l.Style(), ColorEditActiveText), RealColor(l.bg, l.Style(), ColorEditActiveBack)
	}
	fgSel, bgSel := RealColor(l.fgActive, l.Style(), ColorSelectionText), RealColor(l.bgActive, l.Style(), ColorSelectionBack)

	for curr <= maxCurr && dy <= maxDy {
		f, b := fg, bg
		if curr == l.currSelection {
			f, b = fgSel, bgSel
		}

		SetTextColor(f)
		SetBackColor(b)
		FillRect(l.pos.GetX(), l.pos.GetY()+dy, int(l.width.Get()-1), 1, ' ')
		str := SliceColorized(l.items[curr], 0, int(maxWidth))
		DrawText(l.pos.GetX(), l.pos.GetY()+dy, str)

		curr++
		dy++
	}
}

// Draw repaints the control on its View surface
func (l *ListBox) Draw() {
	if l.hidden {
		return
	}

	PushAttributes()
	defer PopAttributes()

	x, y := l.pos.Get()
	w, h := l.Size()

	fg, bg := RealColor(l.fg, l.Style(), ColorEditText), RealColor(l.bg, l.Style(), ColorEditBack)
	if l.Active() {
		fg, bg = RealColor(l.fg, l.Style(), ColorEditActiveText), RealColor(l.bg, l.Style(), ColorEditActiveBack)
	}
	SetTextColor(fg)
	SetBackColor(bg)
	FillRect(x, y, w, h, ' ')
	l.drawItems()
	l.drawScroll()
}

func (l *ListBox) home() {
	if l.currSelection == 0 {
		return
	}

	if len(l.items) > 0 {
		l.currSelection = 0
	}
	l.topLine = 0

	if l.onSelectItem != nil {
		ev := Event{Y: types.ACoordY(l.currSelection), Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

func (l *ListBox) End() {
	length := len(l.items)

	if length == 0 || l.currSelection == length-1 {
		return
	}

	l.currSelection = length - 1
	if length > int(l.height.Get()) {
		l.topLine = length - int(l.height.Get())
	}

	if l.onSelectItem != nil {
		ev := Event{Y: types.ACoordY(l.currSelection), Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

func (l *ListBox) moveUp(dy int) {
	if l.topLine == 0 && l.currSelection == 0 {
		return
	}

	if l.currSelection == -1 {
		if len(l.items) != 0 {
			l.currSelection = 0
		}
		return
	}

	if l.currSelection < dy {
		l.currSelection = 0
	} else {
		l.currSelection -= dy
	}

	l.EnsureVisible()

	if l.onSelectItem != nil {
		ev := Event{Y: types.ACoordY(l.currSelection), Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

func (l *ListBox) moveDown(dy int) {
	length := len(l.items)

	if length == 0 || l.currSelection == length-1 {
		return
	}

	if l.currSelection+dy >= length {
		l.currSelection = length - 1
	} else {
		l.currSelection += dy
	}

	l.EnsureVisible()

	if l.onSelectItem != nil {
		ev := Event{Y: types.ACoordY(l.currSelection), Msg: l.SelectedItemText()}
		l.onSelectItem(ev)
	}
}

// EnsureVisible makes the currently selected item visible and scrolls the item list if it is required
func (l *ListBox) EnsureVisible() {
	length := len(l.items)

	if length <= int(l.height.Get()) || l.currSelection == -1 {
		return
	}

	diff := l.currSelection - l.topLine
	if diff >= 0 && diff < int(l.height.Get()) {
		return
	}

	if diff < 0 {
		l.topLine = l.currSelection
	} else {
		top := l.currSelection - int(l.height.Get()) + 1
		if length-top > int(l.height.Get()) {
			l.topLine = top
		} else {
			l.topLine = length - int(l.height.Get())
		}
	}
}

// Clear deletes all ListBox items
func (l *ListBox) Clear() {
	l.items = make([]string, 0)
	l.currSelection = -1
	l.topLine = 0
}

func (l *ListBox) processMouseClick(ev Event) bool {
	if ev.Key != term.MouseLeft {
		return false
	}

	dx := ev.X - l.pos.GetX()
	dy := ev.Y - l.pos.GetY()

	if dx == types.ACoordX(l.width.Get()-1) {
		if dy < 0 || int(dy) >= int(l.height.Get()) || len(l.items) < 2 {
			return true
		}

		if dy == 0 {
			l.moveUp(1)
			return true
		}
		if int(dy) == int(l.height.Get())-1 {
			l.moveDown(1)
			return true
		}

		l.buttonPos = int(dy)
		l.recalcPositionByScroll()
		return true
	}

	if dx < 0 || int(dx) >= int(l.width.Get()) || dy < 0 || int(dy) >= int(l.height.Get()) {
		return true
	}

	if int(dy) >= len(l.items) {
		return true
	}

	l.SelectItem(l.topLine + int(dy))
	WindowManager().BeginUpdate()
	onSelFunc := l.onSelectItem
	WindowManager().EndUpdate()
	if onSelFunc != nil {
		ev := Event{Y: types.ACoordY(l.topLine) + dy, Msg: l.SelectedItemText()}
		onSelFunc(ev)
	}

	return true
}

func (l *ListBox) recalcPositionByScroll() {
	newPos := ItemByThumbPosition(l.buttonPos, len(l.items), int(l.height.Get()))
	if newPos < 1 {
		return
	}

	l.currSelection = newPos
	l.EnsureVisible()
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (l *ListBox) ProcessEvent(event Event) bool {
	if !l.Active() || !l.Enabled() {
		return false
	}

	switch event.Type {
	case EventKey:
		if l.onKeyPress != nil {
			res := l.onKeyPress(event.Key)
			if res {
				return true
			}
		}

		switch event.Key {
		case term.KeyHome:
			l.home()
			return true
		case term.KeyEnd:
			l.End()
			return true
		case term.KeyArrowUp:
			l.moveUp(1)
			return true
		case term.KeyArrowDown:
			l.moveDown(1)
			return true
		case term.KeyPgdn:
			l.moveDown(int(l.height.Get()))
			return true
		case term.KeyPgup:
			l.moveUp(int(l.height.Get()))
			return true
		case term.KeyCtrlM:
			if l.currSelection != -1 && l.onSelectItem != nil {
				ev := Event{Y: types.ACoordY(l.currSelection), Msg: l.SelectedItemText()}
				l.onSelectItem(ev)
			}
		default:
			return false
		}
	case EventMouse:
		return l.processMouseClick(event)
	}

	return false
}

// own methods

// AddItem adds a new item to item list.
// Returns true if the operation is successful
func (l *ListBox) AddItem(item string) bool {
	l.items = append(l.items, item)
	return true
}

// SelectItem selects item which number in the list equals
// id. If the item exists the ListBox scrolls the list to
// make the item visible.
// Returns true if the item is selected successfully
func (l *ListBox) SelectItem(id int) bool {
	if len(l.items) <= id || id < 0 {
		return false
	}

	l.currSelection = id
	l.EnsureVisible()
	return true
}

// Item returns item text by its index.
// If index is out of range an empty string and false are returned
func (l *ListBox) Item(id int) (string, bool) {
	if len(l.items) <= id || id < 0 {
		return "", false
	}

	return l.items[id], true
}

// FindItem looks for an item in list which text equals
// to text, by default the search is casesensitive.
// Returns item number in item list or -1 if nothing is found.
func (l *ListBox) FindItem(text string, caseSensitive bool) int {
	for idx, itm := range l.items {
		if itm == text || (caseSensitive && strings.EqualFold(itm, text)) {
			return idx
		}
	}

	return -1
}

// PartialFindItem looks for an item in list which text starts from
// the given substring, by default the search is casesensitive.
// Returns item number in item list or -1 if nothing is found.
func (l *ListBox) PartialFindItem(text string, caseSensitive bool) int {
	if !caseSensitive {
		text = strings.ToLower(text)
	}

	for idx, itm := range l.items {
		if caseSensitive {
			if strings.HasPrefix(itm, text) {
				return idx
			}
		} else {
			low := strings.ToLower(itm)
			if strings.HasPrefix(low, text) {
				return idx
			}
		}
	}

	return -1
}

// SelectedItem returns currently selected item id
func (l *ListBox) SelectedItem() int {
	return l.currSelection
}

// SelectedItemText returns text of currently selected item or empty sting if nothing is
// selected or ListBox is empty.
func (l *ListBox) SelectedItemText() string {
	if l.currSelection == -1 {
		return ""
	}

	return l.items[l.currSelection]
}

// RemoveItem deletes an item which number is id in item list
// Returns true if item is deleted
func (l *ListBox) RemoveItem(id int) bool {
	if id < 0 || id >= len(l.items) {
		return false
	}

	l.items = append(l.items[:id], l.items[id+1:]...)
	return true
}

// OnSelectItem sets a callback that is called every time
// the selected item is changed
func (l *ListBox) OnSelectItem(fn func(Event)) {
	l.onSelectItem = fn
}

// OnKeyPress sets the callback that is called when a user presses a Key while
// the controls is active. If a handler processes the key it should return
// true. If handler returns false it means that the default handler will
// process the key
func (l *ListBox) OnKeyPress(fn func(term.Key) bool) {
	l.onKeyPress = fn
}

// ItemCount returns the number of items in the ListBox
func (l *ListBox) ItemCount() int {
	return len(l.items)
}
