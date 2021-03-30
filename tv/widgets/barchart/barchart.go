// Package barchart -- полоска
package barchart

import (
	"fmt"
	"sync/atomic"

	xs "github.com/huandu/xstrings"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/barchart/bardata"
	"github.com/prospero78/goTV/tv/widgets/barchart/bardatacell"
	"github.com/prospero78/goTV/tv/widgets/widgetbase"

	"github.com/prospero78/goTV/tv/types"
)

// TBarChart -- граик представляет сгруппированные данные с помощью
// прямоугольных столбиков. Может быть монохромным - поведени епо умолчанию.
// Обепечивается доступ к индивидуальным цветам каждого столбика
// и исопльзует отрисовку настраиваемых столбиков на цветном экране с
// с зависмостью от значения столбика
// Все столбики имеют одинаковую ширину: устанвливается константой BarSize
// (в случае AutoSize=false), или автоматически вычисляется игнорируя
// BarSize. Столбик тогда не выравнивается в графике и не отображается.
// TBarChart отображает вертикально оси а значениями слева, если ValueWidth
// больше 0, горизонатльная ось с полоской названий (ShowTitles=true,
// включает отображение заголовков --ShowMarks = true ) и график показывает
// легенду справа, если LegendWidth больше 3.
// Если LegendWidth больше графика, то легеда не отображается.
// Тоже самое применимо к ValueWidth.
type TBarChart struct {
	widgetbase.TWidgetBase
	data        []*bardata.TBarData
	autosize    bool
	gap         int32
	barWidth    int32
	legendWidth int32
	valueWidth  int32
	showMarks   bool
	showTitles  bool
	fnDrawCell  func(*bardatacell.TBarDataCell)
}

// New -- возвращает новый *TBarChart.
// view -- TView управляет виджетом
// parent -- контейнер, в который помещается TBarChart. Также TView может
//   отображать родителя в тоже самое время.
// weidthMin и hightMin -- минимальные размеры виджета.
// scale - пут ьмасштабирования виджета когда родитель изменяет размер.
//   Используйте константу DoNotScale если виджет дложен сохранить
//   оригинальный размер.
func New(parent types.IWidget,
	weidthMin, hightMin int,
	scale int) *TBarChart {
	c := new(TBarChart)
	c.TWidgetBase = widgetbase.New()

	if weidthMin == cons.AutoSize {
		weidthMin = 10
	}
	if hightMin == cons.AutoSize {
		hightMin = 5
	}

	c.SetParent(parent)

	c.SetSize(weidthMin, hightMin)
	c.SetConstraints(weidthMin, hightMin)
	c.SetTabStop(true)
	c.showTitles = true
	c.barWidth = 3
	c.data = make([]*bardata.TBarData, 0)
	c.SetScale(scale)

	if parent != nil {
		parent.AddChild(c)
	}

	return c
}

// Draw repaints the control on its View surface
func (sf *TBarChart) Draw() {
	if sf.IsHidden() {
		return
	}

	sf.block.RLock()
	defer sf.block.RUnlock()

	PushAttributes()
	defer PopAttributes()

	fg, bg := RealColor(sf.fg, sf.Style(), cons.ColorBarChartText), RealColor(sf.bg, sf.Style(), cons.ColorBarChartBack)
	SetTextColor(fg)
	SetBackColor(bg)

	FillRect(sf.GetX(), sf.GetY(), sf.GetWidth(), sf.GetHidth(), ' ')

	if len(sf.data) == 0 {
		return
	}

	sf.drawRulers()
	sf.drawValues()
	sf.drawLegend()
	sf.drawBars()
}

func (sf *TBarChart) barHeight() int {
	if sf.showTitles {
		return sf.GetHidth() - 2
	}
	return sf.GetHidth()
}

