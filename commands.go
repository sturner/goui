package main

import (
    "regexp" 
    //"text/template"
    "os/exec"
    "strings"
    "bytes"
    "log"
    "fmt"
)

type Command interface {
    GetName() string
    CanProcess(cmdText string) bool
    Execute(cmdText string, appCtx *AppContext) *CommandResult
}

type CommandResult struct {
    Args []string
}

type BaseCommand struct {
    Name string
    regex *regexp.Regexp 
}

func newBaseCommand(name, cmdExpression string) *BaseCommand {
     regex := regexp.MustCompile(cmdExpression)
     return &BaseCommand{Name: name, regex: regex}
}

func (r *BaseCommand) GetName() string {
    return r.Name
}

func (r *BaseCommand) CanProcess(cmdText string) bool {
    return r.regex.FindAllStringSubmatch(cmdText, -1) != nil
}

type ShellCommand struct {
    *BaseCommand
    template *TemplateEvaluator
}

func (r *ShellCommand) Execute(cmdText string, ctx *AppContext) *CommandResult {
    matches := r.regex.FindAllStringSubmatchIndex(cmdText, -1)
    cmdEndIndex := matches[0][1]
    log.Printf("%+v\n", cmdEndIndex)
    // split up the command and evaluate the expression
    cmdArray := strings.Split(cmdText[cmdEndIndex+1: len(cmdText)], " ")
    data := struct {
        Args []string
    }{
        cmdArray,
    }
    log.Printf("%+v\n", data)

    cmdString := r.template.ExecuteWithCtx(ctx, data)  
    log.Printf("%+v\n", cmdString)
    cmdArray = strings.Split(cmdString, " ")
    log.Printf("%+v\n", cmdArray)
    cmd := exec.Command(cmdArray[0], cmdArray[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
            log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
    return nil
}

func NewShellCommand(name, cmdExpression, shellExpression string) Command {
     shellExp := NewTemplateEvaluator(name, shellExpression)
     return &ShellCommand{newBaseCommand(name, cmdExpression), shellExp}
}

type CommandProcessor struct {
    Processors map[string]Command
}

func (c *CommandProcessor) Process(cmd string) Command {
    for _,v := range(c.Processors) {
        if v.CanProcess(cmd) {
            return v 
        }
    }
    return nil
}

func (c *CommandProcessor) Register(cmd Command) {
    c.Processors[cmd.GetName()] = cmd
}

func NewCommandProcessor() *CommandProcessor {
    return &CommandProcessor{Processors: make(map[string]Command)}
}
