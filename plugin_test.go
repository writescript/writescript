package writescript

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlugin_ParseSource_Simple(t *testing.T) {
	plugin := Plugin{}
	err := plugin.ParseSource(`console.log('hello test')`)
	if err != nil || len(plugin.ImportURLs) != 0 || len(plugin.ImportCodeStack) != 0 || len(plugin.Js) != 1 {
		t.Error("ParseSource failed")
	}
}

func TestPlugin_ParseSource_OneImport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "var foo = 'test'")
	}))
	defer server.Close()

	plugin := Plugin{}
	err := plugin.ParseSource(`#import ` + server.URL + `
console.log('hello test')`)
	if err != nil || len(plugin.ImportURLs) != 1 || len(plugin.ImportCodeStack) != 1 || len(plugin.Js) != 1 {
		t.Error("ParseSource failed")
	}
}

func TestPlugin_ParseSource_MultipleImports(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "var foo = 'test'")
	}))
	defer server.Close()

	plugin := Plugin{}
	err := plugin.ParseSource(`#import ` + server.URL + `/test1
#import ` + server.URL + `/test2
console.log('hello test')`)
	if err != nil || len(plugin.ImportURLs) != 2 || len(plugin.ImportCodeStack) != 2 || len(plugin.Js) != 1 {
		t.Error("ParseSource failed")
	}
}

func TestPlugin_ParseSource_ImportRequestFailed(t *testing.T) {
	plugin := Plugin{}
	err := plugin.ParseSource(`#import http://wrong.url
console.log('hello test')`)
	if err == nil {
		t.Error("ParseSource failed")
	}
}
