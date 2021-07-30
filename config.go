package main

import (
	//"bufio"

	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type CommandConfig struct {
	Name             string     `yaml:"name"`
	Regex            string     `yaml:"regex"`
	ResultType       string     `yaml:"resultType"`
	ResultKey        string     `yaml:"resultKey"`
	FilterExpression string     `yaml:"filterExpression"`
	ShellExpression  string     `yaml:"shellExpression"`
	PassthruSourceId string     `yaml:"sourceId"`
	RestSourceConfig string     `yaml:"restSource"`
	Static           string     `yaml:"static"`
	ViewId           string     `yaml:"viewId"`
	Help             HelpConfig `yaml:"help"`
}

type HelpConfig struct {
	Syntax      string `yaml:"syntax"`
	Description string `yaml:"description"`
}

type RestSourceConfig struct {
	GetUrlExpression   string `yaml:"getUrlExpression"`
	TokenExpression    string `yaml:"tokenExpression"`
	AuthScheme         string `yaml:"authScheme"`
	Username           string `yaml:"username"`
	PasswordExpression string `yaml:"passwordExpression"`
}

type ApplicationConfig struct {
	Name     string          `yaml:"name"`
	Pages    []PageConfig    `yaml:"pages"`
	Commands []CommandConfig `yaml:"commands"`
	Data     map[string]interface{}
}

type PageConfig struct {
	Name     string       `yaml:"name"`
	Id       string       `yaml:"id"`
	Shortcut string       `yaml:"shortcut"`
	Layout   LayoutConfig `yaml:"layout"`
	Views    []ViewConfig `yaml:"views"`
}

type LayoutConfig struct {
	Dir        string                  `yaml:"dir"`
	Containers []ContainerLayoutConfig `yaml:"containers"`
	Views      []ViewLayoutConfig      `yaml:"views"`
}

type ContainerLayoutConfig struct {
	FixedSize  int                     `yaml:"fixedSize"`
	Proportion int                     `yaml:"proportion"`
	Dir        string                  `yaml:"dir"`
	Containers []ContainerLayoutConfig `yaml:"containers"`
	Views      []ViewLayoutConfig      `yaml:"views"`
}

type ViewLayoutConfig struct {
	ViewId     string `yaml:"viewId"`
	FixedSize  int    `yaml:"fixedSize"`
	Proportion int    `yaml:"proportion"`
}

type TableItemConfig struct {
	HeaderExpression string `yaml:"headerExpression"`
	DataExpression   string `yaml:"dataExpression"`
}

type DataFormConfig struct {
	Fields []DataFormFieldConfig `yaml:"fields"`
}

type DataFormFieldConfig struct {
	Id              string `yaml:"id"`
	X               int    `yaml:"x"`
	Y               int    `yaml:"y"`
	Orientation     string `yaml:"orientation"` // v - vertical, h - horizontal
	LabelExpression string `yaml:"labelExpression"`
	LabelWidth      int    `yaml:"labelWidth"`
	ValueExpression string `yaml:"valueExpression"`
	ValueWidth      int    `yaml:"valueWidth"`
}

type ViewConfig struct {
	Id       string         `yaml:"id"`
	Name     string         `yaml:"name"`
	Shortcut string         `yaml:"shortcut"`
	DataPath string         `yaml:"dataPath"`
	Table    TableConfig    `yaml:"table"`
	Form     DataFormConfig `yaml:"form"`
	Static   string         `yaml:"static"`
}

type TableConfig struct {
	SelectExpression string            `yaml:"selectExpression"`
	Columns          []TableItemConfig `yaml:"columns"`
}

type DataConfig struct {
	Id   string      `yaml:"id"`
	Data interface{} `yaml:"data"`
}

func loadConfig(name string) (*ApplicationConfig, error) {
	baseDir := filepath.Dir(name)
	pagesDir := filepath.Join(baseDir, "pages")
	commandsDir := filepath.Join(baseDir, "commands")
	dataDir := filepath.Join(baseDir, "data")

	var app ApplicationConfig
	err := LoadYamlFromFile(name, &app)
	if err != nil {
		return nil, err
	}

	if FileExists(pagesDir) {
		err := loadFromDir(pagesDir, func(file string) error {
			var filePages []PageConfig
			loadErr := LoadYamlFromFile(file, &filePages)
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
			loadErr := LoadYamlFromFile(file, &fileCommands)
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
	if FileExists(dataDir) {
		err := loadFromDir(dataDir, func(file string) error {
			yamlData, loadErr := LoadGenericYamlFromFile(file)
			if loadErr != nil {
				return loadErr
			}
			app.Data = yamlData
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

func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}
