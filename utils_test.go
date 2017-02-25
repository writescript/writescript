package writescript

import (
	"testing"
)

// tests

func TestUtils_IsValueInList(t *testing.T) {
	var testData = []struct {
		val   string
		list  []string
		exist bool
	}{
		{"", []string{}, false},
		{"", []string{""}, true},
		{"", []string{"", "", ""}, true},
		{"1", []string{"123"}, false},
		{"1", []string{"1", "2", "3"}, true},
		{"1", []string{"data", "to"}, false},
		{"1", []string{"data", "to", "1"}, true},
		{"123", []string{"data", "to", "1"}, false},
		{"123", []string{"data", "to", "123"}, true},
		{"1 2 3", []string{"data", "to", "1"}, false},
		{"1 2 3", []string{"data", "to", "1 2 3"}, true},
		{"dolor", []string{"Lorem", "ipsum", "dolor", "sit", "amet,", "consectetur", "adipisicing", "elit,", "sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua.", "Ut", "enim", "ad", "minim", "veniam,", "quis", "nostrud", "exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea", "commodo", "consequat.", "Duis", "aute", "irure", "dolor", "in", "reprehenderit", "in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla", "pariatur.", "Excepteur", "sint", "occaecat", "cupidatat", "non", "proident,", "sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id", "est", "laborum."}, true},
	}

	for i, tt := range testData {
		exist := IsValueInList(tt.val, tt.list)
		if exist != tt.exist {
			t.Errorf("IsValueInList %v failed %q\n", i, tt.val)
		}
	}
}

// benchmarks
