package writescript

import (
	"github.com/paulvollmer/go-verbose"
	"os"
	"testing"
)

var debug2 = verbose.New(os.Stdout, false)

// tests

func TestWritescript_Empty(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "" {
		t.Error("result not correct", err)
	}
}

func TestWritescript_HeaderOn(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("", "", true, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "// written by writescript v"+Version+"\n// DO NOT EDIT!\n\n" {
		t.Error("header on failed", err)
	}
}

func TestWritescript_writeln_empty(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln()", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "\n" {
		t.Error("writeln empty failed", err)
	}
}

func TestWritescript_writeln(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello')", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "hello\n" {
		t.Error("writeln failed", err)
	}
}

func TestWritescript_write_empty(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("write()", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "" {
		t.Error("write failed", err)
	}
}

func TestWritescript_write(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("write('hello')", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "hello\n" {
		t.Error("write failed", err)
	}
}

func TestWritescript_pushLevel(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("pushLevel();write('hello')", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "\thello\n" {
		t.Error("pushLevel failed", err)
	}
}

func TestWritescript_popLevel(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("pushLevel();write('hello');popLevel();writeln('world')", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "\thello\nworld\n" {
		t.Error("pushLevel failed", err)
	}
}

func TestWritescript_getLevel(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln(getLevel());", "", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "0\n" {
		t.Error("getLevel failed", err)
	}
}

func TestWritescript_setLevel(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello');setLevel(3);writeln('world');", "", false, *debug2)
	if string(ws.Content.Get("\n", "-")) != "hello\n---world\n" {
		t.Error("setLevel failed", err)
	}
}

func TestWritescript_PluginAndEmptyDataObject(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello')", "{}", false, *debug2)
	if err != nil || string(ws.Content.Get("\n", "\t")) != "hello\n" {
		t.Error("result not correct")
	}
}

func TestWritescript_PluginBroken(t *testing.T) {
	var ws = WriteScript{}
	err := ws.Process("writeln('hello'", "", false, *debug2)
	if err == nil {
		t.Error("failed, no error was detected")
	}
}

// benchmarks

var resultBench []byte

func benchmarkWritescript(c string, b *testing.B) {
	var r []byte
	for n := 0; n < b.N; n++ {
		var ws = WriteScript{}
		err := ws.Process(c, "", false, *debug2)
		if err != nil {
			panic(err)
		}
		r = ws.Content.Get("\n", "\t")
	}
	resultBench = r
}

func Benchmark_Writescript_small(b *testing.B) {
	benchmarkWritescript("write('hello')", b)
}

func Benchmark_Writescript_large(b *testing.B) {
	benchmarkWritescript(`write('hello')
	writeln('-world')
	writeln('Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.')
	writeln('end...')`, b)
}
