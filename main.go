package main

import (
	"github.com/AndreasMarcec/silvanus/ui"
)

func main() {
	t := ui.Create()
	ui.InitTui(t)
	t.UpdateTable()
	t.Run()
}
