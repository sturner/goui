package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const Row = "row"
const Col = "col"

var configMaster string
var MainContext AppContext

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	flag.StringVar(&configMaster, "c", "", "Configuration file")
}

func main() {

	flag.Parse()
	if len(configMaster) <= 0 {
		fmt.Printf("Need config file [-c] option\n")
		os.Exit(1)
	}

	appConfig, err := loadAppConfig(configMaster)
	if err != nil {
		fmt.Printf("Error loading config file i%s\n", configMaster)
		os.Exit(2)
	}
	rootPages := tview.NewPages()
	app := tview.NewApplication()
	appContext := BuildAppContext(appConfig, app, rootPages)
	MainContext = appContext
	controller := NewController(appContext, appConfig)
	BuildLayoutFromConfig(appConfig, appContext, rootPages)
	rootPages.SwitchToPage(appConfig.Pages[0].Id)
	inputHandler := NewInputHandler(app, appContext, controller.processCommand)
	basicWindow := createBasicWindow(inputHandler.inputView, rootPages)

	frame := tview.NewFrame(basicWindow).AddText("Course Catalog", true, tview.AlignCenter, tcell.ColorGreen)

	if err := app.SetRoot(frame, true).SetFocus(basicWindow).Run(); err != nil {
		panic(err)
	}
}

func loadAppConfig(config string) (*ApplicationConfig, error) {
	return loadConfig(config)
}

func createInputField() *tview.TextView {
	return tview.NewTextView()
}

func createMainView() *tview.Flex {
	return tview.NewFlex()
}

func createBasicWindow(input *tview.TextView, main tview.Primitive) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(main, 0, 1, false).
		AddItem(input, 1, 0, false)
}

func createPagesFromConfig(config *ApplicationConfig) []*Page {
	pages := make([]*Page, 0)
	for _, page := range config.Pages {
		log.Printf("Building page [%s]\n", page.Id)
		views := make([]View, 0)
		for _, view := range page.Views {
			log.Printf("Building view [%s]\n", view.Id)
			v := createViewFromConfig(view)
			if v != nil {
				views = append(views, v)
				log.Printf("Finished building view [%s]\n", view.Id)
			} else {
				log.Printf("Invalid view [%s]\n", view.Id)
			}
		}
		pages = append(pages, NewPage(page.Id, page.Name, page.Shortcut, views))
		log.Printf("Finished building page [%s]\n", page.Id)
	}
	log.Printf("Finished building all pages [%+v]\n", pages)
	return pages
}

func createViewFromConfig(viewConfig ViewConfig) View {

	if len(viewConfig.Table.Columns) > 0 {
		return NewTableFromConfig(viewConfig.Id, viewConfig.Name, viewConfig.Shortcut, viewConfig.DataPath, viewConfig.Table)
	} else if len(viewConfig.Static) > 0 {
		return NewPlaceholder(viewConfig.Id, viewConfig.Name, viewConfig.Shortcut, viewConfig.DataPath)
	} else if len(viewConfig.Form.Fields) > 0 {
		return NewDataForm(viewConfig.Id, viewConfig.Name, viewConfig.Shortcut, viewConfig.Form)
	}
	return nil
}

type AppContext interface {
	RegisterData(key string, data interface{})
	GetData(key string) interface{}
	GetDataMap() map[string]interface{}
	GetPageById(id string) *Page
	GetPageByShortcut(shortcut string) *Page
	GetView(key string) (View, *Page)
	SwitchPage(pageShortcut string)
	FocusOnViewShortcut(viewShortcut string)
	FocusOnViewId(viewId string)
	RegisterArgs(args []string)
	GetArguments() []string
	Quit()
}

type BaseAppContext struct {
	data   map[string]interface{}
	pages  []*Page
	vPages *tview.Pages
	app    *tview.Application
}

func (b *BaseAppContext) RegisterData(key string, data interface{}) {
	b.data[key] = data
}

func (b *BaseAppContext) GetData(key string) interface{} {
	return b.data[key]
}

func (b *BaseAppContext) GetDataMap() map[string]interface{} {
	return b.data
}

