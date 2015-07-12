package writescript

import (
	"io/ioutil"
	"net/http"
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

func (p *Plugin) Init(src string) (err error) {
	p.ImportURLs, p.ImportCodeStack, p.Js, err = PluginParseSource(src)
	// fmt.Println("ImportURLs", p.ImportURLs)
	// fmt.Println("ImportCodeStack", p.ImportCodeStack)
	// fmt.Println("js", strings.Join(p.Js, "\n"))
	return err
}

func PluginParseSource(src string) ([]string, []string, []string, error) {
	tmpImportURLs := []string{}
	tmpImportCodeStack := []string{}
	tmpJavascript := []string{}

	pluginLines := strings.Split(src, "\n")
	for _, v := range pluginLines {
		// fmt.Println("k", k, "v", v)

		if strings.Contains(v, KEYWORD_IMPORT) {
			// get the url...
			tmpUrl := strings.Split(v, KEYWORD_IMPORT)

			// check if import already exists, or is not at the list of known urls
			if len(tmpImportURLs) == 0 || !IsValueInList(tmpUrl[1], tmpImportURLs) {
				tmpImportURLs = append(tmpImportURLs, tmpUrl[1])
				data, err := PluginRequest(tmpUrl[1])
				if err != nil {
					return tmpImportURLs, tmpImportCodeStack, tmpJavascript, err
				}
				tmpImportCodeStack = append(tmpImportCodeStack, string(data))
			}
		} else {
			tmpJavascript = append(tmpJavascript, v)
		}
	}
	return tmpImportURLs, tmpImportCodeStack, tmpJavascript, nil
}

func PluginRequest(url string) ([]byte, error) {
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
