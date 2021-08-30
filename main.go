package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type timeInterval struct {
	startTime time.Time
	endTime   time.Time
}

type record struct {
	name          string
	timeIntervals []timeInterval
}

func main() {
	timerRunning := false
	activeRecord := 0
	records := []record{
		{name: "Record 1"},
		{name: "Record 2"},
		{name: "Record 3"},
	}

	timeView := tview.NewTextView()
	timeView.SetBorder(true).SetTitle("Time")

	recordList := tview.NewList().ShowSecondaryText(false)
	recordList.SetBorder(true).SetTitle("Records").SetBorderPadding(0, 0, 1, 1)
	// recordList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortut rune) {
	// 	timeView.Clear()
	// 	timeView.SetText("Time for " + records[index])
	// })

	renderTimeView := func() {
		timeViewOutput := ""

		for currentRecord, record := range records {
			timeViewOutput += record.name + "\n"
			runTime := 0 * time.Second

			for timeIntervalIndex, timeInterval := range record.timeIntervals {
				if timerRunning && currentRecord == activeRecord && timeIntervalIndex == len(record.timeIntervals)-1 {
					runTime += time.Now().UTC().Sub(timeInterval.startTime)
					timeViewOutput += "D: " + fmt.Sprintf("%v", time.Now().UTC().Sub(timeInterval.startTime)) + "\n"
				} else {
					runTime += timeInterval.endTime.Sub(timeInterval.startTime)
					timeViewOutput += "D: " + fmt.Sprintf("%v", timeInterval.endTime.Sub(timeInterval.startTime)) + "\n"
				}
				// timeViewOutput += "S: " + timeInterval.startTime.Format(time.RFC3339) + "\nE: " + timeInterval.endTime.Format(time.RFC3339) + "\n"
			}

			timeViewOutput += " " + fmt.Sprintf("%.0f", runTime.Seconds()) + "\n"
		}

		timeView.SetText(timeViewOutput)
	}

	recordList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if timerRunning {
			records[activeRecord].timeIntervals[len(records[activeRecord].timeIntervals)-1].endTime = time.Now().UTC()
			records[index].timeIntervals = append(records[index].timeIntervals, timeInterval{startTime: time.Now().UTC()})

			renderTimeView()
		}

		activeRecord = index
	})

	for _, record := range records {
		recordList.AddItem(record.name, "", 0, nil)
	}

	mainView := tview.NewFlex().
		AddItem(recordList, 0, 3, true).
		AddItem(timeView, 0, 1, false)

	footer := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(" a) add record s) start/stop timer [stopped]")

	footer.SetBorderPadding(0, 0, 1, 1)

	app := tview.NewApplication()

	addRecordForm := tview.NewForm().
		AddInputField("New Record Name", "", 20, nil, nil).
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
			inputField := addRecordForm.GetFormItem(0).(*tview.InputField)
			newRecordName := inputField.GetText()
			records = append(records, record{name: newRecordName})
			recordList.AddItem(newRecordName, "", 0, nil)
			renderTimeView()
			pages.RemovePage("modal").ShowPage("main")
			inputField.SetText("")
		}).
		AddButton("Quit", func() {
			addRecordForm.GetFormItem(0).(*tview.InputField).SetText("")
			pages.RemovePage("modal").ShowPage("main")
		})

	go func() {
		for {
			renderTimeView()
			time.Sleep(1 * time.Second)
			app.Draw()
		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' && !pages.HasPage("modal") {
			pages.AddAndSwitchToPage("modal", modal, true).ShowPage("main")
			app.SetFocus(addRecordForm.SetFocus(0))
			return nil
		}

		if event.Rune() == 's' && mainView.HasFocus() {
			if !timerRunning {
				records[activeRecord].timeIntervals = append(records[activeRecord].timeIntervals, timeInterval{startTime: time.Now().UTC()})
			}

			if timerRunning {
				records[activeRecord].timeIntervals[len(records[activeRecord].timeIntervals)-1].endTime = time.Now().UTC()
			}

			timerRunning = !timerRunning
			if timerRunning {
				footer.SetText(" a) add record s) start/stop timer [running]")
			} else {
				footer.SetText(" a) add record s) start/stop timer [stopped]")
			}
			renderTimeView()
		}

		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
