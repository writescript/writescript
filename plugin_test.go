package writescript

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlugin_isValueInList(t *testing.T) {
	tmp1 := IsValueInList("test", []string{"data", "to"})
	if tmp1 != false {
		t.Error("isValueInList failed")
	}

	tmp2 := IsValueInList("test", []string{"data", "to", "test"})
	if tmp2 != true {
		t.Error("isValueInList failed")
	}
}

func TestPlugin_ParseSource_Simple(t *testing.T) {
	urls, pluginStack, js := ParseSource(`console.log('hello test')`)
	if len(urls) != 0 || len(pluginStack) != 0 || len(js) != 1 {
		t.Error("Error at ParseSource")
	}
}

func TestPlugin_ParseSource_OneImport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "var foo = 'test'")
	}))
	defer server.Close()

	urls, pluginStack, js := ParseSource(`#import ` + server.URL + `
console.log('hello test')`)
	if len(urls) != 1 || len(pluginStack) != 1 || len(js) != 1 {
		t.Error("Error at ParseSource")
	}
}

func TestPlugin_ParseSource_MultipleImports(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "var foo = 'test'")
	}))
	defer server.Close()

	urls, pluginStack, js := ParseSource(`#import ` + server.URL + `/test1
#import ` + server.URL + `/test2
console.log('hello test')`)
	if len(urls) != 2 || len(pluginStack) != 2 || len(js) != 1 {
		t.Error("Error at ParseSource")
	}
}

// func TestPlugin_ParseSource(t *testing.T) {
// 	plugin := Plugin{}
// 	plugin.Init("")
// 	if len(plugin.ImportURLs) != 0 || len(plugin.ImportCodeStack) != 0 || len(plugin.Js) != 0 {
// 		t.Error("Error at empty plugin initialisation")
// 	}
// }
