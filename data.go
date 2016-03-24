package writescript

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.pedge.io/pkg/yaml"
)

// Data initialize the data source for the writescript process.
type Data struct {
	// the writescript.Process need the data formatted as JSON.
	// this variable the process func can consume.
	JSON string
}

// Init initialize a new data source
func (d *Data) Init(src string) {
	switch d.CheckSource(src) {
	case SourceUnknown:
		d.JSON = src
		break
	case SourceFileJSON:
		d.ReadJSON(src)
		break
	case SourceDataJSON:
		// if source is empty, set the JSON to an empty object
		if src == "" {
			d.JSON = "{}"
		} else {
			d.JSON = src
		}
		break
	case SourceFileYAML:
		d.ReadYAML(src)
		break
	default:
		d.JSON = src
		break
	}
}

const (
	// SourceUnknown enum
	SourceUnknown = iota
	// SourceDataJSON enum
	SourceDataJSON
	// SourceFileJSON enum
	SourceFileJSON
	// SourceFileYAML enum
	SourceFileYAML
)

// CheckSource check if the src is a .json or .yml file or a string with data
func (d *Data) CheckSource(src string) int {
	// fmt.Println("CheckSource", src)
	tmpSourceType := SourceUnknown

	srcExt := filepath.Ext(src)
	switch srcExt {
	case ".json", ".JSON":
		tmpSourceType = SourceFileJSON
		break
	case ".yml", ".YML", ".yaml", ".YAML":
		tmpSourceType = SourceFileYAML
		break
	default:
		if src == "" {
			tmpSourceType = SourceDataJSON
		} else if string(src[0]) == "{" && string(src[len(src)-1]) == "}" {
			// check if the source looks like an json object / array
			tmpSourceType = SourceDataJSON
		} else if string(src[0]) == "[" && string(src[len(src)-1]) == "]" {
			tmpSourceType = SourceDataJSON
		}
	}
	return tmpSourceType
}

// ReadJSON read a json file
func (d *Data) ReadJSON(path string) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read json file", path)
		fmt.Println(err)
		os.Exit(1)
	}
	d.JSON = string(jsonBytes)
}

// ReadYAML read a yaml file
func (d *Data) ReadYAML(path string) {
	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read yaml file", path)
		fmt.Println(err)
		os.Exit(20)
	}

	// format yaml to json
	tmpJSON, errY2J := pkgyaml.ToJSON(yamlBytes, pkgyaml.ToJSONOptions{})
	if errY2J != nil {
		fmt.Println("decode yaml failed", errY2J)
		os.Exit(21)
	}
	d.JSON = string(tmpJSON)
}
