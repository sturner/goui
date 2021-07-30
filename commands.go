package main

import (
	"regexp"
	"strings"

	//"text/template"
	"os/exec"
	//"strings"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	//"io"
)

var SpaceSplitter = regexp.MustCompile(`\s+`)

const (
	OutputJson = iota
	OutputCsv
	OutputNone
)

type Command interface {
	GetName() string
	CanProcess(cmdText string) bool
	GetArguments(cmdText string) []string
	Execute(cmdText string, appCtx AppContext) (*CommandResult, error)
	GetResultType() int
	GetResultKey() string
	GetViewId() string
	GetHelp() HelpConfig
}

type CommandResult struct {
	Type   int
	Data   interface{}
	Key    string
	ViewId string
}

func NewJsonResult(jsonString, key, viewId string) *CommandResult {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	return NewParsedJsonResult(jsonMap, key, viewId)
}

func NewParsedJsonResult(data interface{}, key, viewId string) *CommandResult {
	return &CommandResult{OutputJson, data, key, viewId}
}

type BaseCommand struct {
	Name       string
	regex      *regexp.Regexp
	resultType int
	resultKey  string
	filter     Filter
	viewId     string
	help       HelpConfig
}

func newBaseCommand(config CommandConfig) *BaseCommand {
	regex := regexp.MustCompile(config.Regex)
	var filter Filter
	if len(config.FilterExpression) > 0 {
		filter = NewJqFilter(config.FilterExpression)
	} else {
		filter = &EmptyFilter{}
	}
	return &BaseCommand{Name: config.Name, regex: regex, resultType: getResultType(config.ResultType),
		resultKey: config.ResultKey, filter: filter, viewId: config.ViewId, help: config.Help}

}

/*
func newBaseCommand(name string, resultType int, resultKey, cmdExpression, filterExpression, viewId string) *BaseCommand {
	regex := regexp.MustCompile(cmdExpression)
	var filter Filter
	if len(filterExpression) > 0 {
		filter = NewJqFilter(filterExpression)
	} else {
		filter = &EmptyFilter{}
	}
	return &BaseCommand{Name: name, regex: regex, resultType: resultType, resultKey: resultKey, filter: filter, viewId: viewId}
}
*/

func (r *BaseCommand) GetName() string {
	return r.Name
}

func (r *BaseCommand) GetResultType() int {
	return r.resultType
}

func (r *BaseCommand) GetResultKey() string {
	return r.resultKey
}

func (r *BaseCommand) CanProcess(cmdText string) bool {
	log.Printf("checking  command [%s] [%+v]\n", cmdText, r.regex)
	return r.regex.FindAllStringSubmatch(cmdText, -1) != nil
}
func (r *BaseCommand) GetViewId() string {
	return r.viewId
}

func (b *BaseCommand) ParseAndFilterString(data string, ctx AppContext) (*CommandResult, error) {
	jsonResult := NewJsonResult(data, b.GetResultKey(), b.GetViewId())
	return b.ParseAndFilter(jsonResult.Data, ctx)
}

func (b *BaseCommand) GetArguments(cmdText string) []string {
	result := make([]string, 0)
	log.Printf("Getting command arguments [%s]\n", cmdText)
	matches := b.regex.FindStringSubmatch(cmdText)
	argsLen := len(matches)
	if argsLen > 0 {
		return matches[1:]
	}
	log.Printf("Got arguments [%+v]\n", result)
	return result
}

func (b *BaseCommand) ParseAndFilter(data interface{}, ctx AppContext) (*CommandResult, error) {
	if b.resultType == OutputJson {
		jsonResult := NewParsedJsonResult(data, b.GetResultKey(), b.GetViewId())
		log.Printf("Filtering data [%+v]\n", jsonResult.Data)
		filtered := b.filter.Filter(jsonResult.Data, ctx)
		log.Printf("Filtered data [%+v]\n", filtered)
		return NewParsedJsonResult(filtered, b.GetResultKey(), b.GetViewId()), nil
	}
	return nil, nil
}

func (b *BaseCommand) GetHelp() HelpConfig {
	return b.help
}

type ShellCommand struct {
	*BaseCommand
	template *TemplateEvaluator
}

