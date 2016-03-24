package writescript

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
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
				data, err := RequestPlugin(tmpURL[1])
				if err != nil {
					return err
				}
				p.ImportCodeStack = append(p.ImportCodeStack, string(data))
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
func LoadPlugin(src string) ([]byte, error) {
	var err error
	var dataReturn []byte

	switch PluginIsType(src) {

	case PluginTypeUnknown:
		dataReturn = []byte("")
		err = errors.New("No Plugin was set")
		break

	case PluginTypeFile:
		dataReturn, err = ioutil.ReadFile(src)
		break

	case PluginTypeURL:
		dataReturn, err = RequestPlugin(src)
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

	tmpExt := filepath.Ext(src)
	tmpURL, urlErr := url.Parse(src)

	if urlErr == nil && tmpURL.Scheme == "http" || tmpURL.Scheme == "https" {
		theType = PluginTypeURL
	} else if tmpExt == ".js" {
		theType = PluginTypeFile
	} else if src != "" {
		theType = PluginTypeString
	}

	return theType
}
