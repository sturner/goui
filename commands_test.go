package main

import (
	//"io/ioutil"
	"fmt"
	"testing"
)

func TestShellCommandRegex(t *testing.T) {

	ctx := &BaseAppContext{data: make(map[string]interface{}), pages: nil, app: nil, vPages: nil}
	args := []string{"shell_result.json"}
	ctx.RegisterArgs(args)
	cmdRegex := "show window \\w"
	cmdText := "show window shell_result.json"
	cmd := NewShellCommand("name", OutputJson, "result", cmdRegex, "", "view", "cat  ./test/{{index .Context.args 0}}")
	fmt.Printf("New shell command")
	canProcess := cmd.CanProcess(cmdText)
	if !canProcess {
		t.Errorf("Should have processed")
	}
	result, err := cmd.Execute(cmdText, ctx)
	if err != nil {
		t.Errorf("Error executing command [%+v]\n", err)
	}
	if result.Data == nil {
		t.Errorf("No data returned")
	}
	if result.Type != OutputJson {
		t.Errorf("Wrong result type")
	}
}
