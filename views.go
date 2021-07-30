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
	dataPath         *JsonPathEvaluator
	selectExpression *TemplateEvaluator
	columns          []*TableColumn
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
			cell := t.GetCell(j+1, x)
			cell.Reference = row
		}
	}
	t.SetFixed(1, len(t.columns))

}

func (t *Table) HandleSelect(row, column int) {

	if t.selectExpression != nil {
		cell := t.GetCell(row, column)
		cmd := t.selectExpression.ExecuteWithDataAndCtx(cell.Reference, MainContext)
		log.Printf("Selected command [%s]\n", cmd)
	}
}

func NewTableColumn(headerExpression, dataExpression string) *TableColumn {
	return &TableColumn{NewTemplateEvaluator(headerExpression),
		NewTemplateEvaluator(dataExpression)}
}

func NewTable(id, name, shortcut, dataPath string, columns []*TableColumn, selectExpression string) *Table {
	tbl := tview.NewTable()
	tbl.SetBorder(true)
	tbl.SetTitle(name + "(" + shortcut + ")")
	jsonPath := NewJsonPathEvaluator(dataPath)
	baseView := NewBaseView(id, name, shortcut)
	newTable := &Table{Table: tbl, BaseView: baseView, dataPath: jsonPath, columns: columns}
	newTable.Table.SetSelectedFunc(newTable.HandleSelect)
	if len(selectExpression) > 0 {
		newTable.selectExpression = NewTemplateEvaluator(selectExpression)
	}
	return newTable
}

func NewTableFromConfig(id, name, shortcut, dataPath string, tableConfig TableConfig) *Table {
	columns := make([]*TableColumn, 0)
	for _, colConf := range tableConfig.Columns {
		col := &TableColumn{headerTemplate: NewTemplateEvaluator(colConf.HeaderExpression),
			dataTemplate: NewTemplateEvaluator(colConf.DataExpression)}
		columns = append(columns, col)
	}
	return NewTable(id, name, shortcut, dataPath, columns, tableConfig.SelectExpression)
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

type DataForm struct {
	*tview.Grid
	*BaseView
	fields []*DataFormField
}

type DataFormField struct {
	id            string
	labelTemplate *TemplateEvaluator
	valueTemplate *TemplateEvaluator
	field         *LabelValue
}

func (t *DataForm) DrawView(appCtx AppContext) {
	for _, field := range t.fields {
		label := field.labelTemplate.ExecuteWithCtx(appCtx)
		value := field.valueTemplate.ExecuteWithCtx(appCtx)
		field.field.SetLabel(label)
		field.field.SetValue(value)
	}
}

func NewDataForm(id, name, shortcut string, config DataFormConfig) *DataForm {
	grid := tview.NewGrid()
	grid.SetBorder(true)
	grid.SetTitle(name + "(" + shortcut + ")")
	baseView := NewBaseView(id, name, shortcut)

	log.Printf("Drawing DataForm\n")

	fields := make([]*DataFormField, 0)
	for _, labelConfig := range config.Fields {

		log.Printf("Drawing Label %s \n", labelConfig.Id)
		dir := LabelHorizontal
		if labelConfig.Orientation == "v" {
			dir = LabelVertical
		}

		labelTemplate := NewTemplateEvaluator(labelConfig.LabelExpression)
		valueTemplate := NewTemplateEvaluator(labelConfig.ValueExpression)
		labelValue := NewLabelValue(dir)
		grid.AddItem(labelValue, labelConfig.X, labelConfig.Y, 1, 1, -1, -1, false)
		fields = append(fields, &DataFormField{id: labelConfig.Id, labelTemplate: labelTemplate,
			valueTemplate: valueTemplate, field: labelValue})

	}

	return &DataForm{BaseView: baseView, Grid: grid, fields: fields}
}
