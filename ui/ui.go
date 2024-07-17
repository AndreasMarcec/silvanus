package ui

import (
	awsilvanus "github.com/AndreasMarcec/silvanus/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Tui struct {
	App      *tview.Application
	LogsView *tview.TextView
	Flex     *tview.Flex
	Table    *tview.Table
	Client   *awsilvanus.FunctionWrapper
}

func (t *Tui) Run() {
	if err := t.App.SetRoot(t.Flex, true).SetFocus(t.Table).EnableMouse(true).Run(); err != nil {
		panic(err.Error())
	}
}

func (t *Tui) Debug(debugText string) {
	t.LogsView.SetText(debugText)
}

func (t *Tui) UpdateTable() {
	lorem := t.Client.ListFunctions(100)
	rows := len(lorem)

	// Fill the cells of the table
	for r := 0; r < rows; r++ {
		color := tcell.ColorWhite
		t.Table.SetCell(r, 0,
			tview.NewTableCell(aws.ToString(lorem[r].FunctionName)).
				SetTextColor(color).
				SetAlign(tview.AlignCenter))
	}

	t.Table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			t.App.Stop()
		}
		if key == tcell.KeyEnter {
			t.Table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		cell := t.Table.GetCell(row, column)
		text := t.Client.GetLogs(cell.Text)
		t.LogsView.SetText(text)
		cell.SetTextColor(tcell.ColorRed)
		t.Table.SetSelectable(false, false)
	})
}

func InitTui(t *Tui) {
	t.App = tview.NewApplication()
	t.Flex = tview.NewFlex()
	t.LogsView = tview.NewTextView()
	t.Table = tview.NewTable().SetBorders(true)

	// TODO Cleanup
	t.Client = &awsilvanus.FunctionWrapper{}
	t.Client.LambdaClient = t.Client.InitLambdaClient()

	t.LogsView.SetBorder(true)
	t.Flex.AddItem(t.LogsView, 0, 1, false)
	t.Flex.AddItem(t.Table, 0, 1, true)
}

func Create() *Tui {
	return new(Tui)
}
