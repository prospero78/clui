package main

import (
	"github.com/prospero78/goTV/tv"
)

func customColored(d *tv.BarDataCell) {
	part := d.TotalMax / 3
	if d.ID%2 == 0 {
		switch {
		case d.Value <= part:
			d.Fg = tv.ColorGreen
		case d.Value > 2*part:
			d.Fg = tv.ColorRed
		default:
			d.Fg = tv.ColorBlue
		}
	} else {
		d.Ch = '#'
		switch {
		case d.Value <= part:
			d.Fg = tv.ColorGreenBold
		case d.Value > 2*part:
			d.Fg = tv.ColorRedBold
		default:
			d.Fg = tv.ColorBlueBold
		}
	}
}

func createView() *tv.BarChart {

	view := tv.AddWindow(0, 0, 10, 7, "BarChart Demo")
	bch := tv.CreateBarChart(view, 40, 12, 1)

	frmChk := tv.CreateFrame(view, 8, 5, tv.BorderNone, tv.Fixed)
	frmChk.SetPack(tv.Vertical)
	chkTitles := tv.CreateCheckBox(frmChk, tv.AutoSize, "Show Titles", tv.Fixed)
	chkMarks := tv.CreateCheckBox(frmChk, tv.AutoSize, "Show Marks", tv.Fixed)
	chkTitles.SetState(1)
	chkLegend := tv.CreateCheckBox(frmChk, tv.AutoSize, "Show Legend", tv.Fixed)
	chkValues := tv.CreateCheckBox(frmChk, tv.AutoSize, "Show Values", tv.Fixed)
	chkValues.SetState(1)
	chkFixed := tv.CreateCheckBox(frmChk, tv.AutoSize, "Fixed Width", tv.Fixed)
	chkGap := tv.CreateCheckBox(frmChk, tv.AutoSize, "No Gap", tv.Fixed)
	chkMulti := tv.CreateCheckBox(frmChk, tv.AutoSize, "MultiColored", tv.Fixed)
	chkCustom := tv.CreateCheckBox(frmChk, tv.AutoSize, "Custom Colors", tv.Fixed)

	tv.ActivateControl(view, chkTitles)

	chkTitles.OnChange(func(state int) {
		if state == 0 {
			chkMarks.SetEnabled(false)
			bch.SetShowTitles(false)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		} else if state == 1 {
			chkMarks.SetEnabled(true)
			bch.SetShowTitles(true)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		}
	})
	chkMarks.OnChange(func(state int) {
		if state == 0 {
			bch.SetShowMarks(false)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		} else if state == 1 {
			bch.SetShowMarks(true)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		}
	})
	chkLegend.OnChange(func(state int) {
		if state == 0 {
			bch.SetLegendWidth(0)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		} else if state == 1 {
			bch.SetLegendWidth(10)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		}
	})
	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		} else if state == 1 {
			bch.SetValueWidth(5)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		}
	})
	chkMulti.OnChange(func(state int) {
		if state == 0 {
			d := []tv.BarData{
				{Value: 80, Title: "80%"},
				{Value: 50, Title: "50%"},
				{Value: 150, Title: ">100%"},
			}
			bch.SetData(d)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		} else if state == 1 {
			d := []tv.BarData{
				{Value: 80, Title: "80%", Fg: tv.ColorBlue},
				{Value: 50, Title: "50%", Fg: tv.ColorGreen, Ch: 'X'},
				{Value: 150, Title: ">100%", Fg: tv.ColorYellow},
			}
			bch.SetData(d)
			tv.PutEvent(tv.Event{Type: tv.EventRedraw})
		}
	})
	chkFixed.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoSize(true)
		} else if state == 1 {
			bch.SetAutoSize(false)
		}
		tv.PutEvent(tv.Event{Type: tv.EventRedraw})
	})
	chkGap.OnChange(func(state int) {
		if state == 1 {
			bch.SetBarGap(0)
		} else if state == 0 {
			bch.SetBarGap(1)
		}
		tv.PutEvent(tv.Event{Type: tv.EventRedraw})
	})
	chkCustom.OnChange(func(state int) {
		if state == 0 {
			bch.OnDrawCell(nil)
		} else if state == 1 {
			bch.OnDrawCell(customColored)
		}
		tv.PutEvent(tv.Event{Type: tv.EventRedraw})
	})

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	tv.InitLibrary()
	defer tv.DeinitLibrary()

	b := createView()
	b.SetBarGap(1)
	d := []tv.BarData{
		{Value: 80, Title: "80%"},
		{Value: 50, Title: "50%"},
		{Value: 150, Title: ">100%"},
	}
	b.SetData(d)
	b.SetValueWidth(5)
	b.SetAutoSize(true)

	// start event processing loop - the main core of the library
	tv.MainLoop()
}

func main() {
	mainLoop()
}
