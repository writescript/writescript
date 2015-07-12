package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Data initialize the data source for the writescript process.
type Data struct {
	// the writescript.Process need the data formatted as JSON.
	// this variable the process func can consume.
	JSON string
}

func (d *Data) Init(src string) {
	switch d.CheckSource(src) {
	case SOURCE_UNKNOWN:
		d.JSON = src
		break
	case SOURCE_JSON_FILE:
		d.ReadJson(src)
		break
	case SOURCE_JSON_DATA:
		// if source is empty, set the JSON to an empty object
		if src == "" {
			d.JSON = "{}"
		} else {
			d.JSON = src
		}
		break
	case SOURCE_YAML_FILE:
		d.ReadYaml(src)
		break
	default:
		d.JSON = src
		break
	}
}

const (
	SOURCE_UNKNOWN = iota
	SOURCE_JSON_FILE
	SOURCE_JSON_DATA
	SOURCE_YAML_FILE
)

// check if the src is a .json or .yml file or a string with data
func (d *Data) CheckSource(src string) int {
	// fmt.Println("CheckSource", src)
	tmpSourceType := SOURCE_UNKNOWN

	srcExt := filepath.Ext(src)
	switch srcExt {
	case ".json", ".JSON":
		tmpSourceType = SOURCE_JSON_FILE
		break
	case ".yml", ".YML", ".yaml", ".YAML":
		tmpSourceType = SOURCE_YAML_FILE
		break
	default:
		if src == "" {
			tmpSourceType = SOURCE_JSON_DATA
		} else if string(src[0]) == "{" && string(src[len(src)-1]) == "}" {
			// check if the source looks like an json object / array
			tmpSourceType = SOURCE_JSON_DATA
		} else if string(src[0]) == "[" && string(src[len(src)-1]) == "]" {
			tmpSourceType = SOURCE_JSON_DATA
		}
	}
	return tmpSourceType
}

func (d *Data) ReadJson(path string) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read json file", path)
		fmt.Println(err)
		os.Exit(1)
	}
	d.JSON = string(jsonBytes)
}

func (d *Data) ReadYaml(path string) {
	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read yaml file", path)
		fmt.Println(err)
		os.Exit(1)
	}
	// format yaml to json
	var t map[string]interface{}
	err = yaml.Unmarshal(yamlBytes, &t)
	if err != nil {
		fmt.Println("parse yaml file failed", err)
		os.Exit(2)
	}

	// fmt.Printf("T %#v\n", t)
	j, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error", err)
	}

	d.JSON = string(j)
}