func (sf *TBarChart) drawBars() {
	if len(sf.data) == 0 {
		return
	}

	start, width := sf.calculateBarArea()
	if width < 2 {
		return
	}

	barW := sf.calculateBarWidth()
	if barW == 0 {
		return
	}

	coeff, max := sf.calculateMultiplier()
	if coeff == 0.0 {
		return
	}

	PushAttributes()
	defer PopAttributes()

	h := sf.barHeight()
	pos := start
	parts := []rune(SysObject(cons.ObjBarChart))
	fg, bg := TextColor(), BackColor()

	for idx, d := range sf.data {
		if pos+barW > start+width {
			break
		}

		fColor, bColor := d.Fg, d.Bg
		ch := d.Ch
		if fColor == cons.ColorDefault {
			fColor = fg
		}
		if bColor == cons.ColorDefault {
			bColor = bg
		}
		if ch == 0 {
			ch = parts[0]
		}

		barH := int(d.Value * coeff)
		if sf.fnDrawCell == nil {
			SetTextColor(fColor)
			SetBackColor(bColor)
			FillRect(sf.GetX()+pos, sf.GetY()+h-barH, barW, barH, ch)
		} else {
			cellDef := TBarDataCell{Item: d.Title, ID: idx,
				Value: 0, BarMax: d.Value, TotalMax: max,
				Fg: fColor, Bg: bColor, Ch: ch}
			for dy := 0; dy < barH; dy++ {
				req := cellDef
				req.Value = max * float64(dy+1) / float64(h)
				sf.fnDrawCell(&req)
				SetTextColor(req.Fg)
				SetBackColor(req.Bg)
				for dx := 0; dx < barW; dx++ {
					PutChar(sf.GetX()+pos+dx, sf.GetY()+h-1-dy, req.Ch)
				}
			}
		}

		if sf.showTitles {
			SetTextColor(fg)
			SetBackColor(bg)
			if sf.showMarks {
				c := parts[7]
				PutChar(sf.GetX()+pos+barW/2, sf.GetY()+h, c)
			}
			var s string
			shift := 0
			if xs.Len(d.Title) > barW {
				s = CutText(d.Title, barW)
			} else {
				shift, s = AlignText(d.Title, barW, cons.AlignCenter)
			}
			DrawRawText(sf.GetX()+pos+shift, sf.GetY()+h+1, s)
		}

		pos += barW + int(sf.BarGap())
	}
}

func (sf *TBarChart) drawLegend() {
	pos, width := sf.calculateBarArea()
	if pos+width >= sf.GetWidth()-3 {
		return
	}

	PushAttributes()
	defer PopAttributes()
	fg, bg := RealColor(sf.fg, sf.Style(), ColorBarChartText), RealColor(sf.bg, sf.Style(), cons.ColorBarChartBack)

	parts := []rune(SysObject(ObjBarChart))
	defRune := parts[0]
	for idx, d := range sf.data {
		if idx >= sf.GetHidth() {
			break
		}

		c := d.Ch
		if c == 0 {
			c = defRune
		}
		SetTextColor(d.Fg)
		SetBackColor(d.Bg)
		PutChar(sf.GetX()+pos+width, sf.GetY()+idx, c)
		s := CutText(fmt.Sprintf(" - %v", d.Title), int(sf.LegendWidth()))
		SetTextColor(fg)
		SetBackColor(bg)
		DrawRawText(sf.GetX()+pos+width+1, sf.GetY()+idx, s)
	}
}

func (sf *TBarChart) drawValues() {
	valVal := int(sf.ValueWidth())
	if valVal <= 0 {
		return
	}

	pos, _ := sf.calculateBarArea()
	if pos == 0 {
		return
	}

	h := sf.barHeight()
	coeff, max := sf.calculateMultiplier()
	if max == coeff {
		return
	}

	dy := 0
	format := fmt.Sprintf("%%%v.2f", valVal)
	for dy < h-1 {
		v := float64(h-dy) / float64(h) * max
		s := fmt.Sprintf(format, v)
		s = CutText(s, valVal)
		DrawRawText(sf.GetX(), sf.GetY()+dy, s)

		dy += 2
	}
}

func (sf *TBarChart) drawRulers() {
	if int(sf.ValueWidth()) <= 0 && int(sf.LegendWidth()) <= 0 && !sf.showTitles {
		return
	}

	pos, vWidth := sf.calculateBarArea()

	parts := []rune(SysObject(ObjBarChart))
	h := sf.barHeight()

	if pos > 0 {
		pos--
		vWidth++
	}

	// horizontal and vertical lines, corner
	cH, cV, cC := parts[1], parts[2], parts[5]

	if pos > 0 {
		for dy := 0; dy < h; dy++ {
			PutChar(sf.GetX()+pos, sf.GetY()+dy, cV)
		}
	}
	if sf.showTitles {
		for dx := 0; dx < vWidth; dx++ {
			PutChar(sf.GetX()+pos+dx, sf.GetY()+h, cH)
		}
	}
	if pos > 0 && sf.showTitles {
		PutChar(sf.GetX()+pos, sf.GetY()+h, cC)
	}
}

func (sf *TBarChart) calculateBarArea() (int, int) {
	w := sf.GetWidth()
	pos := 0

	valVal := int(sf.ValueWidth())
	if valVal < w/2 {
		w = w - valVal - 1
		pos = valVal + 1
	}

	legVal := int(sf.LegendWidth())
	if legVal < w/2 {
		w -= legVal
	}

	return pos, w
}

func (sf *TBarChart) calculateBarWidth() int {
	if len(sf.data) == 0 {
		return 0
	}

	if !sf.autosize {
		return int(sf.MinBarWidth())
	}

	w := sf.GetWidth()
	legVal := int(sf.LegendWidth())
	valVal := int(sf.ValueWidth())
	if valVal < w/2 {
		w = w - valVal - 1
	}
	if legVal < w/2 {
		w -= legVal
	}

	dataCount := len(sf.data)
	gapVal := int(sf.BarGap())
	barVal := int(sf.MinBarWidth())
	minSize := dataCount*barVal + (dataCount-1)*gapVal
	if minSize >= w {
		return barVal
	}

	sz := (w - (dataCount-1)*gapVal) / dataCount
	if sz == 0 {
		sz = 1
	}

	return sz
}

