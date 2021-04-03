// Package main -- пример для демонстрации goTV
package main

import (
	"github.com/prospero78/goTV/tv"
)

func main() {
	tv.InitLibrary()
	defer tv.DeinitLibrary()
	view := tv.AddWindow(0, 0, 10, 7, "Привет мир!", false, false)
	btnQuit := tv.CreateButton(view, 15, 4, "Hi", 1, false, false)
	btnQuit.OnClick(func(ev tv.Event) {
		go tv.Stop()
	})
	tv.MainLoop()
}
