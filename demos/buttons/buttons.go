package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
)

func createView() {
	view := tv.AddWindow(0, 0, 10, 7, "Button` Demo")
	view.SetTitleButtons(cons.ButtonMaximize | cons.ButtonClose)

	frmViews := tv.CreateFrame(view, 8, 5, cons.BorderNone, cons.Fixed)
	frmViews.SetPack(cons.Horizontal)
	frmFull := tv.CreateFrame(frmViews, 8, 5, cons.BorderThin, cons.Fixed)
	frmFull.SetPack(cons.Vertical)
	frmFull.SetTitle("Full")
	frmHalf := tv.CreateFrame(frmViews, 8, 5, cons.BorderThin, cons.Fixed)
	frmHalf.SetPack(cons.Vertical)
	frmHalf.SetTitle("Half")
	frmNone := tv.CreateFrame(frmViews, 8, 5, cons.BorderThin, cons.Fixed)
	frmNone.SetPack(cons.Vertical)
	frmNone.SetTitle("None")

	btnF1 := tv.CreateButton(frmFull, cons.AutoSize, 4, "First", cons.Fixed)
	btnF2 := tv.CreateButton(frmFull, cons.AutoSize, 4, "Second", cons.Fixed)
	btnF3 := tv.CreateButton(frmFull, cons.AutoSize, 4, "Quit", cons.Fixed)
	btnF1.SetShadowType(cons.ShadowFull)
	btnF2.SetShadowType(cons.ShadowFull)
	btnF3.SetShadowType(cons.ShadowFull)
	btnH1 := tv.CreateButton(frmHalf, cons.AutoSize, 4, "First", cons.Fixed)
	btnH2 := tv.CreateButton(frmHalf, cons.AutoSize, 4, "Second", cons.Fixed)
	btnH3 := tv.CreateButton(frmHalf, cons.AutoSize, 4, "Quit", cons.Fixed)
	btnH1.SetShadowType(cons.ShadowHalf)
	btnH2.SetShadowType(cons.ShadowHalf)
	btnH3.SetShadowType(cons.ShadowHalf)
	btnN1 := tv.CreateButton(frmNone, cons.AutoSize, 4, "First", cons.Fixed)
	btnN2 := tv.CreateButton(frmNone, cons.AutoSize, 4, "Second", cons.Fixed)
	btnN3 := tv.CreateButton(frmNone, cons.AutoSize, 4, "Quit", cons.Fixed)
	btnN1.SetShadowType(cons.ShadowNone)
	btnN2.SetShadowType(cons.ShadowNone)
	btnN3.SetShadowType(cons.ShadowNone)

	btnF3.OnClick(func(ev cons.Event) {
		go tv.Stop()
	})
	btnH3.OnClick(func(ev cons.Event) {
		go tv.Stop()
	})
	btnN3.OnClick(func(ev cons.Event) {
		go tv.Stop()
	})

	tv.ActivateControl(view, btnF1)
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	tv.InitLibrary()
	defer tv.DeinitLibrary()

	createView()

	// start event processing loop - the main core of the library
	tv.MainLoop()
}

func main() {
	mainLoop()
}
