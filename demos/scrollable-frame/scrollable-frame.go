package main

import (
	"fmt"

	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/types"
)

func main() {
	tv.InitLibrary()
	defer tv.DeinitLibrary()

	wnd := tv.AddWindow(0, 0, 60, types.AutoSize, "Scrollable frame")
	wnd.SetSizable(false)

	frm := tv.CreateFrame(wnd, 50, 12, tv.BorderNone, types.Fixed)
	frm.SetPack(tv.Vertical)
	frm.SetScrollable(true)

	for i := 0; i < 10; i++ {
		label := fmt.Sprintf("Button %d - press to quit", i)
		btn := tv.CreateButton(frm, 40, types.AutoSize, label, 1)

		btn.OnClick(func(ev tv.Event) {
			go tv.Stop()
		})
	}

	tv.MainLoop()
}
