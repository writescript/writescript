package writescript

import (
	"testing"
)

//
// test empty plugin and data string
//
func TestGeneratorEmpty(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("", "", "", false)
	if err != nil {
		t.Error("failed", err)
	}
	if ws.Content.GetString("\n", "\t") != "" {
		t.Error("result not correct")
	}
}

//
// test simple plugin
//
func TestGeneratorPlugin(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello')", "", "", false)
	if err != nil {
		t.Error("failed", err)
	}
	if ws.Content.GetString("\n", "\t") != "hello\n" {
		t.Error("result not correct")
	}
}

//
// test simple plugin with empty data object
//
func TestGeneratorPluginData(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello')", "{}", "", false)
	if err != nil {
		t.Error("failed", err)
	}
	if ws.Content.GetString("\n", "\t") != "hello\n" {
		t.Error("result not correct")
	}
}

//
// test broken plugin
//
func TestGeneratorPluginBroken(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello'", "", "", false)
	if err == nil {
		t.Error("failed, no error was detected")
	}
}

// func TestGeneratorWrite(t *testing.T) {
// 	_ = generator.Process(`write('foo')`, "")
// 	fmt.Println(generator.Content.GetString("\n", "\t"))
// }
