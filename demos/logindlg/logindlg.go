package main

import (
	ui "github.com/prospero78/goTV/tv"
)

func createView() {
	view := ui.AddWindow(0, 0, 30, 7, "Login dialog", false, false)
	view.SetPack(ui.Vertical)
	view.SetGaps(0, 1)
	view.SetPaddings(2, 2)

	frmOpts := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmOpts.SetPack(ui.Horizontal)
	cbCheck := ui.CreateCheckBox(frmOpts, 0, "Use callback to test data", ui.Fixed, true)

	ui.CreateLabel(view, 0, 0, "Correct credentials", ui.Fixed)

	frmCreds := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmCreds.SetPack(ui.Horizontal)
	frmCreds.SetGaps(1, 0)
	ui.CreateLabel(frmCreds, 0, 0, "Username", ui.Fixed)
	edUser := ui.CreateEditField(frmCreds, 8, "", 1)
	ui.CreateLabel(frmCreds, 0, 0, "Password", ui.Fixed)
	edPass := ui.CreateEditField(frmCreds, 8, "", 1)

	lbRes := ui.CreateLabel(view, 0, 0, "Result:", ui.Fixed)

	frmBtns := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Fixed)
	frmBtns.SetPack(ui.Horizontal)
	btnDlg := ui.CreateButton(frmBtns, 0, 4, "Login", ui.Fixed, true, false)
	btnQuit := ui.CreateButton(frmBtns, 0, 4, "Quit", ui.Fixed, true, false)
	ui.CreateFrame(frmBtns, 1, 1, ui.BorderNone, 1)

	ui.ActivateControl(view, edUser)

	btnDlg.OnClick(func(ev ui.Event) {
		dlg := ui.CreateLoginDialog(
			"Enter credentials",
			edUser.Title(),
		)

		if cbCheck.State() == 1 {
			dlg.OnCheck(func(u, p string) bool {
				return u == edUser.Title() && p == edPass.Title()
			})
		} else {
			dlg.OnCheck(nil)
		}

		dlg.OnClose(func() {
			if dlg.Action == ui.LoginCanceled {
				lbRes.SetTitle("Result:\nDialog canceled")
				return
			}

			if dlg.Action == ui.LoginInvalid {
				lbRes.SetTitle("Result:\nInvalid username or password")
				return
			}

			if dlg.Action == ui.LoginOk {
				if cbCheck.State() == 1 {
					lbRes.SetTitle("Result:\nLogged in successfully")
				} else {
					lbRes.SetTitle("Result:\nEntered [" + dlg.Username + ":" + dlg.Password + "]")
				}
				return
			}
		})
	})

	btnQuit.OnClick(func(ev ui.Event) {
		go ui.Stop()
	})
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	createView()

	// start event processing loop - the main core of the library
	ui.MainLoop()
}

func main() {
	mainLoop()
}
