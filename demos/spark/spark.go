package main

import (
	"math/rand"
	"time"

	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func createView() *tv.SparkChart {

	view := tv.AddWindow(0, 0, 10, 7, "BarChart Demo")
	bch := tv.CreateSparkChart(view, 25, 12, 1)
	bch.SetTop(20)

	frmChk := tv.CreateFrame(view, 8, 5, cons.ABorderNone, cons.Fixed)
	frmChk.SetPack(cons.Vertical)
	chkValues := tv.CreateCheckBox(frmChk, cons.AutoSize, "Show Values", cons.Fixed)
	chkValues.SetState(0)
	chkHilite := tv.CreateCheckBox(frmChk, cons.AutoSize, "Hilite peaks", cons.Fixed)
	chkHilite.SetState(1)
	chkAuto := tv.CreateCheckBox(frmChk, cons.AutoSize, "Auto scale", cons.Fixed)
	chkAuto.SetState(1)

	tv.ActivateControl(view, chkValues)

	chkValues.OnChange(func(state int) {
		if state == 0 {
			bch.SetValueWidth(0)
		} else if state == 1 {
			bch.SetValueWidth(5)
		}
		tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
	})
	chkHilite.OnChange(func(state int) {
		if state == 0 {
			bch.SetHilitePeaks(false)
		} else if state == 1 {
			bch.SetHilitePeaks(true)
		}
		tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
	})
	chkAuto.OnChange(func(state int) {
		if state == 0 {
			bch.SetAutoScale(false)
		} else if state == 1 {
			bch.SetAutoScale(true)
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
	b.SetData([]float64{1, 2, 3, 4, 5, 6, 6, 7, 5, 8, 9})

	ticker := time.NewTicker(time.Millisecond * 200).C
	go func() {
		for {
			select {
			case <-ticker:
				b.AddData(float64(rand.Int31n(20)))
				tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
			}
		}
	}()

	// start event processing loop - the main core of the library
	tv.MainLoop()
}

func main() {
	mainLoop()
}
