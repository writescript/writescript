package writescript

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
	"os"
	"path/filepath"
	"github.com/paulvollmer/go-verbose"
	"strings"
)

const (
	// KeywordImport define the syntax for import
	KeywordImport = "#import "
)

// Plugin store url and code for the Plugin
type Plugin struct {
	ImportURLs      []string
	ImportCodeStack []string
	Js              []string // here we store the main plugin code
}

// ParseSource parse a source string
func (p *Plugin) ParseSource(src string) error {
	pluginLines := strings.Split(src, "\n")
	for _, v := range pluginLines {
		// fmt.Println("k", k, "v", v)

		if strings.Contains(v, KeywordImport) {
			// get the url...
			tmpURL := strings.Split(v, KeywordImport)

			// check if import already exists, or is not at the list of known urls
			if len(p.ImportURLs) == 0 || !IsValueInList(tmpURL[1], p.ImportURLs) {
				p.ImportURLs = append(p.ImportURLs, tmpURL[1])

				// file
				u, err := url.Parse(tmpURL[1])

				if err != nil {
					fmt.Println("Error", err)
					os.Exit(128)
				}
				if u.Scheme == "" {
					// read file
					dat, err := ioutil.ReadFile(tmpURL[1])
					if err != nil {
						fmt.Println(red("Error read file ", tmpURL[1]))
						fmt.Println(red(err))
						os.Exit(127)
					}
					// TODO: recursive
					p.ImportCodeStack = append(p.ImportCodeStack, string(dat))
					// fmt.Println("DATA:", string(dat))
				} else if u.Scheme == "http" || u.Scheme == "https" {
					// http request
					data, err := RequestPlugin(tmpURL[1])
					if err != nil {
						fmt.Println("Error", err)
						return err
					}
					p.ImportCodeStack = append(p.ImportCodeStack, string(data))
				}


			}
		} else {
			p.Js = append(p.Js, v)
		}
	}
	return nil
}

// RequestPlugin load a plugin source over http get
func RequestPlugin(url string) ([]byte, error) {
	resp, errReq := http.Get(url)
	if errReq != nil {
		return []byte{}, errReq
	}
	defer resp.Body.Close()
	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		return []byte{}, errBody
	}
	return body, nil
}

// LoadPlugin and return the source as a byte array
func LoadPlugin(src string, debug verbose.Verbose) ([]byte, error) {
	var err error
	var dataReturn []byte

	srcType := PluginIsType(src)
	debug.Println("--> type", srcType)
	switch srcType {

	case PluginTypeUnknown:
		dataReturn = []byte("")
		err = errors.New("No Plugin was set")
		fmt.Println(red("Error", err))
		break

	case PluginTypeFile:
		dataReturn, err = ioutil.ReadFile(src)
		if err != nil {
			fmt.Println(red("Error", err))
			os.Exit(127)
		}
		debug.Println("rile read", string(dataReturn))
		break

	case PluginTypeURL:
		dataReturn, err = RequestPlugin(src)
		if err != nil {
			fmt.Println(red("Error", err))
			os.Exit(127)
		}
		break

	case PluginTypeString:
		dataReturn = []byte(src)
		break
	}

	// fmt.Println("==> pluginBytes:", string(pluginBytes))
	return dataReturn, err
}

const (
	// PluginTypeUnknown enum
	PluginTypeUnknown = iota
	// PluginTypeFile enum
	PluginTypeFile
	// PluginTypeURL enum
	PluginTypeURL
	// PluginTypeString enum
	PluginTypeString
)

// PluginIsType return the PluginType
func PluginIsType(src string) int {
	theType := PluginTypeUnknown

	tmpExt := strings.ToLower(filepath.Ext(src))
	tmpURL, urlErr := url.Parse(src)

	if urlErr == nil && tmpURL.Scheme == "http" || tmpURL.Scheme == "https" {
		theType = PluginTypeURL
	} else if tmpExt == ".js" {
		theType = PluginTypeFile
	} else if src != "" {
		theType = PluginTypeString
	} else {
		fmt.Println(red("Error, unknown plugin type"))
	}

	return theType
}
