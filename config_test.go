package main

import (
	"io/ioutil"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	_, err := loadConfig("test/configs/simple-parent-detail/master.yml")
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func loadFile(name string) string {
	dat, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(dat)
}
