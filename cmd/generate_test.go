package cmd

import (
	. "github.com/aandryashin/matchers"
	"sort"
	"testing"
)

func TestParseInputFile(t *testing.T) {
	input, err := parseInputFile("../test-data/input.json")
	AssertThat(t, err, Is{nil})
	AssertThat(t, len(input.Hosts), EqualTo{3})
	AssertThat(t, len(input.Quota), EqualTo{1})
}

func TestParseHostPattern(t *testing.T) {
	AssertThat(t, parseHostPattern("test-host"), EqualTo{[]string{"test-host"}})
	AssertThat(t, parseHostPattern("test[1:3]-host"), EqualTo{[]string{"test1-host", "test2-host", "test3-host"}})
	AssertThat(t, parseHostPattern("[1:3]-host"), EqualTo{[]string{"1-host", "2-host", "3-host"}})
	AssertThat(t, parseHostPattern("host-[1:3]"), EqualTo{[]string{"host-1", "host-2", "host-3"}})
	AssertThat(t, parseHostPattern("host-[01:03]"), EqualTo{[]string{"host-01", "host-02", "host-03"}})
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
	AssertThat(t, browser.DefaultPlatform, EqualTo{"LINUX"})

	versions := browser.Versions
	AssertThat(t, len(versions), EqualTo{4})

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Number < versions[j].Number
	})

	firstVersion := versions[0]
	AssertThat(t, firstVersion.Number, EqualTo{"33.0"})
	secondVersion := versions[1]
	AssertThat(t, secondVersion.Number, EqualTo{"42.0"})
	thirdVersion := versions[2]
	AssertThat(t, thirdVersion.Number, EqualTo{"43.0"})
	AssertThat(t, thirdVersion.Platform, EqualTo{"WINDOWS"})
	fourthVersion := versions[3]
	AssertThat(t, fourthVersion.Number, EqualTo{"45.0"})
	AssertThat(t, fourthVersion.Platform, EqualTo{"LINUX"})

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

	thirdRegions := thirdVersion.Regions
	AssertThat(t, len(thirdRegions), EqualTo{1})
	region := thirdRegions[0]
	AssertThat(t, region.Name == "provider-1", Is{true})

	AssertThat(t, len(region.Hosts), EqualTo{5})

	fourthRegions := fourthVersion.Regions
	AssertThat(t, len(fourthRegions), EqualTo{1})
	fourthRegion := fourthRegions[0]
	AssertThat(t, fourthRegion.Name == "some-dc", Is{true})
	AssertThat(t, len(fourthRegion.Hosts), EqualTo{1})
	vncHost := fourthRegion.Hosts[0]
	AssertThat(t, vncHost.VNC, EqualTo{"ws://selenoid-host.example.com:4444/vnc"})

	for _, host := range region.Hosts {
		AssertThat(t, host.Username, EqualTo{"user1"})
		AssertThat(t, host.Password, EqualTo{"Password1"})
	}
}

func TestParseVersionPlatform(t *testing.T) {
	v, p := parseVersionPlatform("some-string")
	AssertThat(t, v, EqualTo{"some-string"})
	AssertThat(t, p, EqualTo{""})

	v, p = parseVersionPlatform("version@platform")
	AssertThat(t, v, EqualTo{"version"})
	AssertThat(t, p, EqualTo{"platform"})

	v, p = parseVersionPlatform("version@platform@platform")
	AssertThat(t, v, EqualTo{"version"})
	AssertThat(t, p, EqualTo{"platform@platform"})
}

func TestPreProcessVNC(t *testing.T) {
	AssertThat(t, preProcessVNC("selenoid-host.example.com", 4444, "selenoid"), EqualTo{"ws://selenoid-host.example.com:4444/vnc"})
	AssertThat(t, preProcessVNC("vnc-host.example.com", 5900, "vnc://$hostName:5900"), EqualTo{"vnc://vnc-host.example.com:5900"})
}

func TestGetPorts(t *testing.T) {
	AssertThat(t, getPorts(4444, ""), EqualTo{[]int{4444}})
	AssertThat(t, getPorts(4444, "4445"), EqualTo{[]int{4445}})
	AssertThat(t, getPorts(4444, "444[5:8]"), EqualTo{[]int{4445, 4446, 4447, 4448}})
	AssertThat(t, getPorts(4444, "44[5:8]4"), EqualTo{[]int{4454, 4464, 4474, 4484}})
	AssertThat(t, len(getPorts(4444, "NaN")), EqualTo{0})
}
