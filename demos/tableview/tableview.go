package main

import (
	"fmt"

	"github.com/prospero78/goTV/tv"
	"github.com/prospero78/goTV/tv/cons"
	"github.com/prospero78/goTV/tv/widgets/event"
)

func createView() *tv.TableView {

	view := tv.AddWindow(0, 0, 10, 7, "TableView Demo")
	bch := tv.CreateTableView(view, 25, 12, 1)
	tv.ActivateControl(view, bch)

	return bch
}

const rowCount = 15

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	tv.InitLibrary()
	defer tv.DeinitLibrary()

	b := createView()
	b.SetShowLines(true)
	b.SetShowRowNumber(true)
	b.SetRowCount(rowCount)
	cols := []tv.Column{
		{Title: "Text", Width: 5, Alignment: cons.AlignLeft},
		{Title: "Number", Width: 10, Alignment: cons.AlignRight},
		{Title: "Misc", Width: 12, Alignment: cons.AlignCenter},
		{Title: "Long", Width: 50, Alignment: cons.AlignLeft},
		{Title: "Last", Width: 8, Alignment: cons.AlignLeft},
	}
	b.SetColumns(cols)
	colCount := len(cols)

	values := make([]string, rowCount*colCount)
	for r := 0; r < rowCount; r++ {
		for c := 0; c < colCount; c++ {
			values[r*colCount+c] = fmt.Sprintf("%v:%v", r, c)
		}
	}

	b.OnDrawCell(func(info *tv.ColumnDrawInfo) {
		info.Text = values[info.Row*colCount+info.Col]
	})

	b.OnAction(func(ev tv.TableEvent) {
		btns := []string{"Close", "Dismiss"}
		var action string
		switch ev.Action {
		case cons.TableActionSort:
			action = "Sort table"
		case cons.TableActionEdit:
			c := ev.Col
			r := ev.Row
			oldVal := values[r*colCount+c]
			dlg := tv.CreateEditDialog(
				fmt.Sprintf("Editing value: %s", oldVal), "New value", oldVal,
			)
			dlg.OnClose(func() {
				if dlg.Result() == cons.DialogButton1 {
					newText := dlg.EditResult()
					values[r*colCount+c] = newText
					tv.PutEvent(event.TEvent{Type: cons.EventRedraw})
				}
			})
			return
		case cons.TableActionNew:
			action = "Add new row"
		case cons.TableActionDelete:
			action = "Delete row"
		default:
			action = "Unknown action"
		}

		dlg := tv.CreateConfirmationDialog(
			"<c:blue>"+action,
			"Click any button or press <c:yellow>SPACE<c:> to close the dialog",
			btns, cons.DialogButton1)
		dlg.OnClose(func() {})
	})

	// start event processing loop - the main core of the library
	tv.MainLoop()
}

func main() {
	mainLoop()
}
