package main

import (
	ui "github.com/prospero78/goTV/tv"
)

func createView() {
	view := ui.AddWindow(0, 0, 10, 7, "Button` Demo", false, false)
	view.SetTitleButtons(ui.ButtonMaximize | ui.ButtonClose)

	frmViews := ui.CreateFrame(view, 8, 5, ui.BorderNone, ui.Fixed)
	frmViews.SetPack(ui.Horizontal)
	frmFull := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
	frmFull.SetPack(ui.Vertical)
	frmFull.SetTitle("Full")
	frmHalf := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
	frmHalf.SetPack(ui.Vertical)
	frmHalf.SetTitle("Half")
	frmNone := ui.CreateFrame(frmViews, 8, 5, ui.BorderThin, ui.Fixed)
	frmNone.SetPack(ui.Vertical)
	frmNone.SetTitle("None")

	btnF1 := ui.CreateButton(frmFull, 0, 4, "First", ui.Fixed, true, false)
	btnF2 := ui.CreateButton(frmFull, 0, 4, "Second", ui.Fixed, true, false)
	btnF3 := ui.CreateButton(frmFull, 0, 4, "Quit", ui.Fixed, true, false)
	btnF1.SetShadowType(ui.ShadowFull)
	btnF2.SetShadowType(ui.ShadowFull)
	btnF3.SetShadowType(ui.ShadowFull)
	btnH1 := ui.CreateButton(frmHalf, 0, 4, "First", ui.Fixed, true, false)
	btnH2 := ui.CreateButton(frmHalf, 0, 4, "Second", ui.Fixed, true, false)
	btnH3 := ui.CreateButton(frmHalf, 0, 4, "Quit", ui.Fixed, true, false)
	btnH1.SetShadowType(ui.ShadowHalf)
	btnH2.SetShadowType(ui.ShadowHalf)
	btnH3.SetShadowType(ui.ShadowHalf)
	btnN1 := ui.CreateButton(frmNone, 0, 4, "First", ui.Fixed, true, false)
	btnN2 := ui.CreateButton(frmNone, 0, 4, "Second", ui.Fixed, true, false)
	btnN3 := ui.CreateButton(frmNone, 0, 4, "Quit", ui.Fixed, true, false)
	btnN1.SetShadowType(ui.ShadowNone)
	btnN2.SetShadowType(ui.ShadowNone)
	btnN3.SetShadowType(ui.ShadowNone)

	btnF3.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
	btnH3.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
	btnN3.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})

	ui.ActivateControl(view, btnF1)
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
