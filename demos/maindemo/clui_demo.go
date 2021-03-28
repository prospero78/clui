/*
Demo includes:
    - How to intialize and run the application
    - How to stop the application
    - How to use Control's events (Button ones)
    - How to change theme on the fly
    - How to use dialogs
    - How to make composer refresh the screen
    - How to intercept Enter key(term.KeyCtrlM) in EditField(ListBox is the same)
*/
package main

import (
	"fmt"
	"strconv"

	term "github.com/nsf/termbox-go"
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func updateProgress(value string, pb *tv.ProgressBar) {
	v, _ := strconv.Atoi(value)
	pb.SetValue(v)
}

func changeTheme(lb *tv.ListBox, btn *tv.Button, tp int) {
	items := tv.ThemeNames()
	dlgType := cons.SelectDialogRadio
	if tp == 1 {
		dlgType = cons.SelectDialogList
	}

	curr := -1
	for i, tName := range items {
		if tName == tv.CurrentTheme() {
			curr = i
			break
		}
	}

	selDlg := tv.CreateSelectDialog("Choose a theme", items, curr, dlgType)
	selDlg.OnClose(func() {
		if selDlg.Result() == cons.DialogButton1 {
			idx := selDlg.Value()
			lb.AddItem(fmt.Sprintf("Selected item: %v", selDlg.Value()))
			lb.SelectItem(lb.ItemCount() - 1)
			if idx != -1 {
				tv.SetCurrentTheme(items[idx])
			}
		}

		btn.SetEnabled(true)
		// ask the composer to repaint all windows
		tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
	})
}

func createView() {

	view := tv.AddWindow(0, 0, 20, 7, "Theme Manager Demo")

	frmLeft := tv.CreateFrame(view, 8, 4, cons.BorderNone, 1)
	frmLeft.SetPack(cons.Vertical)
	frmLeft.SetGaps(cons.KeepValue, 1)
	frmLeft.SetPaddings(1, 1)

	frmTheme := tv.CreateFrame(frmLeft, 8, 1, cons.BorderNone, cons.Fixed)
	frmTheme.SetGaps(1, cons.KeepValue)
	checkBox := tv.CreateCheckBox(frmTheme, cons.AutoSize, "Use ListBox", cons.Fixed)
	btnTheme := tv.CreateButton(frmTheme, cons.AutoSize, 4, "Select theme", cons.Fixed)
	tv.CreateFrame(frmLeft, 1, 1, cons.BorderNone, 1)

	frmPb := tv.CreateFrame(frmLeft, 8, 1, cons.BorderNone, cons.Fixed)
	tv.CreateLabel(frmPb, 1, 1, "[", cons.Fixed)
	pb := tv.CreateProgressBar(frmPb, 20, 1, 1)
	pb.SetLimits(0, 10)
	pb.SetTitle("{{value}} of {{max}}")
	tv.CreateLabel(frmPb, 1, 1, "]", cons.Fixed)

	edit := tv.CreateEditField(frmLeft, 5, "0", cons.Fixed)

	frmEdit := tv.CreateFrame(frmLeft, 8, 1, cons.BorderNone, cons.Fixed)
	frmEdit.SetPaddings(1, 1)
	frmEdit.SetGaps(1, cons.KeepValue)
	btnSet := tv.CreateButton(frmEdit, cons.AutoSize, 4, "Set", cons.Fixed)
	btnStep := tv.CreateButton(frmEdit, cons.AutoSize, 4, "Step", cons.Fixed)
	tv.CreateFrame(frmEdit, 1, 1, cons.BorderNone, 1)
	btnQuit := tv.CreateButton(frmEdit, cons.AutoSize, 4, "Quit", cons.Fixed)

	logBox := tv.CreateListBox(view, 28, 5, cons.Fixed)

	tv.ActivateControl(view, edit)

	edit.OnKeyPress(func(key term.Key, ch rune) bool {
		if key == term.KeyCtrlM {
			v := edit.Title()
			logBox.AddItem(fmt.Sprintf("New PB value(KeyPress): %v", v))
			logBox.SelectItem(logBox.ItemCount() - 1)
			updateProgress(v, pb)
			return true
		}
		return false
	})
	btnTheme.OnClick(func(ev event.TEvent) {
		btnTheme.SetEnabled(false)
		tp := checkBox.State()
		changeTheme(logBox, btnTheme, tp)
	})
	btnSet.OnClick(func(ev event.TEvent) {
		v := edit.Title()
		logBox.AddItem(fmt.Sprintf("New ProgressBar value: %v", v))
		logBox.SelectItem(logBox.ItemCount() - 1)
		updateProgress(v, pb)
	})
	btnStep.OnClick(func(ev event.TEvent) {
		go pb.Step()
		logBox.AddItem("ProgressBar step")
		logBox.SelectItem(logBox.ItemCount() - 1)
		tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
	})
	btnQuit.OnClick(func(ev event.TEvent) {
		go tv.Stop()
	})
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	tv.InitLibrary()
	defer tv.DeinitLibrary()

	tv.SetThemePath("themes")

	createView()

	// start event processing loop - the main core of the library
	tv.MainLoop()
}

func main() {
	mainLoop()
}
