package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	data, err := loadConfig("test/configs/simple-parent-detail/master.yml")
	if err != nil {
		t.Errorf("%+v", err)
	}

	fmt.Printf("%+v", data)

}

func loadFile(name string) string {
	dat, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(dat)
}
