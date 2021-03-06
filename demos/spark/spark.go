package main

import (
	"math/rand"
	"time"

	ui "github.com/prospero78/goTV/tv"
)

func createView() *ui.SparkChart {

	view := ui.AddWindow(0, 0, 10, 7, "BarChart Demo", false, false)
	bch := ui.CreateSparkChart(view, 25, 12, 1)
	bch.SetTop(20)

	frmChk := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmChk.SetPack(ui.Vertical)
	chkValues := ui.CreateCheckBox(frmChk, 0, "Show Values", ui.Fixed, true)
	chkValues.SetState(0)
	chkHilite := ui.CreateCheckBox(frmChk, 0, "Hilite peaks", ui.Fixed, true)
	chkHilite.SetState(1)
	chkAuto := ui.CreateCheckBox(frmChk, 0, "Auto scale", ui.Fixed, true)
	chkAuto.SetState(1)

	ui.ActivateControl(view, chkValues)

	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
		} else if state == 1 {
			bch.SetValueWidth(5)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	chkHilite.OnChange(func(state int) {
		if state == 0 {
			bch.SetHilitePeaks(false)
		} else if state == 1 {
			bch.SetHilitePeaks(true)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})
	chkAuto.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoScale(false)
		} else if state == 1 {
			bch.SetAutoScale(true)
		}
		ui.PutEvent(ui.Event{Type: ui.EventRedraw})
	})

	return bch
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	b := createView()
	b.SetData([]float64{1, 2, 3, 4, 5, 6, 6, 7, 5, 8, 9})

	ticker := time.NewTicker(time.Millisecond * 200).C
	go func() {
		for range ticker {
			b.AddData(float64(rand.Int31n(20))) //nolint
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		}
	}()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
