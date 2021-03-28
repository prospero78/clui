package main

import (
	ui "github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/types"
)

func createView() {
	view := ui.AddWindow(0, 0, 10, 7, "EditField Demo")
	view.SetTitleButtons(ui.ButtonMaximize | ui.ButtonClose)

	frmChk := ui.CreateFrame(view, 8, 5, ui.BorderNone, types.Fixed)
	frmChk.SetPack(ui.Vertical)
	frmChk.SetPaddings(1, 1)
	frmChk.SetGaps(1, 1)
	ui.CreateLabel(frmChk, types.AutoSize, types.AutoSize, "Enter password:", types.Fixed)
	edFld := ui.CreateEditField(frmChk, 20, "", types.Fixed)
	edFld.SetPasswordMode(true)
	chkPass := ui.CreateCheckBox(frmChk, types.AutoSize, "Show Password", types.Fixed)

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