func (b *BaseAppContext) GetPageById(id string) *Page {
	log.Printf("Looking for page id [%s]\n", id)
	for _, p := range b.pages {
		log.Printf("Checking page id [%s]\n", p.Id)
		if id == p.Id {
			log.Printf("Found page [%s]\n", id)
			return p
		}
	}
	return nil
}

func (b *BaseAppContext) GetPageByShortcut(shortcut string) *Page {
	log.Printf("Looking for page shortcut [%s]\n", shortcut)
	for _, p := range b.pages {
		log.Printf("Checking page shortcut [%s]\n", p.Shortcut)
		if shortcut == p.Shortcut {
			log.Printf("Found page [%+v]\n", p)
			return p
		}
	}
	return nil
}
func (b *BaseAppContext) GetView(key string) (View, *Page) {
	for _, p := range b.pages {
		for _, v := range p.Views {
			if strings.Compare(v.GetId(), key) == 0 {
				return v, p
			}
		}
	}
	return nil, nil
}
func (b *BaseAppContext) Quit() {
	b.app.Stop()
}
func (b *BaseAppContext) SwitchPage(pageShortcut string) {
	p := b.GetPageByShortcut(pageShortcut)
	if p != nil {
		b.vPages.SwitchToPage(p.Id)
	}
}

func (b *BaseAppContext) FocusOnViewId(viewId string) {
	v, p := b.GetView(viewId)
	log.Printf("Focus on  page :[%+v] - view [%+v]\n", p, v)
	b.SwitchPage(p.Shortcut)
	b.FocusOnViewShortcut(v.GetShortcut())

}

func (b *BaseAppContext) FocusOnViewShortcut(viewShortcut string) {

	pageId, _ := b.vPages.GetFrontPage()
	log.Printf("Focusing on page [%s]\n", pageId)
	p := b.GetPageById(pageId)
	log.Printf("Got page [%+v]\n", p)
	if p != nil {
		for _, view := range p.Views {
			log.Printf("Got page view  [%+v]\n", view)
			if view.GetShortcut() == viewShortcut {
				log.Printf("Found page view [%+v], setting focus\n", view)
				b.app.SetFocus(view)
			}
		}
	}
}
func (b *BaseAppContext) RegisterArgs(args []string) {
	b.data["args"] = args
}

func (b *BaseAppContext) GetArguments() []string {
	if data, ok := b.data["args"]; ok {
		if args, ok := data.([]string); ok {
			return args
		}
	}
	return make([]string, 0)
}

func BuildAppContext(config *ApplicationConfig, app *tview.Application, rootPages *tview.Pages) AppContext {
	pages := createPagesFromConfig(config)
	baseContext := &BaseAppContext{data: config.Data, pages: pages, app: app, vPages: rootPages}
	return baseContext
}

func BuildLayoutFromConfig(config *ApplicationConfig, ctx AppContext, rootPages *tview.Pages) {
	for _, pageConfig := range config.Pages {
		log.Printf("Building page layout [%s]\n", pageConfig.Id)
		page := ctx.GetPageById(pageConfig.Id)
		viewLayout := buildLayout(pageConfig.Layout.Dir, pageConfig.Layout.Containers, pageConfig.Layout.Views, ctx)
		page.Flex = viewLayout
		log.Printf("Finished building page layout [%+v] for page [%+v]\n", viewLayout, page)
		rootPages.AddPage(pageConfig.Id, page, true, false)
	}
}

func getDirection(dir string) int {
	if strings.Compare(dir, Row) != 0 && strings.Compare(dir, Col) != 0 {
		log.Fatalf("Invalid layout direction: [%s], valid values are [\"%s\", \"%s\"]", dir, Row, Col)
	}
	layoutDirection := tview.FlexColumn
	if strings.Compare(dir, Row) == 0 {
		layoutDirection = tview.FlexRow
	}
	return layoutDirection
}

