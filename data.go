package writescript

import (
	"io/ioutil"
	"path/filepath"

	"go.pedge.io/pkg/yaml"
	"github.com/paulvollmer/go-verbose"
)

// Data initialize the data source for the writescript process.
type Data struct {
	// the writescript.Process need the data formatted as JSON.
	// this variable the process func can consume.
	JSON []byte
}

// Init initialize a new data source
func (d *Data) Init(src string, debug verbose.Verbose) {
	debug.Println("==> check type")
	switch d.CheckSource(src) {
	case SourceUnknown:
		debug.Println("--> unknown type")
		d.JSON = []byte("{}")
		break

	case SourceFileJSON:
		debug.Println("--> type json")
		d.ReadJSON(src)
		break

	case SourceDataJSON:
		// if source is empty, set the JSON to an empty object
		if src == "" {
			debug.Println("--> empty data")
			d.JSON = []byte("{}")
		} else {
			debug.Println("--> set source")
			d.JSON = []byte(src)
		}
		break

	case SourceFileYAML:
		debug.Println("--> type yaml")
		d.ReadYAML(src)
		break

	default:
		d.JSON = []byte("{}")
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
func (d *Data) CheckSource(src string) (r int) {
	r = SourceUnknown

	srcExt := filepath.Ext(src)
	switch srcExt {
	case ".json", ".JSON":
		r = SourceFileJSON
		break
	case ".yml", ".YML", ".yaml", ".YAML":
		r = SourceFileYAML
		break
	default:
		if src == "" {
			r = SourceDataJSON
		} else if string(src[0]) == "{" && string(src[len(src)-1]) == "}" {
			// check if the source looks like an json object / array
			r = SourceDataJSON
		} else if string(src[0]) == "[" && string(src[len(src)-1]) == "]" {
			r = SourceDataJSON
		}
	}
	return r
}

// ReadJSON read a json file
func (d *Data) ReadJSON(path string) error {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		// fmt.Println("cannot read json file", path)
		// fmt.Println(err)
		// os.Exit(1)
		return err
	}
	d.JSON = jsonBytes
	return nil
}

// ReadYAML read a yaml file
func (d *Data) ReadYAML(path string) error {
	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		// fmt.Println("cannot read yaml file", path)
		// fmt.Println(err)
		// os.Exit(20)
		return err
	}

	// format yaml to json
	tmpJSON, errY2J := pkgyaml.ToJSON(yamlBytes, pkgyaml.ToJSONOptions{})
	if errY2J != nil {
		// fmt.Println("decode yaml failed", errY2J)
		// os.Exit(21)
		return errY2J
	}
	d.JSON = tmpJSON
	return nil
}
