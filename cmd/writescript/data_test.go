package main

import (
	"testing"
)

func TestData_CheckSource(t *testing.T) {
	data := Data{}

	if data.CheckSource("") != SOURCE_JSON_DATA {
		t.Error("CheckSource empty failed")
	}

	if data.CheckSource("test.json") != SOURCE_JSON_FILE {
		t.Error("CheckSource json failed")
	}

	if data.CheckSource("test.JSON") != SOURCE_JSON_FILE {
		t.Error("CheckSource JSON failed")
	}

	if data.CheckSource("test.yml") != SOURCE_YAML_FILE {
		t.Error("CheckSource yml failed")
	}

	if data.CheckSource("test.yaml") != SOURCE_YAML_FILE {
		t.Error("CheckSource yaml failed")
	}

	if data.CheckSource("test.YML") != SOURCE_YAML_FILE {
		t.Error("CheckSource YML failed")
	}

	if data.CheckSource("test.YAML") != SOURCE_YAML_FILE {
		t.Error("CheckSource YAML failed")
	}

	if data.CheckSource(`{"some":"data"}`) != SOURCE_JSON_DATA {
		t.Error("CheckSource data json object failed")
	}

	if data.CheckSource(`["some","data"]`) != SOURCE_JSON_DATA {
		t.Error("CheckSource data json array failed")
	}

	if data.CheckSource(`some: data`) != SOURCE_UNKNOWN {
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
	data.Init("../../fixture/testdata.json")
	if data.JSON != `{"name":"testdata","description":"some data for testing"}
` {
		t.Error("Data Init JSON File failed")
	}
}

func TestData_Init_YAML_File(t *testing.T) {
	data := Data{}
	data.Init("../../fixture/testdata.yml")
	if data.JSON != `{"description":"some data for testing","name":"testdata"}` {
		t.Error("Data Init YAML File failed")
	}
}
