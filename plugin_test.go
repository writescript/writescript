package writescript

import (
	"github.com/paulvollmer/go-verbose"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var debugPlugin = verbose.New(os.Stdout, false)

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

func Test_PluginIsType(t *testing.T) {
	var testData = []struct {
		src string
		typ int
	}{
		{"", PluginTypeUnknown},
		{"test.js", PluginTypeFile},
		{"http://test.js", PluginTypeURL},
		{"test string", PluginTypeString},
	}

	for i, tt := range testData {
		if PluginIsType(tt.src) != tt.typ {
			t.Errorf("Source (%v) type is not correct type\n", i)
		}
	}
}

//
// test LoadPlugin function
//
func Test_LoadPlugin_Empty(t *testing.T) {
	result1, err := LoadPlugin("", *debugPlugin)
	if string(result1) != "" && err == nil {
		t.Error("returned plugin is not empty message", err)
	}
}

func Test_LoadPlugin_File(t *testing.T) {
	result2, err := LoadPlugin("./fixture/testplugin.js", *debugPlugin)
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
	result3, err := LoadPlugin(server.URL, *debugPlugin)
	if string(result3) != "ok..." {
		t.Error("returned plugin is incorrect", err)
	}
}

func Test_LoadPlugin_StringSource(t *testing.T) {
	result, err := LoadPlugin("var foo = 1337;", *debugPlugin)
	if string(result) != "var foo = 1337;" {
		t.Error("returned plugin source incorrect", err)
	}
}

// benchmarks
