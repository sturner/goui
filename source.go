package main

import "github.com/go-resty/resty/v2"

var Client = resty.New()

type Source interface {
	Retrieve(ctx AppContext)
}

type RestSource struct {
	getTemplate *TemplateEvaluator
}

func NewRestSource(getExpression string) *RestSource {
	getExp := NewTemplateEvaluator(getExpression)
	src := RestSource{getTemplate: getExp}
	return &src
}

func (c *RestSource) Retrieve(ctx AppContext) *RestSource {
	return nil

}
