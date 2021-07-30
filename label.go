package main

import (
	"fmt"
	//"log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	LabelHorizontal = iota
	LabelVertical
)

type LabelValue struct {
	*tview.Box
	direction int
	label     string
	value     string
}

func NewLabelValue(direction int) *LabelValue {
	return &LabelValue{
		Box:       tview.NewBox(),
		direction: direction,
	}
}

func (r *LabelValue) SetLabel(label string) {
	r.label = label
}

func (r *LabelValue) GetLabel() string {
	return r.label
}

func (r *LabelValue) SetValue(value string) {
	r.value = value
}

func (r *LabelValue) GetValue() string {
	return r.value
}

func (r *LabelValue) Draw(screen tcell.Screen) {
	r.Box.DrawForSubclass(screen, r)
	x, y, width, height := r.GetInnerRect()
	if r.direction == LabelHorizontal {
		line := fmt.Sprintf(`%s  %s`, r.label, r.value)
		tview.Print(screen, line, x, y, width, tview.AlignLeft, tcell.ColorGreen)
	} else {
		if height > 1 {
			label := fmt.Sprintf(`%s`, r.label)
			value := fmt.Sprintf(`%s`, r.value)
			tview.Print(screen, label, x, y, width, tview.AlignLeft, tcell.ColorYellow)
			tview.Print(screen, value, x, y+1, width, tview.AlignLeft, tcell.ColorGreen)
		}
	}
}
