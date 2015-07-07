package writescript

import (
	"testing"
)

// initialize an instance we can test
var ws = WriteScript{}

//
// test empty plugin and data string
//
func TestGeneratorEmpty(t *testing.T) {
	err := ws.Process("", "", "", false)
	if err != nil {
		t.Error("failed", err)
	}
	println(ws.Content.AsString("\n", "\t"))
	if ws.Content.AsString("\n", "\t") != "" {
		t.Error("result not correct")
	}
}

//
// test simple plugin
//
func TestGeneratorPlugin(t *testing.T) {
	err := ws.Process("writeln('hello')", "", "", false)
	if err != nil {
		t.Error("failed", err)
	}
	if ws.Content.AsString("\n", "\t") != "hello\n" {
		t.Error("result not correct")
	}
}

//
// test simple plugin with empty data object
//
func TestGeneratorPluginData(t *testing.T) {
	err := ws.Process("writeln('hello')", "{}", "", false)
	if err != nil {
		t.Error("failed", err)
	}
	if ws.Content.AsString("\n", "\t") != "hello\n" {
		t.Error("result not correct")
	}
}

//
// test broken plugin
//
func TestGeneratorPluginBroken(t *testing.T) {
	err := ws.Process("writeln('hello'", "", "", false)
	if err == nil {
		t.Error("failed, no error was detected")
	}
}
