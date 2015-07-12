package writescript

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlugin_PluginIsValueInList(t *testing.T) {
	tmp1 := PluginIsValueInList("test", []string{"data", "to"})
	if tmp1 != false {
		t.Error("IsValueInList failed")
	}

	tmp2 := PluginIsValueInList("test", []string{"data", "to", "test"})
	if tmp2 != true {
		t.Error("IsValueInList failed")
	}
}

func TestPlugin_PluginParseSource_Simple(t *testing.T) {
	urls, pluginStack, js, err := PluginParseSource(`console.log('hello test')`)
	if err != nil || len(urls) != 0 || len(pluginStack) != 0 || len(js) != 1 {
		t.Error("ParseSource failed")
	}
}

func TestPlugin_PluginParseSource_OneImport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "var foo = 'test'")
	}))
	defer server.Close()

	urls, pluginStack, js, err := PluginParseSource(`#import ` + server.URL + `
console.log('hello test')`)
	if err != nil || len(urls) != 1 || len(pluginStack) != 1 || len(js) != 1 {
		t.Error("ParseSource failed")
	}
}

func TestPlugin_PluginParseSource_MultipleImports(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "var foo = 'test'")
	}))
	defer server.Close()

	urls, pluginStack, js, err := PluginParseSource(`#import ` + server.URL + `/test1
#import ` + server.URL + `/test2
console.log('hello test')`)
	if err != nil || len(urls) != 2 || len(pluginStack) != 2 || len(js) != 1 {
		t.Error("ParseSource failed")
	}
}

func TestPlugin_PluginParseSource_ImportRequestFailed(t *testing.T) {
	_, _, _, err := PluginParseSource(`#import http://wrong.url
console.log('hello test')`)
	if err == nil {
		t.Error("ParseSource failed")
	}
}
