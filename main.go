package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	records := []string{"Record 1", "Record 2", "Record 3"}

	timeView := tview.NewTextView().SetText("Some Text")
	timeView.SetBorder(true).SetTitle("Time")

	recordList := tview.NewList().ShowSecondaryText(false)
	recordList.SetBorder(true).SetTitle("Records").SetBorderPadding(0, 0, 1, 1)
	recordList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortut rune) {
		timeView.Clear()
		timeView.SetText("Time for " + records[index])
	})

	for _, record := range records {
		recordList.AddItem(record, "", 0, nil)
	}

	mainView := tview.NewFlex().
		AddItem(recordList, 0, 3, true).
		AddItem(timeView, 0, 1, false)

	footer := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(" a) add record s) start timer e) end timer")

	footer.SetBorderPadding(0, 0, 1, 1)

	app := tview.NewApplication()

	newRecordName := ""

	addRecordForm := tview.NewForm().
		AddInputField("New Record Name", "", 20, nil, func(text string) {
			newRecordName = text
		}).
		SetButtonsAlign(tview.AlignRight)

	addRecordForm.SetBorder(true).SetBorderColor(tcell.ColorBlue)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(addRecordForm, 7, 1, false).
			AddItem(nil, 0, 1, false), 40, 1, false).
		AddItem(nil, 0, 1, false)

	grid := tview.NewGrid().
		SetRows(0, 1).
		AddItem(mainView, 0, 0, 1, 1, 0, 0, true).
		AddItem(footer, 1, 0, 1, 1, 0, 0, false)

	pages := tview.NewPages().
		AddPage("main", grid, true, true)

	addRecordForm.
		AddButton("Save", func() {
			footer.SetText(newRecordName)
			pages.RemovePage("modal").ShowPage("main")
		}).
		AddButton("Quit", func() {
			pages.RemovePage("modal").ShowPage("main")
		})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' && !pages.HasPage("modal") {
			pages.AddAndSwitchToPage("modal", modal, true).ShowPage("main")
			app.SetFocus(addRecordForm)
			return nil
		}

		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
