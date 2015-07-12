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
	KEYWORD_IMPORT = "#import "
)

type Plugin struct {
	ImportURLs      []string
	ImportCodeStack []string
	Js              []string // here we store the main plugin code
}

func (p *Plugin) ParseSource(src string) error {
	pluginLines := strings.Split(src, "\n")
	for _, v := range pluginLines {
		// fmt.Println("k", k, "v", v)

		if strings.Contains(v, KEYWORD_IMPORT) {
			// get the url...
			tmpUrl := strings.Split(v, KEYWORD_IMPORT)

			// check if import already exists, or is not at the list of known urls
			if len(p.ImportURLs) == 0 || !IsValueInList(tmpUrl[1], p.ImportURLs) {
				p.ImportURLs = append(p.ImportURLs, tmpUrl[1])
				data, err := p.request(tmpUrl[1])
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

func (p *Plugin) request(url string) ([]byte, error) {
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

	case PluginTypeString:
		dataReturn = []byte(src)
		break
	}

	// TODO: check if plugin keyword exists (in registry)
	// XXX: library of default plugins
	// fmt.Println("search if plugin is embedded...")
	// switch src {
	// case "golang-const":
	// 	// println("...golang-const")
	// 	dataReturn = PLUGIN_GOLANG_CONST
	// 	break
	// case "golang-cli":
	// 	// println("...golang-cli")
	// 	dataReturn = PLUGIN_GOLANG_CLI
	// 	break
	// default:
	// println("... default")
	// }

	// fmt.Println("==> pluginBytes:", string(pluginBytes))
	return dataReturn, err
}

const (
	PluginTypeUnknown = iota
	PluginTypeFile
	PluginTypeURL
	PluginTypeString
)

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
