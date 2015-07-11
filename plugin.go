package writescript

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func (p *Plugin) Init(src string) {
	p.ImportURLs, p.ImportCodeStack, p.Js = ParseSource(src)
	// fmt.Println("ImportURLs", p.ImportURLs)
	// fmt.Println("ImportCodeStack", p.ImportCodeStack)
	// fmt.Println("js", strings.Join(p.Js, "\n"))
}

func ParseSource(src string) ([]string, []string, []string) {
	tmpImportURLs := []string{}
	tmpImportCodeStack := []string{}
	tmpJavascript := []string{}

	pluginLines := strings.Split(src, "\n")
	for _, v := range pluginLines {
		// fmt.Println("k", k, "v", v)

		if strings.Contains(v, KEYWORD_IMPORT) {
			// get the url...
			tmpUrl := strings.Split(v, KEYWORD_IMPORT)

			// check if import already exists
			if len(tmpImportURLs) == 0 {
				// 	fmt.Println("FIRST PLUGIN IMPORT", importUrl)
				tmpImportURLs = append(tmpImportURLs, tmpUrl[1])
				data := Request(tmpUrl[1])
				// fmt.Println("DATA", string(data))
				tmpImportCodeStack = append(tmpImportCodeStack, string(data))
			} else {
				// fmt.Println("check if url exists at array... '" + tmpUrl[1] + "'")

				// Is url already Known?
				if !isValueInList(tmpUrl[1], tmpImportURLs) {
					// fmt.Println("ADD NEW URL...")
					tmpImportURLs = append(tmpImportURLs, tmpUrl[1])
					data := Request(tmpUrl[1])
					tmpImportCodeStack = append(tmpImportCodeStack, string(data))
				}
			}

		} else {
			tmpJavascript = append(tmpJavascript, v)
		}
	}

	return tmpImportURLs, tmpImportCodeStack, tmpJavascript
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func Request(url string) []byte {
	resp, errReq := http.Get(url)
	if errReq != nil {
		fmt.Println("cannot import source", url)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		fmt.Println("require read body failed", errBody)
		os.Exit(2)
	}
	// fmt.Println("file", string(body))
	return body
}
