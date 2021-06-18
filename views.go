package main

import (
	"log"

	"github.com/rivo/tview"
)

type View interface {
	tview.Primitive
	DrawView(appCtx AppContext)
	GetId() string
	GetName() string
	GetShortcut() string
	GetDataPath() string
}

type BaseView struct {
	Name     string
	Id       string
	Shortcut string
	DataPath string
}

func (b *BaseView) GetId() string {
	return b.Id
}
func (b *BaseView) GetName() string {
	return b.Name
}
func (b *BaseView) GetShortcut() string {
	return b.Shortcut
}
func (b *BaseView) GetDataPath() string {
	return b.DataPath
}

func NewBaseView(id, name, shortcut string) *BaseView {
	return &BaseView{Id: id, Name: name, Shortcut: shortcut}
}

type Table struct {
	*tview.Table
	*BaseView
	dataPath *JsonPathEvaluator
	columns  []*TableColumn
}

type TableColumn struct {
	headerTemplate *TemplateEvaluator
	dataTemplate   *TemplateEvaluator
}

func (t *Table) DrawView(appCtx AppContext) {
	log.Printf("Drawing table\n")
	t.Clear()
	t.SetSelectable(true, false)
	dataRows, err := t.dataPath.ExecuteWithCtx(appCtx, nil)
	if err != nil {
		log.Printf("Error retrieving data for table [%+v]\n", err)
		return
	}
	for i, c := range t.columns {
		headerValue := c.headerTemplate.ExecuteWithCtx(appCtx)
		t.SetCellSimple(0, i, headerValue)
	}
	for j, row := range dataRows.([]interface{}) {
		for x, col := range t.columns {
			dataValue := col.dataTemplate.Execute(row)
			t.SetCellSimple(j+1, x, dataValue)
		}
	}
	t.SetFixed(1, len(t.columns))
	//func (r *JsonPathEvaluator) ExecuteWithCtx(ctx AppContext, data interface{}) (interface{}, error) {
	//t.SetCellSimple(row, column int, text string)

}

func NewTableColumn(headerExpression, dataExpression string) *TableColumn {
	return &TableColumn{NewTemplateEvaluator(headerExpression),
		NewTemplateEvaluator(dataExpression)}
}

func NewTable(id, name, shortcut, dataPath string, columns []*TableColumn) *Table {
	tbl := tview.NewTable()
	tbl.SetBorder(true)
	tbl.SetTitle(name + "(" + shortcut + ")")
	jsonPath := NewJsonPathEvaluator(dataPath)
	baseView := NewBaseView(id, name, shortcut)
	return &Table{Table: tbl, BaseView: baseView, dataPath: jsonPath, columns: columns}
}

func NewTableFromConfig(id, name, shortcut, dataPath string, columnConfig []TableItemConfig) *Table {
	columns := make([]*TableColumn, 0)
	for _, colConf := range columnConfig {
		col := &TableColumn{headerTemplate: NewTemplateEvaluator(colConf.HeaderExpression),
			dataTemplate: NewTemplateEvaluator(colConf.DataExpression)}
		columns = append(columns, col)
	}
	return NewTable(id, name, shortcut, dataPath, columns)
}

type Placeholder struct {
	*tview.TextView
	*BaseView
}

func (t *Placeholder) DrawView(appCtx AppContext) {

}

func NewPlaceholder(id, name, shortcut, dataPath string) *Placeholder {
	textView := tview.NewTextView()
	textView.SetText("Lorem ipsum")
	textView.SetBorder(true)
	textView.SetTitle(name + "(" + shortcut + ")")
	baseView := NewBaseView(id, name, shortcut)
	return &Placeholder{TextView: textView, BaseView: baseView}
}
