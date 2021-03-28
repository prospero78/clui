package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/types"
)

func createView() {
	view := tv.AddWindow(0, 0, 30, 7, "Login dialog")
	view.SetPack(tv.Vertical)
	view.SetGaps(0, 1)
	view.SetPaddings(2, 2)

	frmOpts := tv.CreateFrame(view, 1, 1, tv.BorderNone, types.Fixed)
	frmOpts.SetPack(tv.Horizontal)
	cbCheck := tv.CreateCheckBox(frmOpts, types.AutoSize, "Use callback to test data", types.Fixed)

	tv.CreateLabel(view, types.AutoSize, types.AutoSize, "Correct credentials", types.Fixed)

	frmCreds := tv.CreateFrame(view, 1, 1, tv.BorderNone, types.Fixed)
	frmCreds.SetPack(tv.Horizontal)
	frmCreds.SetGaps(1, 0)
	tv.CreateLabel(frmCreds, types.AutoSize, types.AutoSize, "Username", types.Fixed)
	edUser := tv.CreateEditField(frmCreds, 8, "", 1)
	tv.CreateLabel(frmCreds, types.AutoSize, types.AutoSize, "Password", types.Fixed)
	edPass := tv.CreateEditField(frmCreds, 8, "", 1)

	lbRes := tv.CreateLabel(view, types.AutoSize, types.AutoSize, "Result:", types.Fixed)

	frmBtns := tv.CreateFrame(view, 1, 1, tv.BorderNone, types.Fixed)
	frmBtns.SetPack(tv.Horizontal)
	btnDlg := tv.CreateButton(frmBtns, types.AutoSize, 4, "Login", types.Fixed)
	btnQuit := tv.CreateButton(frmBtns, types.AutoSize, 4, "Quit", types.Fixed)
	tv.CreateFrame(frmBtns, 1, 1, tv.BorderNone, 1)

	tv.ActivateControl(view, edUser)

	btnDlg.OnClick(func(ev tv.Event) {
		dlg := tv.CreateLoginDialog(
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
			if dlg.Action == tv.LoginCanceled {
				lbRes.SetTitle("Result:\nDialog canceled")
				return
			}

			if dlg.Action == tv.LoginInvalid {
				lbRes.SetTitle("Result:\nInvalid username or password")
				return
			}

			if dlg.Action == tv.LoginOk {
				if cbCheck.State() == 1 {
					lbRes.SetTitle("Result:\nLogged in successfully")
				} else {
					lbRes.SetTitle("Result:\nEntered [" + dlg.Username + ":" + dlg.Password + "]")
				}
				return
			}
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
