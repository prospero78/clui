// Package main -- пример для демонстрации goTV
package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func main() {
	tv.InitLibrary()
	defer tv.DeinitLibrary()
	view := tv.AddWindow(0, 0, 10, 7, "Привет мир!")
	btnQuit := tv.CreateButton(view, 15, 4, "Hi", 1)
	btnQuit.OnClick(func(ev event.TEvent) {
		go tv.Stop()
	})
	tv.MainLoop()
}
