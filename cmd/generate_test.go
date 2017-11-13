package cmd

import (
	. "github.com/aandryashin/matchers"
	"sort"
	"testing"
)

func TestParseInputFile(t *testing.T) {
	input, err := parseInputFile("../test-data/input.json")
	AssertThat(t, err, Is{nil})
	AssertThat(t, len(input.Hosts), EqualTo{2})
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
	AssertThat(t, len(versions), EqualTo{3})

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Number < versions[j].Number
	})

	firstVersion := versions[0]
	AssertThat(t, firstVersion.Number == "33.0", Is{true})
	secondVersion := versions[1]
	AssertThat(t, secondVersion.Number == "42.0", Is{true})
	thirdVersion := versions[2]
	AssertThat(t, thirdVersion.Number == "43.0", Is{true})

	firstRegions := firstVersion.Regions
	AssertThat(t, len(firstRegions), EqualTo{2})
	firstRegion := firstRegions[0]
	AssertThat(t, firstRegion.Name == "region-a" || firstRegion.Name == "region-b", Is{true})

	secondRegion := firstRegions[1]
	if firstRegion.Name == "region-a" {
		AssertThat(t, secondRegion.Name, EqualTo{"region-b"})
	} else {
		AssertThat(t, secondRegion.Name, EqualTo{"region-a"})
	}
	AssertThat(t, len(firstRegion.Hosts), EqualTo{20})
	AssertThat(t, len(secondRegion.Hosts), EqualTo{20})

	for i := 0; i < len(firstRegion.Hosts); i++ {
		firstHost := firstRegion.Hosts[i]
		secondHost := secondRegion.Hosts[i]

		AssertThat(t, firstHost.Username == "" && secondHost.Username == "", Is{true})
		AssertThat(t, firstHost.Password == "" && secondHost.Password == "", Is{true})
	}

	secondRegions := thirdVersion.Regions
	AssertThat(t, len(secondRegions), EqualTo{1})
	region := secondRegions[0]
	AssertThat(t, region.Name == "provider-1", Is{true})

	AssertThat(t, len(region.Hosts), EqualTo{5})

	for _, host := range region.Hosts {
		AssertThat(t, host.Username, EqualTo{"user1"})
		AssertThat(t, host.Password, EqualTo{"Password1"})
	}
}
