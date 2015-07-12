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

// func (p *Plugin) Init(src string) error {
// 	err := p.ParseSource(src)
// 	// fmt.Println("ImportURLs", p.ImportURLs)
// 	// fmt.Println("ImportCodeStack", p.ImportCodeStack)
// 	// fmt.Println("js", strings.Join(p.Js, "\n"))
// 	return err
// }

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
