package main

import (
	"strings"
	"testing"
)

func TestJsonPath(t *testing.T) {
	path := NewJsonPathEvaluator("$.hello")

	data := map[string]interface{}{
		"hello": "there",
		"data":  1,
	}
	result, err := path.Execute(data)
	if err != nil {
		t.Errorf("Error: %+v\n", err)
		t.Fail()
	}

	if strings.Compare(result.(string), "there") != 0 {
		t.Errorf("Value not correct")
	}
}

func TestJqFilter(t *testing.T) {
	data, err := LoadGenericYamlFromFile("test/configs/simple-parent-detail/data/courses.yaml")
	if err != nil {
		t.Errorf("Error loading YAML %+v\n", err)
	}
	ctx := &BaseAppContext{data: data, pages: nil, app: nil, vPages: nil}
	ctx.RegisterArgs([]string{"ENG-256"})
	t.Logf("Filter data %+v\n", data)
	filterExpression := ".[] | select(.location == $args[0])"
	t.Logf("Filter %+v\n", filterExpression)
	filter := NewJqFilter(filterExpression)
	t.Logf("Filter %+v\n", filter)

	classData := data["classes.json"]

	if _, ok := classData.([]map[string]interface{}); !ok {
		t.Logf("It's not a list of maps\n")
	}

	results := filter.Filter(data["classes.json"], ctx)

	t.Logf("Object Results %+v\n", results)

}
