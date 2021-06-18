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
	filterExpression := ".[] | select(.foo == $args[0])"
	t.Logf("Filter %+v\n", filterExpression)
	filter := NewJqFilter(filterExpression)
	t.Logf("Filter %+v\n", filter)

	input := []interface{}{
		map[string]interface{}{
			"foo": "hello",
		},
		map[string]interface{}{
			"foo": "there",
		},
	}

	ctx := []interface{}{"hello"}

	results := filter.Filter(input, ctx)

	t.Logf("Object Results %+v\n", results)

}
