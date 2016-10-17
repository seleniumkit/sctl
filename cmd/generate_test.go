package cmd

import (
	"testing"
	. "github.com/aandryashin/matchers"
)

func TestParseInputFile(t *testing.T) {
	input, err := parseInputFile("../test-data/input.json")
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

func TestConvert(t *testing.T) {
	input, _ := parseInputFile("../test-data/input.json")
	output := convert(*input)
	_, containsKey := output["test-quota"]
	AssertThat(t, containsKey, Is{true})
	browsers := output["test-quota"].Browsers
	AssertThat(t, len(browsers), EqualTo{1})
	browser := browsers[0]
	AssertThat(t, browser.Name, EqualTo{"firefox"})
	AssertThat(t, browser.DefaultVersion, EqualTo{"33.0"})
	
	versions := browser.Versions
	AssertThat(t, len(versions), EqualTo{2})
	firstVersion := versions[0]
	AssertThat(t, firstVersion.Number == "33.0" || firstVersion.Number == "42.0", Is{true})
	secondVersion := versions[1]
	if (firstVersion.Number == "42.0") {
		AssertThat(t, secondVersion.Number, EqualTo{"33.0"})
	} else {
		AssertThat(t, secondVersion.Number, EqualTo{"42.0"})
	}
	
	regions := firstVersion.Regions
	AssertThat(t, len(regions), EqualTo{2})
	firstRegion := regions[0]
	AssertThat(t, firstRegion.Name == "region-a" || firstRegion.Name == "region-b", Is{true})
	
	secondRegion := regions[1]
	if (firstRegion.Name == "region-a") {
		AssertThat(t, secondRegion.Name, EqualTo{"region-b"})
	} else {
		AssertThat(t, secondRegion.Name, EqualTo{"region-a"})
	}
	AssertThat(t, len(firstRegion.Hosts), EqualTo{20})
	AssertThat(t, len(secondRegion.Hosts), EqualTo{20})
}