package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
)

func main() {

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

       


        textView := createInputField()
        //mainView := createMainView()
        pages := createPages()
        basicWindow := createBasicWindow(textView, pages)
        
	app := tview.NewApplication()
         
	frame := tview.NewFrame(basicWindow).AddText("Course Catalog", true, tview.AlignCenter, tcell.ColorGreen)

	if err := app.SetRoot(frame, true).SetFocus(basicWindow).Run(); err != nil {
		panic(err)
	}
}

func createInputField() *tview.TextView {
    return tview.NewTextView()
}

func createMainView() *tview.Flex {
    return tview.NewFlex()
}

func createBasicWindow(input *tview.TextView, main *tview.Flex) *tview.Flex {
	return tview.NewFlex().
		AddItem(main, 0, 1, false).
		AddItem(input, 1, 0, false)  
}

func createPages() *tview.Flex {
	horLabel := NewLabelValue(LabelHorizontal, "Label", "Value")
	horLabel.SetBorder(true)
	horLabel.SetTitle("Courses")

	vertLabel := NewLabelValue(LabelVertical, "Label", "Value")
	vertLabel.SetBorder(true)
	vertLabel.SetTitle("Course Info")

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(horLabel, 0, 2, false).
			AddItem(vertLabel, 0, 1, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Students"), 0, 2, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Student Info"), 0, 1, false), 0, 1, false)
       return flex
}

type AppContext  struct {

}

