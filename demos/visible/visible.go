package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func main() {
	tv.InitLibrary()
	defer tv.DeinitLibrary()

	view := tv.AddWindow(0, 0, 10, 7, "Hello World!")
	view.SetPack(cons.Vertical)

	frmResize := tv.CreateFrame(view, 8, 6, cons.BorderNone, cons.Fixed)
	frmResize.SetTitle("FrameTop")
	frmResize.SetPack(cons.Horizontal)
	btn1 := tv.CreateButton(frmResize, 8, 5, "Button 1", 1)
	btn2 := tv.CreateButton(frmResize, 8, 5, "Button 2", 1)
	btn3 := tv.CreateButton(frmResize, 8, 5, "Button 3", 1)

	frmBtns := tv.CreateFrame(view, 8, 5, cons.BorderNone, cons.Fixed)
	frmBtns.SetPack(cons.Horizontal)
	frmBtns.SetTitle("FrameBottom")

	btnHide1 := tv.CreateButton(frmBtns, 8, 4, "Hide 1", 1)
	btnHide1.OnClick(func(ev event.TEvent) {
		if btn1.IsVisible() {
			btnHide1.SetTitle("Show 1")
			tv.ActivateControl(view, btn1)
			btn1.SetVisible(false)
		} else {
			btnHide1.SetTitle("Hide 1")
			btn1.SetVisible(true)
		}
	})
	btnHide2 := tv.CreateButton(frmBtns, 8, 4, "Hide 2", 1)
	btnHide2.OnClick(func(ev event.TEvent) {
		if btn2.IsVisible() {
			btnHide2.SetTitle("Show 2")
			tv.ActivateControl(view, btn2)
			btn2.SetVisible(false)
		} else {
			btnHide2.SetTitle("Hide 2")
			btn2.SetVisible(true)
		}
	})
	btnHide3 := tv.CreateButton(frmBtns, 8, 4, "Hide 3", 1)
	btnHide3.OnClick(func(ev event.TEvent) {
		if btn3.IsVisible() {
			btnHide3.SetTitle("Show 3")
			tv.ActivateControl(view, btn3)
			btn3.SetVisible(false)
		} else {
			btnHide3.SetTitle("Hide 3")
			btn3.SetVisible(true)
		}
	})

	btnQuit := tv.CreateButton(frmBtns, 8, 4, "Quit", 1)
	btnQuit.OnClick(func(ev event.TEvent) {
		go tv.Stop()
	})

	tv.MainLoop()
}
