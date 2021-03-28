package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/types"
)

func createView() {
	view := tv.AddWindow(0, 0, 10, 7, "Button` Demo")
	view.SetTitleButtons(tv.ButtonMaximize | tv.ButtonClose)

	frmViews := tv.CreateFrame(view, 8, 5, tv.BorderNone, types.Fixed)
	frmViews.SetPack(tv.Horizontal)
	frmFull := tv.CreateFrame(frmViews, 8, 5, tv.BorderThin, types.Fixed)
	frmFull.SetPack(tv.Vertical)
	frmFull.SetTitle("Full")
	frmHalf := tv.CreateFrame(frmViews, 8, 5, tv.BorderThin, types.Fixed)
	frmHalf.SetPack(tv.Vertical)
	frmHalf.SetTitle("Half")
	frmNone := tv.CreateFrame(frmViews, 8, 5, tv.BorderThin, types.Fixed)
	frmNone.SetPack(tv.Vertical)
	frmNone.SetTitle("None")

	btnF1 := tv.CreateButton(frmFull, types.AutoSize, 4, "First", types.Fixed)
	btnF2 := tv.CreateButton(frmFull, types.AutoSize, 4, "Second", types.Fixed)
	btnF3 := tv.CreateButton(frmFull, types.AutoSize, 4, "Quit", types.Fixed)
	btnF1.SetShadowType(tv.ShadowFull)
	btnF2.SetShadowType(tv.ShadowFull)
	btnF3.SetShadowType(tv.ShadowFull)
	btnH1 := tv.CreateButton(frmHalf, types.AutoSize, 4, "First", types.Fixed)
	btnH2 := tv.CreateButton(frmHalf, types.AutoSize, 4, "Second", types.Fixed)
	btnH3 := tv.CreateButton(frmHalf, types.AutoSize, 4, "Quit", types.Fixed)
	btnH1.SetShadowType(tv.ShadowHalf)
	btnH2.SetShadowType(tv.ShadowHalf)
	btnH3.SetShadowType(tv.ShadowHalf)
	btnN1 := tv.CreateButton(frmNone, types.AutoSize, 4, "First", types.Fixed)
	btnN2 := tv.CreateButton(frmNone, types.AutoSize, 4, "Second", types.Fixed)
	btnN3 := tv.CreateButton(frmNone, types.AutoSize, 4, "Quit", types.Fixed)
	btnN1.SetShadowType(tv.ShadowNone)
	btnN2.SetShadowType(tv.ShadowNone)
	btnN3.SetShadowType(tv.ShadowNone)

	btnF3.OnClick(func(ev tv.Event) {
		go tv.Stop()
	})
	btnH3.OnClick(func(ev tv.Event) {
		go tv.Stop()
	})
	btnN3.OnClick(func(ev tv.Event) {
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
