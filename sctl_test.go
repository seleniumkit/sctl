package main

import (
	"testing"
	. "github.com/aandryashin/matchers"
)

func TestParseInputFile(t *testing.T) {
	input, err := parseInputFile("example/input.json")
	AssertThat(t, err, Is{nil})
	AssertThat(t, len(input.Hosts), EqualTo{1})
	AssertThat(t, len(input.Quota), EqualTo{1})
}

func TestParseHostPattern(t *testing.T) {
	AssertThat(t, parseHostPattern("test-host"), EqualTo{[]string{"test-host"}})
	AssertThat(t, parseHostPattern("test[1:3]-host"), EqualTo{[]string{"test1-host", "test2-host", "test3-host"}})
	AssertThat(t, parseHostPattern("[1:3]-host"), EqualTo{[]string{"1-host", "2-host", "3-host"}})
	AssertThat(t, parseHostPattern("host-[1:3]"), EqualTo{[]string{"host-1", "host-2", "host-3"}})
}