package main

import (
	ui "github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func createView() {
	view := ui.AddWindow(0, 0, 10, 7, "EditField Demo")
	view.SetTitleButtons(cons.ButtonMaximize | cons.ButtonClose)

	frmChk := ui.CreateFrame(view, 8, 5, cons.BorderNone, cons.Fixed)
	frmChk.SetPack(cons.Vertical)
	frmChk.SetPaddings(1, 1)
	frmChk.SetGaps(1, 1)
	ui.CreateLabel(frmChk, cons.AutoSize, cons.AutoSize, "Enter password:", cons.Fixed)
	edFld := ui.CreateEditField(frmChk, 20, "", cons.Fixed)
	edFld.SetPasswordMode(true)
	chkPass := ui.CreateCheckBox(frmChk, cons.AutoSize, "Show Password", cons.Fixed)

	ui.ActivateControl(view, edFld)

	chkPass.OnChange(func(state int) {
		if state == 1 {
			edFld.SetPasswordMode(false)
			ui.PutEvent(event.TEvent{Type: cons.EventRedraw})
		} else if state == 0 {
			edFld.SetPasswordMode(true)
			ui.PutEvent(event.TEvent{Type: cons.EventRedraw})
		}
	})
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	createView()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
