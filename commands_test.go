package main

import (
    //"io/ioutil"
    "testing"
)

func TestShellCommandRegex(t *testing.T) {
    var ctx AppContext
    cmdText := "show window"
    cmd := NewShellCommand("name", cmdText, "ls ./{{index .Command.Args 0}}")  
    canProcess := cmd.CanProcess(cmdText)
    if !canProcess {
        t.Errorf("Should have processed") 
    }
    cmd.Execute(cmdText + " test and this and that", &ctx)
}

