package main

import (
	ui "github.com/prospero78/goTV/tv"
)

func createView() {
	view := ui.AddWindow(0, 0, 10, 7, "EditField Demo", false, false)
	view.SetTitleButtons(ui.ButtonMaximize | ui.ButtonClose)

	frmChk := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmChk.SetPack(ui.Vertical)
	frmChk.SetPaddings(1, 1)
	frmChk.SetGaps(1, 1)
	ui.CreateLabel(frmChk, 0, 0, "Enter password:", ui.Fixed)
	edFld := ui.CreateEditField(frmChk, 20, "", ui.Fixed)
	edFld.SetPasswordMode(true)
	chkPass := ui.CreateCheckBox(frmChk, 0, "Show Password", ui.Fixed, true)

	ui.ActivateControl(view, edFld)

	chkPass.OnChange(func(state int) {
		if state == 1 {
			edFld.SetPasswordMode(false)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
		} else if state == 0 {
			edFld.SetPasswordMode(true)
			ui.PutEvent(ui.Event{Type: ui.EventRedraw})
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