func (sf *TBarChart) calculateMultiplier() (float64, float64) {
	if len(sf.data) == 0 {
		return 0, 0
	}

	h := sf.barHeight()
	if h <= 1 {
		return 0, 0
	}

	max := sf.data[0].Value
	for _, val := range sf.data {
		if val.Value > max {
			max = val.Value
		}
	}

	if max == 0 {
		return 0, 0
	}

	return float64(h) / max, max
}

// AddData appends a new bar to a chart
func (sf *TBarChart) AddData(val TBarData) {
	sf.block.Lock()
	defer sf.block.Unlock()

	sf.data = append(sf.data, val)
}

// ClearData removes all bar from chart
func (sf *TBarChart) ClearData() {
	sf.block.Lock()
	defer sf.block.Unlock()

	sf.data = make([]TBarData, 0)
}

// SetData assign a new bar list to a chart
func (sf *TBarChart) SetData(data []TBarData) {
	sf.block.Lock()
	defer sf.block.Unlock()

	sf.data = make([]TBarData, len(data))
	copy(sf.data, data)
}

// AutoSize returns whether automatic bar width
// calculation is on. If AutoSize is false then all
// bars have width BarWidth. If AutoSize is true then
// bar width is the maximum of three values: BarWidth,
// calculated width that makes all bars fit the
// bar chart area, and 1
func (sf *TBarChart) AutoSize() bool {
	return sf.autosize
}

// SetAutoSize enables or disables automatic bar
// width calculation
func (sf *TBarChart) SetAutoSize(auto bool) {
	sf.block.Lock()
	defer sf.block.Unlock()

	sf.autosize = auto
}

// BarGap returns width of visual gap between two adjacent bars
func (sf *TBarChart) BarGap() int32 {
	return atomic.LoadInt32(&sf.gap)
}

// SetBarGap sets the space width between two adjacent bars
func (sf *TBarChart) SetBarGap(gap int32) {
	atomic.StoreInt32(&sf.gap, gap)
}

// MinBarWidth returns current minimal bar width
func (sf *TBarChart) MinBarWidth() int32 {
	return atomic.LoadInt32(&sf.barWidth)
}

// SetMinBarWidth changes the minimal bar width
func (sf *TBarChart) SetMinBarWidth(size int32) {
	atomic.StoreInt32(&sf.barWidth, size)
}

// ValueWidth returns the width of the area at the left of
// chart used to draw values. Set it to 0 to turn off the
// value panel
func (sf *TBarChart) ValueWidth() int32 {
	return atomic.LoadInt32(&sf.valueWidth)
}

// SetValueWidth changes width of the value panel on the left
func (sf *TBarChart) SetValueWidth(width int32) {
	atomic.StoreInt32(&sf.valueWidth, width)
}

// ShowTitles returns if chart displays horizontal axis and
// bar titles under it
func (sf *TBarChart) ShowTitles() bool {
	return sf.showTitles
}

// SetShowTitles turns on and off horizontal axis and bar titles
func (sf *TBarChart) SetShowTitles(show bool) {
	sf.block.Lock()
	defer sf.block.Unlock()

	sf.showTitles = show
}

// LegendWidth returns width of chart legend displayed at the
// right side of the chart. Set it to 0 to disable legend
func (sf *TBarChart) LegendWidth() int32 {
	return atomic.LoadInt32(&sf.legendWidth)
}

// SetLegendWidth sets new legend panel width
func (sf *TBarChart) SetLegendWidth(width int32) {
	atomic.StoreInt32(&sf.legendWidth, width)
}

// OnDrawCell sets callback that allows to draw multicolored
// bars. BarChart sends the current attrubutes and rune that
// it is going to use to display as well as the current value
// of the bar. A user can change the values of BarDataCell
// depending on some external data or calculations - only
// changing colors and rune makes sense. Changing anything else
// does not affect the chart
func (sf *TBarChart) OnDrawCell(fn func(*TBarDataCell)) {
	sf.block.Lock()
	defer sf.block.Unlock()

	sf.fnDrawCell = fn
}

// ShowMarks returns if horizontal axis has mark under each
// bar. To show marks, ShowTitles must be enabled.
func (sf *TBarChart) ShowMarks() bool {
	return sf.showMarks
}

// SetShowMarks turns on and off marks under horizontal axis
func (sf *TBarChart) SetShowMarks(show bool) {
	sf.block.Lock()
	defer sf.block.Unlock()

	sf.showMarks = show
}
