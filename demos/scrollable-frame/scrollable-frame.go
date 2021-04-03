package main

import (
	"fmt"

	ui "github.com/prospero78/goTV/tv"
)

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	wnd := ui.AddWindow(0, 0, 60, 0, "Scrollable frame", false, true)
	wnd.SetSizable(false)

	frm := ui.CreateFrame(wnd, 50, 12, ui.BorderNone, ui.Fixed)
	frm.SetPack(ui.Vertical)
	frm.SetScrollable(true)

	for i := 0; i < 10; i++ {
		label := fmt.Sprintf("Button %d - press to quit", i)
		btn := ui.CreateButton(frm, 40, ui.AutoSize, label, 1)

		btn.OnClick(func(ev ui.Event) {
			go ui.Stop()
		})
	}

	ui.MainLoop()
}
