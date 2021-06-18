package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v2"
)

type TemplateEvaluator struct {
	template *template.Template
}

type JsonPathEvaluator struct {
	expression gval.Evaluable
}

func (r *TemplateEvaluator) ExecuteWithCtx(ctx AppContext) string {
	argData := struct {
		Context interface{}
	}{
		ctx.GetDataMap(),
	}
	return r.Execute(argData)
}

func (r *TemplateEvaluator) Execute(data interface{}) string {
	log.Printf("Execute: %+v\n", data)
	var buf bytes.Buffer
	err := r.template.Execute(&buf, data)
	if err != nil {
		log.Printf("Error: %+v", err)
	}
	return buf.String()
}

func NewTemplateEvaluator(expression string) *TemplateEvaluator {
	compiledTemplate := template.Must(template.New("template").Parse(expression))
	return &TemplateEvaluator{compiledTemplate}
}

func NewJsonPathEvaluator(expression string) *JsonPathEvaluator {
	path, err := jsonpath.New(expression)
	if err != nil {
		log.Fatal("Invalid json path:" + expression)
	}
	return &JsonPathEvaluator{expression: path}
}

func (r *JsonPathEvaluator) ExecuteWithCtx(ctx AppContext, data interface{}) (interface{}, error) {
	argData := map[string]interface{}{
		"Data": ctx.GetDataMap(),
	}
	return r.Execute(argData)
}

func (r *JsonPathEvaluator) Execute(data interface{}) (interface{}, error) {
	return r.expression(context.Background(), data)
}

type Filter interface {
	Filter(data, ctx interface{}) interface{}
}

type JqFilter struct {
	code *gojq.Code
}

type EmptyFilter struct {
}

func (j *EmptyFilter) Filter(data, ctx interface{}) interface{} {
	log.Printf("EmptyFilter\n")
	return data
}

func NewJqFilter(expression string) Filter {
	query, parseErr := gojq.Parse(expression)
	if parseErr != nil {
		log.Panicf("Invalid filter expression %+v\n", parseErr)
	}
	code, compileErr := gojq.Compile(
		query,
		gojq.WithVariables([]string{"$args"}),
	)
	if compileErr != nil {
		log.Panicf("Error compiling expression %+v\n", compileErr)
	}
	return &JqFilter{code: code}
}

func (j *JqFilter) Filter(data, ctx interface{}) interface{} {
	log.Printf("JqFilter [%+v]\n", j)
	results := make([]interface{}, 0)
	log.Printf("JqFilter data [%+v]\n", data)
	iter := j.code.Run(data, ctx)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			fmt.Println(err)
			continue
		}
		results = append(results, v)
	}
	log.Printf("JqFilter results [%+v]\n", results)

	return results
}

func LoadYamlFromFile(file string, holder interface{}) error {

	log.Printf("Loading file [%s]\n", file)
	if !FileExists(file) {
		log.Printf("File does not exist [%s]\n", file)
		return errors.New("File does not exist or not accessible: [" + file + "]")
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Error reading file [%s], [%+v]\n", file, err)
		return err
	}
	err = yaml.Unmarshal([]byte(data), holder)
	if err != nil {
		log.Printf("Error loading file [%s], [%+v]\n", file, err)
		return err
	}
	return nil
}
