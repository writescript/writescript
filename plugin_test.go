package writescript

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// tests

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

func Test_PluginIsType_Unknown(t *testing.T) {
	if PluginIsType("") != PluginTypeUnknown {
		t.Error("Source type is not unknown")
	}
}

func Test_PluginIsType_File(t *testing.T) {
	if PluginIsType("test.js") != PluginTypeFile {
		t.Error("Source type is not javascript")
	}
}

func Test_PluginIsType_Url(t *testing.T) {
	if PluginIsType("http://test.js") != PluginTypeURL {
		t.Error("Source type is not an url")
	}
}

func Test_PluginIsType_String(t *testing.T) {
	if PluginIsType("test string") != PluginTypeString {
		t.Error("Source type is not a string")
	}
}

//
// test LoadPlugin function
//
func Test_LoadPlugin_Empty(t *testing.T) {
	result1, err := LoadPlugin("")
	if string(result1) != "" && err == nil {
		t.Error("returned plugin is not empty message", err)
	}
}

func Test_LoadPlugin_File(t *testing.T) {
	result2, err := LoadPlugin("./fixture/testplugin.js")
	if string(result2) != "console.log('hello world')\n" {
		t.Error("returned plugin is incorrect", err)
	}
}

func Test_LoadPlugin_Url(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "ok...")
	}))
	defer server.Close()
	result3, err := LoadPlugin(server.URL)
	if string(result3) != "ok..." {
		t.Error("returned plugin is incorrect", err)
	}
}

func Test_LoadPlugin_StringSource(t *testing.T) {
	result, err := LoadPlugin("var foo = 1337;")
	if string(result) != "var foo = 1337;" {
		t.Error("returned plugin source incorrect", err)
	}
}

// benchmarks
