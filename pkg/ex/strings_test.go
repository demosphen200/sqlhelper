package ex

import (
	"testing"
)

func testSplitString(
	t *testing.T,
	str string,
	sep string,
	expectedFirst string,
	expectedSecond string,
	expectedFound bool,
) {
	first, second, found := SplitStringOnce(str, sep)

	if found != expectedFound {
		t.Errorf("expected found %v, actual %v", expectedFound, found)
	}
	if first != expectedFirst {
		t.Errorf("expected first %s, actual %s", expectedFirst, first)
	}
	if second != expectedSecond {
		t.Errorf("expected second %s, actual %s", expectedSecond, second)
	}

}

func TestSplitStringOnceFound(t *testing.T) {
	testSplitString(t, "first#second", "#", "first", "second", true)
}

func TestSplitStringOnceNotFound(t *testing.T) {
	testSplitString(t, "first", "#", "first", "", false)
}

func TestSplitStringOnceLongSep(t *testing.T) {
	testSplitString(t, "first###second", "###", "first", "second", true)
}

func TestSplitStringOnceSepAtEnd(t *testing.T) {
	testSplitString(t, "first#", "#", "first", "", true)
}

func TestSplitStringOnceSepAtStart(t *testing.T) {
	testSplitString(t, "#second", "#", "", "second", true)
}

func TestSplitStringOnceEmptyString(t *testing.T) {
	testSplitString(t, "", "#", "", "", false)
}

func TestSplitStringOnceOnlySep(t *testing.T) {
	testSplitString(t, "#", "#", "", "", true)
}