func buildLayout(dir string, containerConfigs []ContainerLayoutConfig, viewConfigs []ViewLayoutConfig, ctx AppContext) *tview.Flex {
	flexLayoutDir := getDirection(dir)
	flexLayout := tview.NewFlex().SetDirection(flexLayoutDir)
	log.Printf("Building layout\n")

	for _, vLayout := range viewConfigs {
		log.Printf("Building view layout [%s]\n", vLayout.ViewId)
		v, _ := ctx.GetView(vLayout.ViewId)
		if v == nil {
			log.Fatalf("View id %s not found\n", vLayout.ViewId)
		}
		flexLayout.AddItem(v, vLayout.FixedSize, vLayout.Proportion, false)
		log.Printf("Finished view layout[%s]\n", vLayout.ViewId)
	}

	for _, cLayout := range containerConfigs {
		log.Printf("Building container layout\n")
		container := buildLayout(cLayout.Dir, cLayout.Containers, cLayout.Views, ctx)
		flexLayout.AddItem(container, cLayout.FixedSize, cLayout.Proportion, false)
		log.Printf("Finished container layout\n")
	}
	log.Printf("Finished layout\n")
	return flexLayout

}

type Controller struct {
	cmdProcessor *CommandProcessor
	appContext   AppContext
}

func BuildCommandProcessor(appConfig *ApplicationConfig) *CommandProcessor {
	return NewCommandProcessor(appConfig)
}

func NewController(ctx AppContext, appConfig *ApplicationConfig) *Controller {
	return &Controller{cmdProcessor: NewCommandProcessor(appConfig), appContext: ctx}
}

func (c *Controller) processCommand(cmd string) {
	if len(cmd) > 0 {
		cmdResult, err := c.cmdProcessor.Process(cmd, c.appContext)
		if err != nil {
			log.Printf("Error executing command [%s] [%+v]\n", cmd, err)
			return
		}
		if cmdResult != nil {
			log.Printf("Command results: [%+v]\n", cmdResult.Data)
			c.appContext.RegisterData(cmdResult.Key, cmdResult.Data)
			if len(cmdResult.ViewId) > 0 {
				v, _ := c.appContext.GetView(cmdResult.ViewId)
				if v != nil {
					v.DrawView(c.appContext)
					c.appContext.FocusOnViewId(cmdResult.ViewId)
				}
			}
		}
	}
}

type Page struct {
	*tview.Flex
	Id       string
	Name     string
	Shortcut string
	Views    []View
}

func NewPage(id, name, shortcut string, views []View) *Page {
	return &Page{Id: id, Name: name, Shortcut: shortcut, Views: views}
}

type InputHandler struct {
	inputView   *tview.TextView
	buffer      *strings.Builder
	application *tview.Application
	appContext  AppContext
	processFunc func(string)
}

func NewInputHandler(app *tview.Application, appContext AppContext, cmdProcessor func(string)) *InputHandler {
	handler := &InputHandler{inputView: tview.NewTextView(),
		buffer:      &strings.Builder{},
		application: app,
		appContext:  appContext,
		processFunc: cmdProcessor}
	app.SetInputCapture(handler.InputCapture)
	return handler
}

func (i *InputHandler) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	log.Printf("%+v, %q, %+v\n", event.Key(), event.Rune(), event.Modifiers())
	//if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyETX {
	//		i.application.Stop()
	//	return event
	//} else
	if event.Rune() == ':' {
		i.buffer.Reset()
		i.buffer.WriteString(":")
		i.application.SetFocus(i.inputView)
		return nil
	}
	focus := i.application.GetFocus()
	if focus != i.inputView {
		return event
	}
	if event.Key() == tcell.KeyEsc {
		i.buffer.Reset()
		i.inputView.SetText(i.buffer.String())
		i.application.SetFocus(i.inputView)
		return nil
	} else if event.Rune() == ':' {
		i.buffer.Reset()
		i.buffer.WriteString(":")
		i.application.SetFocus(i.inputView)
		return nil
	} else if event.Key() == tcell.KeyEnter {
		if len(i.buffer.String()) == 0 {
			return event
		}
		i.processFunc(i.buffer.String()[1:])
		i.buffer.Reset()
		i.inputView.SetText("")
		return nil
	} else if event.Key() == tcell.KeyDelete || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyDEL {
		str := i.buffer.String()
		i.buffer.Reset()
		if len(str) > 0 {
			i.buffer.WriteString(str[0 : len(str)-1])
			i.inputView.SetText(i.buffer.String())
		}
		return nil
	} else {
		i.buffer.WriteRune(event.Rune())
		i.inputView.SetText(i.buffer.String())
		return nil
	}
}
