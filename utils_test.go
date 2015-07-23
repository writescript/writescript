package writescript

import (
	"testing"
)

func TestUtils_IsValueInList(t *testing.T) {
	tmp1 := IsValueInList("test", []string{"data", "to"})
	if tmp1 != false {
		t.Error("IsValueInList failed")
	}

	tmp2 := IsValueInList("test", []string{"data", "to", "test"})
	if tmp2 != true {
		t.Error("IsValueInList failed")
	}
}
