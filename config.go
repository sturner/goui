package main

import (
	//"bufio"
	"io/fs"
	"fmt"
	"strings"
	"io/ioutil"
	"path/filepath"
	"errors"
	"os"
	"gopkg.in/yaml.v2"
)

type CommandConfig struct {
	Name            string `yaml:"name"`
	Regex           string `yaml:"regex"`
	ShellExpression string `yaml:"shellExpression"`
	ResultType      string `yaml:"resultType"`
	ResultKey       string `yaml:"resultKey"`
}

type ApplicationConfig struct {
	Name  string       `yaml:"name"`
	Pages []PageConfig `yaml:"pages"`
	Commands []CommandConfig `yaml:"commands"`
}

type PageConfig struct {
	Name     string                  `yaml:"name"`
	Id       string                  `yaml:"id"`
	Shortcut string                  `yaml:"shortcut"`
	Layout   []ContainerLayoutConfig `yaml:"layout"`
}

type ContainerLayoutConfig struct {
	Dir        string                  `yaml:"dir"`
	Containers []ContainerLayoutConfig `yaml:"containers"`
	Views      []ViewLayoutConfig      `yaml:"views"`
}

type ViewLayoutConfig struct {
	View       string `yaml:"view"`
	FixedSize  int    `yaml:"fixedSize"`
	Proportion int    `yaml:"proportion"`
}

type PositionConfig struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type TableItemConfig struct {
	HeaderExpression string `yaml:"headerExpression"`
	DataExpression   string `yaml:"dataExpression"`
	Position         int    `yaml:"position"`
}

type FormItemConfig struct {
	X               int    `yaml:"x"`
	Y               int    `yaml:"y"`
	Orientation     string `yaml:"orientation"` // v - vertical, h - horizontal
	LabelExpression string `yaml:"labelExpression"`
	LabelWidth      int    `yaml:"labelWidth"`
	ValueExpression string `yaml:valueExpression"`
	ValueWidth      int    `yaml:"valueWidth"`
}

type ViewConfig struct {
	Name     string            `yaml:"name"`
	Shortcut string            `yaml:"shortcut"`
	Position PositionConfig    `yaml:"position"`
	Table    []TableItemConfig `yaml:"table"`
	Form     []FormItemConfig  `yaml:"form"`
}

func loadConfig(name string) (*ApplicationConfig, error) {
        baseDir := filepath.Dir(name)
        pagesDir := filepath.Join(baseDir, "pages")
        commandsDir := filepath.Join(baseDir, "commands")

	var app ApplicationConfig
        err := loadFromFile(name, &app)
	if err != nil {
		return nil, err
	}

        if FileExists(pagesDir) {
            err := loadFromDir(pagesDir, func(file string) error { 
                 var filePages []PageConfig 
                 loadErr := loadFromFile(file, &filePages)
                 if loadErr != nil {
                     return loadErr
                 } 
                 app.Pages = append(app.Pages, filePages...)
                 return nil
            }) 

            if err != nil {
                return nil, err
            }
        } 
        if FileExists(commandsDir) {
            err := loadFromDir(commandsDir, func(file string) error { 
                 var fileCommands []CommandConfig 
                 loadErr := loadFromFile(file, &fileCommands)
                 if loadErr != nil {
                     return loadErr
                 } 
                 app.Commands = append(app.Commands, fileCommands...)
                 return nil
            }) 

            if err != nil {
                return nil, err
            }
        }

	return &app, nil
}

func loadFromDir(dir string, action func(path string) error) error {

    return filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
                    return action(path)
		}
                return nil
	})
}

func loadFromFile(file string, holder interface{})  error {

        if !FileExists(file) {
           return errors.New("File does not exist or not accessible: [" + file + "]") 
        } 

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(data), holder)
	if err != nil {
		return  err
	}
        return nil
}

func FileExists(file string) bool {
    if _, err := os.Stat(file); err == nil {
         return true
    } 
    return false
}

/*
func loadCommandsFromDir(dir string) []CommandConfig, error {

        err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})

        
	var commands []ApplicationConfig
	err = yaml.Unmarshal([]byte(data), &app)
	if err != nil {
		return nil, err
	}
}
*/

