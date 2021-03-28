package main

import (
	"fmt"

	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func main() {
	tv.InitLibrary()
	defer tv.DeinitLibrary()

	wnd := tv.AddWindow(0, 0, 60, cons.AutoSize, "Scrollable frame")
	wnd.SetSizable(false)

	frm := tv.CreateFrame(wnd, 50, 12, cons.BorderNone, cons.Fixed)
	frm.SetPack(cons.Vertical)
	frm.SetScrollable(true)

	for i := 0; i < 10; i++ {
		label := fmt.Sprintf("Button %d - press to quit", i)
		btn := tv.CreateButton(frm, 40, cons.AutoSize, label, 1)

		btn.OnClick(func(ev event.TEvent) {
			go tv.Stop()
		})
	}

	tv.MainLoop()
}
