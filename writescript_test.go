package writescript

import (
	"testing"
)

func TestWritescript_Empty(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "" {
		t.Error("result not correct", err)
	}
}

func TestWritescript_HeaderOn(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("", "", true)
	if err != nil || ws.Content.GetString("\n", "\t") != "// written by writescript v0.2.1\n// DO NOT EDIT!\n\n" {
		t.Error("header on failed", err)
	}
}

func TestWritescript_writeln_empty(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln()", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "\n" {
		t.Error("writeln empty failed", err)
	}
}

func TestWritescript_writeln(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello')", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "hello\n" {
		t.Error("writeln failed", err)
	}
}

func TestWritescript_write_empty(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("write()", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "" {
		t.Error("write failed", err)
	}
}

func TestWritescript_write(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("write('hello')", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "hello\n" {
		t.Error("write failed", err)
	}
}

func TestWritescript_pushLevel(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("pushLevel();write('hello')", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "\thello\n" {
		t.Error("pushLevel failed", err)
	}
}

func TestWritescript_popLevel(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("pushLevel();write('hello');popLevel();writeln('world')", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "\thello\nworld\n" {
		t.Error("pushLevel failed", err)
	}
}

func TestWritescript_getLevel(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln(getLevel());", "", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "0\n" {
		t.Error("getLevel failed", err)
	}
}

// func TestWritescript_setLevel(t *testing.T) {
// 	var ws = WriteScript{}
// 	err := ws.Process("setLevel(3);writeln('hello');", "", false)
// 	fmt.Println(ws.Content.GetString("\n", "-"))
// 	if err != nil || ws.Content.GetString("\n", "-") != "---hello\n" {
// 		t.Error("setLevel failed", err)
// 	}
// }

func TestWritescript_PluginAndEmptyDataObject(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello')", "{}", false)
	if err != nil || ws.Content.GetString("\n", "\t") != "hello\n" {
		t.Error("result not correct")
	}
}

func TestWritescript_PluginBroken(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello'", "", false)
	if err == nil {
		t.Error("failed, no error was detected")
	}
}
