package main

import (
   "text/template"
    "bytes"
    "log"
)

type TemplateEvaluator struct {
    template *template.Template
}

func (r *TemplateEvaluator) ExecuteWithCtx(ctx *AppContext, data interface{}) string {
    argData := struct {
       Context interface{}
       Command interface{}
    } {
       ctx,
       data,
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

func NewTemplateEvaluator(name, expression string) *TemplateEvaluator {
    compiledTemplate := template.Must(template.New(name).Parse(expression))
    return &TemplateEvaluator{compiledTemplate}
}
