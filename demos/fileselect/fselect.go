package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/types"
)

func createView() {
	view := tv.AddWindow(0, 0, 30, 7, "File select")
	view.SetPack(tv.Vertical)
	view.SetGaps(0, 1)
	view.SetPaddings(2, 2)

	frmPath := tv.CreateFrame(view, 1, 1, tv.BorderNone, types.Fixed)
	frmPath.SetPack(tv.Horizontal)
	tv.CreateLabel(frmPath, types.AutoSize, types.AutoSize, "Initial path", types.Fixed)
	edPath := tv.CreateEditField(frmPath, 16, "", 1)

	frmMask := tv.CreateFrame(view, 1, 1, tv.BorderNone, types.Fixed)
	frmMask.SetPack(tv.Horizontal)
	tv.CreateLabel(frmMask, types.AutoSize, types.AutoSize, "File masks", types.Fixed)
	edMasks := tv.CreateEditField(frmMask, 16, "*", 1)

	frmOpts := tv.CreateFrame(view, 1, 1, tv.BorderNone, types.Fixed)
	frmOpts.SetPack(tv.Horizontal)
	cbDir := tv.CreateCheckBox(frmOpts, types.AutoSize, "Select directory", types.Fixed)
	cbMust := tv.CreateCheckBox(frmOpts, types.AutoSize, "Must exists", types.Fixed)
	tv.CreateFrame(frmOpts, 1, 1, tv.BorderNone, 1)

	lblSelected := tv.CreateLabel(view, 30, 5, "Selected:", types.Fixed)
	lblSelected.SetMultiline(true)

	frmBtns := tv.CreateFrame(view, 1, 1, tv.BorderNone, types.Fixed)
	frmBtns.SetPack(tv.Horizontal)
	btnSet := tv.CreateButton(frmBtns, types.AutoSize, 4, "Select", types.Fixed)
	btnQuit := tv.CreateButton(frmBtns, types.AutoSize, 4, "Quit", types.Fixed)
	tv.CreateFrame(frmBtns, 1, 1, tv.BorderNone, 1)

	tv.ActivateControl(view, edMasks)

	btnSet.OnClick(func(ev tv.Event) {
		s := "Select "
		if cbDir.State() == 1 {
			s += "directory"
		} else {
			s += "file"
		}
		if cbMust.State() == 1 {
			s += "[X]"
		}
		dlg := tv.CreateFileSelectDialog(
			s,
			edMasks.Title(),
			edPath.Title(),
			cbDir.State() == 1,
			cbMust.State() == 1)
		dlg.OnClose(func() {
			if !dlg.Selected {
				lblSelected.SetTitle("Selected:\nNothing")
				return
			}

			var lb string
			if dlg.Exists {
				lb = "Selected existing"
			} else {
				lb = "Create new"
			}

			if cbDir.State() == 0 {
				lb += " file:\n"
			} else {
				lb += " directory:\n"
			}

			lb += dlg.FilePath
			lblSelected.SetTitle(lb)
		})
	})

	btnQuit.OnClick(func(ev tv.Event) {
		go tv.Stop()
	})
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
