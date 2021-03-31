package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/barchart/bardatacell"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func customColored(d *bardatacell.TBarDataCell) {
	part := d.TotalMax / 3
	if d.ID %2 == 0 {
		switch {
		case d.Value <= part:
			d.Fg = cons.ColorGreen
		case d.Value > 2*part:
			d.Fg = cons.ColorRed
		default:
			d.Fg = cons.ColorBlue
		}
	} else {
		d.Ch = '#'
		switch {
		case d.Value <= part:
			d.Fg = cons.ColorGreenBold
		case d.Value > 2*part:
			d.Fg = cons.ColorRedBold
		default:
			d.Fg = cons.ColorBlueBold
		}
	}
}

func createView() *tv.TBarChart {

	view := tv.AddWindow(0, 0, 10, 7, "BarChart Demo")
	bch := tv.CreateBarChart(view, 40, 12, 1)

	frmChk := tv.CreateFrame(view, 8, 5, cons.BorderNone, cons.Fixed)
	frmChk.SetPack(cons.Vertical)
	chkTitles := tv.CreateCheckBox(frmChk, cons.AutoSize, "Show Titles", cons.Fixed)
	chkMarks := tv.CreateCheckBox(frmChk, cons.AutoSize, "Show Marks", cons.Fixed)
	chkTitles.SetState(1)
	chkLegend := tv.CreateCheckBox(frmChk, cons.AutoSize, "Show Legend", cons.Fixed)
	chkValues := tv.CreateCheckBox(frmChk, cons.AutoSize, "Show Values", cons.Fixed)
	chkValues.SetState(1)
	chkFixed := tv.CreateCheckBox(frmChk, cons.AutoSize, "Fixed Width", cons.Fixed)
	chkGap := tv.CreateCheckBox(frmChk, cons.AutoSize, "No Gap", cons.Fixed)
	chkMulti := tv.CreateCheckBox(frmChk, cons.AutoSize, "MultiColored", cons.Fixed)
	chkCustom := tv.CreateCheckBox(frmChk, cons.AutoSize, "Custom Colors", cons.Fixed)

	tv.ActivateControl(view, chkTitles)

	chkTitles.OnChange(func(state int) {
		if state == 0 {
			chkMarks.SetEnabled(false)
			bch.SetShowTitles(false)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		} else if state == 1 {
			chkMarks.SetEnabled(true)
			bch.SetShowTitles(true)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		}
	})
	chkMarks.OnChange(func(state int) {
		if state == 0 {
			bch.SetShowMarks(false)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		} else if state == 1 {
			bch.SetShowMarks(true)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		}
	})
	chkLegend.OnChange(func(state int) {
		if state == 0 {
			bch.SetLegendWidth(0)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		} else if state == 1 {
			bch.SetLegendWidth(10)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		}
	})
	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		} else if state == 1 {
			bch.SetValueWidth(5)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
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
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		} else if state == 1 {
			d := []tv.BarData{
				{Value: 80, Title: "80%", Fg: cons.ColorBlue},
				{Value: 50, Title: "50%", Fg: cons.ColorGreen, Ch: 'X'},
				{Value: 150, Title: ">100%", Fg: cons.ColorYellow},
			}
			bch.SetData(d)
			tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
		}
	})
	chkFixed.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoSize(true)
		} else if state == 1 {
			bch.SetAutoSize(false)
		}
		tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
	})
	chkGap.OnChange(func(state int) {
		if state == 1 {
			bch.SetBarGap(0)
		} else if state == 0 {
			bch.SetBarGap(1)
		}
		tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
	})
	chkCustom.OnChange(func(state int) {
		if state == 0 {
			bch.OnDrawCell(nil)
		} else if state == 1 {
			bch.OnDrawCell(customColored)
		}
		tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
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