func (r *ShellCommand) Execute(cmdText string, ctx AppContext) (*CommandResult, error) {

	cmdString := r.template.ExecuteWithCtx(ctx)
	cmdArray := SpaceSplitter.Split(cmdString, -1)
	cmd := exec.Command(cmdArray[0], cmdArray[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return r.ParseAndFilterString(out.String(), ctx)

}

func NewShellCommand(config CommandConfig, shellExpression string) Command {
	shellExp := NewTemplateEvaluator(shellExpression)
	return &ShellCommand{newBaseCommand(config), shellExp}
}

type PassthruCommand struct {
	*BaseCommand
	SourceDataId string
}

func (p *PassthruCommand) Execute(cmdText string, ctx AppContext) (*CommandResult, error) {
	log.Printf("Passthru command [%s]\n", cmdText)
	log.Printf("Passthru command data [%s]\n", p.SourceDataId)
	data := ctx.GetData(p.SourceDataId)
	return p.ParseAndFilter(data, ctx)
}

func NewPassthruCommand(config CommandConfig, sourceDataId string) Command {
	return &PassthruCommand{newBaseCommand(config), sourceDataId}
}

type CommandProcessor struct {
	Processors map[string]Command
}

func (c *CommandProcessor) Process(cmd string, appCtx AppContext) (*CommandResult, error) {
	log.Printf("Processing command [%s]\n", cmd)
	for _, v := range c.Processors {
		log.Printf("Checking command [%+v]\n", v)
		if v.CanProcess(cmd) {
			args := v.GetArguments(cmd)
			appCtx.RegisterArgs(args)
			return v.Execute(cmd, appCtx)
		}
	}
	return nil, nil
}

func CreateCommand(config CommandConfig) Command {
	if len(config.ShellExpression) > 0 {
		return NewShellCommand(config, config.ShellExpression)
	}
	return NewPassthruCommand(config, config.PassthruSourceId)
}

func getResultType(resultType string) int {
	if strings.Compare(resultType, "json") == 0 {
		return OutputJson
	}
	return -1
}

func (c *CommandProcessor) Register(cmd Command) {
	log.Printf("Registering command [%s]\n", cmd.GetName())
	c.Processors[cmd.GetName()] = cmd
}

func NewCommandProcessor(appConfig *ApplicationConfig) *CommandProcessor {

	processor := &CommandProcessor{Processors: make(map[string]Command)}

	// register the system commands first
	sysCommands := make([]Command, 0)
	sysCommands = append(sysCommands, NewQuitCommand(), NewPageCommand(), NewFocusCommand())

	for _, sysCmd := range sysCommands {
		processor.Register(sysCmd)
	}

	for _, cmd := range appConfig.Commands {
		newCmd := CreateCommand(cmd)
		processor.Register(newCmd)
	}
	return processor
}

// Built-in commands that can't be overridden

type QuitCommand struct {
	*BaseCommand
}

func (q *QuitCommand) Execute(cmdText string, ctx AppContext) (*CommandResult, error) {
	log.Printf("Quitting application")
	ctx.Quit()
	return nil, nil
}

func NewQuitCommand() Command {
	config := CommandConfig{Name: "quit", Regex: "q", ResultType: "", FilterExpression: "", ViewId: ""}
	return &QuitCommand{newBaseCommand(config)}
}

type PageCommand struct {
	*BaseCommand
}

func (q *PageCommand) Execute(cmdText string, ctx AppContext) (*CommandResult, error) {
	args := ctx.GetArguments()
	pageId := args[0]
	log.Printf("Processing page command [%s] - [%s]\n", cmdText, args[0])
	ctx.SwitchPage(pageId)
	return nil, nil
}

func NewPageCommand() Command {
	config := CommandConfig{Name: "page", Regex: "p \\w+", ResultType: "", FilterExpression: "", ViewId: ""}
	return &PageCommand{newBaseCommand(config)}
}

type FocusCommand struct {
	*BaseCommand
}

func (q *FocusCommand) Execute(cmdText string, ctx AppContext) (*CommandResult, error) {
	results := strings.Split(cmdText, " ")
	viewShortcut := results[1]
	log.Printf("Processing focus command [%s] - [%s]\n", cmdText, viewShortcut)
	ctx.FocusOnViewShortcut(viewShortcut)
	return nil, nil
}

func NewFocusCommand() Command {
	config := CommandConfig{Name: "focus", Regex: "f \\w+", ResultType: "", FilterExpression: "", ViewId: ""}
	return &FocusCommand{newBaseCommand(config)}

}
