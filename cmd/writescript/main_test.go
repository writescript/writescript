package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSourceIsTypeUnknown(t *testing.T) {
	if SourceIsType("") != SourceTypeUnknown {
		t.Error("Source type is not unknown")
	}
}

func TestSourceIsTypeJs(t *testing.T) {
	if SourceIsType("test.js") != SourceTypeJavascript {
		t.Error("Source type is not javascript")
	}
}

func TestSourceIsTypeJson(t *testing.T) {
	if SourceIsType("test.json") != SourceTypeJSON {
		t.Error("Source type is not json")
	}
}

func TestSourceIsTypeUrl(t *testing.T) {
	if SourceIsType("http://test.js") != SourceTypeURL {
		t.Error("Source type is not an url")
	}
}

func TestSourceIsTypeString(t *testing.T) {
	if SourceIsType("test string") != SourceTypeString {
		t.Error("Source type is not a string")
	}
}

//
// test ReadPlugin function
//
func TestReadPluginEmpty(t *testing.T) {
	result1, err := ReadPlugin("")
	if string(result1) != "" && err == nil {
		t.Error("returned plugin is not empty message", err)
	}
}

func TestReadPluginFileJs(t *testing.T) {
	result2, err := ReadPlugin("../../fixture/testplugin.js")
	if string(result2) != "console.log('hello world')\n" {
		t.Error("returned plugin is incorrect", err)
	}
}

func TestReadPluginFileJson(t *testing.T) {
	_, err := ReadPlugin("test.json")
	if err == nil {
		t.Error("error not correct", err)
	}
}

func TestReadPluginUrl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, "ok...")
	}))
	defer server.Close()
	result3, err := ReadPlugin(server.URL)
	if string(result3) != "ok..." {
		t.Error("returned plugin is incorrect", err)
	}
}

//
// test ReadData function
//
func TestReadData(t *testing.T) {
	result1, err := ReadData("")
	if result1 != "{}" {
		t.Error("returned data string is not empty object", err)
	}

	result2, err := ReadData("{}")
	if result2 != "{}" {
		t.Error("returned data string is not {}", err)
	}

	result3, err := ReadData("../../fixture/testdata.json")
	if result3 != `{  "name": "testdata",  "description": "some data for testing"}` {
		t.Error("returned data string incorrect", err)
	}

	_, err = ReadData("not-exist.json")
	if err == nil {
		t.Error("no error returned", err)
	}
}
