package writescript

import (
	"testing"
)

// tests

func TestData_CheckSource(t *testing.T) {
	data := Data{}

	if data.CheckSource("") != SourceDataJSON {
		t.Error("CheckSource empty failed")
	}

	if data.CheckSource("test.json") != SourceFileJSON {
		t.Error("CheckSource json failed")
	}

	if data.CheckSource("test.JSON") != SourceFileJSON {
		t.Error("CheckSource JSON failed")
	}

	if data.CheckSource("test.yml") != SourceFileYAML {
		t.Error("CheckSource yml failed")
	}

	if data.CheckSource("test.yaml") != SourceFileYAML {
		t.Error("CheckSource yaml failed")
	}

	if data.CheckSource("test.YML") != SourceFileYAML {
		t.Error("CheckSource YML failed")
	}

	if data.CheckSource("test.YAML") != SourceFileYAML {
		t.Error("CheckSource YAML failed")
	}

	if data.CheckSource(`{"some":"data"}`) != SourceDataJSON {
		t.Error("CheckSource data json object failed")
	}

	if data.CheckSource(`["some","data"]`) != SourceDataJSON {
		t.Error("CheckSource data json array failed")
	}

	if data.CheckSource(`some: data`) != SourceUnknown {
		t.Error("CheckSource data unknown failed")
	}
}

func TestData_Init_Empty(t *testing.T) {
	data := Data{}
	data.Init("")
	if data.JSON != "{}" {
		t.Error("Data Init Empty failed")
	}
}

func TestData_Init_JSON_File(t *testing.T) {
	data := Data{}
	data.Init("./fixture/testdata.json")
	if data.JSON != `{"name":"testdata","description":"some data for testing"}
` {
		t.Error("Data Init JSON File failed")
	}
}

func TestData_Init_YAML_File(t *testing.T) {
	data := Data{}
	data.Init("./fixture/testdata.yml")
	if data.JSON != `{"description":"some data for testing","name":"testdata"}` {
		t.Error("Data Init YAML File failed")
	}
}

// benchmarks
