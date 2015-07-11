package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/writescript/writescript"
)

// main cli tool
func main() {

	//
	// meta
	//
	app := cli.NewApp()
	app.Name = "writescript"
	app.Version = writescript.Version
	app.Usage = "writescript pulgin based source generator"
	app.Author = "Paul Vollmer"
	app.Email = "paul.vollmer@fh-potsdam.de"

	//
	// flags
	//
	app.Flags = []cli.Flag{
		// javascript code
		cli.StringFlag{
			Name:  "plugin, p",
			Value: "",
			Usage: "the generator plugin as file", // TODO: implement 	keywords
		},
		// json data
		cli.StringFlag{
			Name:  "data, d",
			Value: "",
			Usage: "the data as file or json string",
		},
		cli.StringFlag{
			Name:  "linebreak, l",
			Value: "\\n",
			Usage: "the linebreak for each row",
		},
		cli.StringFlag{
			Name:  "whitespace, w",
			Value: "\t",
			Usage: "the level whitespace",
		},
		cli.BoolFlag{
			Name:  "header-off, H",
			Usage: "disables header output",
		},
	}

	//
	// process
	//
	app.Action = func(c *cli.Context) {
		//
		// cli flags
		//
		flagPlugin := c.String("plugin")
		flagData := c.String("data")
		flagLinebreak := c.String("linebreak")
		flagWhitespace := c.String("whitespace")
		flagHeaderOff := c.Bool("header-off")

		// read plugin
		pluginBytes, err := ReadPlugin(flagPlugin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// read data
		dataBytes, err := ReadData(flagData)
		if err != nil {
			fmt.Println("//", err)
			fmt.Println("// try to run...")
		}

		// run the generator
		ws := writescript.WriteScript{}
		err = ws.Process(string(pluginBytes), dataBytes, !flagHeaderOff)
		if err != nil {
			fmt.Println("writescript plugin error!\n", err)
			os.Exit(1)
		}

		if flagLinebreak == "\\n" {
			flagLinebreak = "\n"
		}
		fmt.Println(ws.Content.GetString(flagLinebreak, flagWhitespace))
	}

	app.Run(os.Args)
}

const (
	SourceTypeUnknown    = 0
	SourceTypeString     = 1
	SourceTypeJavascript = 2
	SourceTypeJSON       = 3
	SourceTypeURL        = 4
	ExtensionJavascript  = ".js"
	ExtensionJSON        = ".json"
)

func SourceIsType(src string) int {
	theType := SourceTypeUnknown

	tmpExt := filepath.Ext(src)
	tmpURL, urlErr := url.Parse(src)

	if urlErr == nil && tmpURL.Scheme == "http" || tmpURL.Scheme == "https" {
		theType = SourceTypeURL
	} else if tmpExt == ExtensionJavascript {
		theType = SourceTypeJavascript
	} else if tmpExt == ExtensionJSON {
		theType = SourceTypeJSON
	} else if src != "" {
		theType = SourceTypeString
	}

	return theType
}

//
// ReadPlugin and return as byte array
//
func ReadPlugin(src string) ([]byte, error) {
	var err error
	var dataReturn []byte

	switch SourceIsType(src) {
	case SourceTypeUnknown:
		dataReturn = []byte("")
		err = errors.New("No Plugin was set")
		break
	case SourceTypeString:
		dataReturn = []byte(src)
		break
	case SourceTypeJavascript:
		dataReturn, err = ioutil.ReadFile(src)
		break
	case SourceTypeJSON:
		dataReturn = []byte("")
		err = errors.New("JSON as Plugin not supported")
		break
	case SourceTypeURL:
		resp, errReq := http.Get(src)
		if errReq != nil {
			err = errReq
		}
		defer resp.Body.Close()
		body, errBody := ioutil.ReadAll(resp.Body)
		if errBody != nil {
			err = errBody
		}
		dataReturn = body
		break
	}

	return dataReturn, err
}

//
// ReadData and return as string. (this must be formatted as json)
//
func ReadData(src string) (string, error) {
	if src == "" {
		return "{}", nil
	}
	dataTmp := ""
	// check if data-flag is a .json file or data as string
	if filepath.Ext(src) == ".json" {
		dataBytes, err := ioutil.ReadFile(src)
		if err != nil {
			return "", err
		}
		dataTmp = string(dataBytes)
	} else {
		dataTmp = src
	}
	// remove linebreaks (sonst error im js)
	return strings.Replace(dataTmp, "\n", "", -1), nil
}
