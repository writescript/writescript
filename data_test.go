package writescript

import (
	"github.com/paulvollmer/go-verbose"
	"os"
	"testing"
)

var debug = verbose.New(os.Stdout, false)

// tests

func TestData_CheckSource(t *testing.T) {
	var testData = []struct {
		src string
		t   int
	}{
		{"", SourceDataJSON},
		{"test.json", SourceFileJSON},
		{"test.JSON", SourceFileJSON},
		{`{"some":"data"}`, SourceDataJSON},
		{`["some","data"]`, SourceDataJSON},
		{"test.yml", SourceFileYAML},
		{"test.YML", SourceFileYAML},
		{"test.yaml", SourceFileYAML},
		{"test.YAML", SourceFileYAML},
		{"some: data", SourceUnknown},
	}

	for i, tt := range testData {
		data := Data{}
		if data.CheckSource(tt.src) != tt.t {
			t.Errorf("CheckSource %q (%v) failed\n", tt.src, i)
		}
	}
}

func TestData_Init_Empty(t *testing.T) {
	data := Data{}
	data.Init("", *debug)
	if string(data.JSON) != "{}" {
		t.Error("Data Init Empty failed")
	}
}

func TestData_Init_JSON_File(t *testing.T) {
	data := Data{}
	data.Init("./fixture/testdata.json", *debug)
	if string(data.JSON) != `{"name":"testdata","description":"some data for testing"}
` {
		t.Error("Data Init JSON File failed")
	}
}

func TestData_Init_YAML_File(t *testing.T) {
	data := Data{}
	data.Init("./fixture/testdata.yml", *debug)
	if string(data.JSON) != `{"description":"some data for testing","name":"testdata"}` {
		t.Error("Data Init YAML File failed")
	}
}

// func TestData_Init_XML_File(t *testing.T) {
// 	data := Data{}
// 	data.Init("./fixture/testdata.xml")
// 	if string(data.JSON) != `{"description":"some data for testing","name":"testdata"}` {
// 		t.Error("Data Init XML File failed")
// 	}
// }

// func TestData_Init_CSV_File(t *testing.T) {
// 	data := Data{}
// 	data.Init("./fixture/testdata.csv")
// 	if string(data.JSON) != `[[1,2,3],[4,5,6]]` {
// 		t.Error("Data Init CSV File failed")
// 	}
// }

// benchmarks

var resultDataBench []byte

func benchmarkData(c string, b *testing.B) {
	var r []byte
	for n := 0; n < b.N; n++ {
		data := Data{}
		data.Init(c, *debug)
		r = data.JSON
	}
	resultDataBench = r
}

func Benchmark_Data_empty(b *testing.B) {
	benchmarkData("", b)
}
func Benchmark_Data_string(b *testing.B) {
	benchmarkData(`{"name":"testdata","description":"some data for testing"}`, b)
}
func Benchmark_Data_json(b *testing.B) {
	benchmarkData("./fixture/testdata.json", b)
}
func Benchmark_Data_yaml(b *testing.B) {
	benchmarkData("./fixture/testdata.yml", b)
}
