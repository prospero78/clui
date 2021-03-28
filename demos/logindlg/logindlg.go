package main

import (
	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func createView() {
	view := tv.AddWindow(0, 0, 30, 7, "Login dialog")
	view.SetPack(cons.Vertical)
	view.SetGaps(0, 1)
	view.SetPaddings(2, 2)

	frmOpts := tv.CreateFrame(view, 1, 1, cons.BorderNone, cons.Fixed)
	frmOpts.SetPack(cons.Horizontal)
	cbCheck := tv.CreateCheckBox(frmOpts, cons.AutoSize, "Use callback to test data", cons.Fixed)

	tv.CreateLabel(view, cons.AutoSize, cons.AutoSize, "Correct credentials", cons.Fixed)

	frmCreds := tv.CreateFrame(view, 1, 1, cons.BorderNone, cons.Fixed)
	frmCreds.SetPack(cons.Horizontal)
	frmCreds.SetGaps(1, 0)
	tv.CreateLabel(frmCreds, cons.AutoSize, cons.AutoSize, "Username", cons.Fixed)
	edUser := tv.CreateEditField(frmCreds, 8, "", 1)
	tv.CreateLabel(frmCreds, cons.AutoSize, cons.AutoSize, "Password", cons.Fixed)
	edPass := tv.CreateEditField(frmCreds, 8, "", 1)

	lbRes := tv.CreateLabel(view, cons.AutoSize, cons.AutoSize, "Result:", cons.Fixed)

	frmBtns := tv.CreateFrame(view, 1, 1, cons.BorderNone, cons.Fixed)
	frmBtns.SetPack(cons.Horizontal)
	btnDlg := tv.CreateButton(frmBtns, cons.AutoSize, 4, "Login", cons.Fixed)
	btnQuit := tv.CreateButton(frmBtns, cons.AutoSize, 4, "Quit", cons.Fixed)
	tv.CreateFrame(frmBtns, 1, 1, cons.BorderNone, 1)

	tv.ActivateControl(view, edUser)

	btnDlg.OnClick(func(ev event.TEvent) {
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

	btnQuit.OnClick(func(ev event.TEvent) {
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
