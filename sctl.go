package main

import (
	"flag"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/xml"
	"regexp"
	"strconv"
	"path"
	"os"
)

const (
	fileMode = 0644
)

var (
	inputFilePath = flag.String("inputFile", "input.json", "path to input file")
	outputDirectory = flag.String("outputDirectory", ".", "output directory")
	dryRun = flag.Bool("dryRun", false, "whether to send output to stdout instead of writing files")
)

// Input data

type JsonHost struct {
	Port int `json:"port"`
	Count int `json:"count"`
}

type JsonRegion map[string] JsonHost

type JsonRegions map[string]JsonRegion

type JsonHosts map[string]JsonRegions

type JsonVersions map[string] string

type JsonBrowser struct {
	DefaultVersion string `json:"defaultVersion"`
	Versions       JsonVersions `json:"versions"`
}

type JsonQuota map[string]JsonBrowser

type JsonInput struct {
	Hosts JsonHosts `json:"hosts"`
	Quota map[string]JsonQuota `json:"quota"`
}

// Output data

type XmlBrowsers struct {
	XMLName  xml.Name  `xml:"urn:config.gridrouter.qatools.ru browsers"`
	Browsers []XmlBrowser `xml:"browser"`
}

type XmlBrowser struct {
	Name           string    `xml:"name,attr"`
	DefaultVersion string    `xml:"defaultVersion,attr"`
	Versions       []XmlVersion `xml:"version"`
}

type XmlVersion struct {
	Number  string   `xml:"number,attr"`
	Regions []XmlRegion `xml:"region"`
}

type XmlHosts []XmlHost

type XmlRegion struct {
	Name  string `xml:"name,attr"`
	Hosts XmlHosts  `xml:"host"`
}

type XmlHost struct {
	Name   string `xml:"name,attr"`
	Port   int    `xml:"port,attr"`
	Count  int    `xml:"count,attr"`
}

func init() {
	flag.Parse()
}

func main() {
	input, err := parseInputFile(*inputFilePath)
	if (err != nil) {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	files := convert(*input)
	for name, browsers := range files {
		if err := output(name, browsers, *outputDirectory); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func convert(input JsonInput) map[string] XmlBrowsers {
	ret := make(map[string] XmlBrowsers)
	hostsMap := input.Hosts
	quotaMap := input.Quota
	for quotaName, quota := range quotaMap {
		browsers := []XmlBrowser{}
		for browserName, browser := range quota {
			xmlVersions := []XmlVersion{}
			for versionName, hostsRef := range browser.Versions {
				regions := hostsMap[hostsRef]
				if (regions != nil) {
					xmlVersion := XmlVersion{
						Number: versionName,
						Regions: jsonRegionsToXmlRegions(regions),
					}
					xmlVersions = append(xmlVersions, xmlVersion)
				} else {
					fmt.Printf("Missing host reference %s for browser %s:%s:%s", hostsRef, quotaName, browserName, versionName)
					os.Exit(1)
				}
			}
			xmlBrowser := XmlBrowser{
				Name: browserName,
				DefaultVersion: browser.DefaultVersion,
				Versions: xmlVersions,
			}
			browsers = append(browsers, xmlBrowser)
		}
		ret[quotaName] = XmlBrowsers{
			Browsers: browsers,
		}
	}
	return ret
}

func jsonRegionsToXmlRegions(regions JsonRegions) []XmlRegion {
	xmlRegions := []XmlRegion{}
	for regionName, region := range regions {
		xmlHosts := XmlHosts{}
		for hostPattern, host := range region {
			hostNames := parseHostPattern(hostPattern)
			for _, hostName := range hostNames {
				xmlHosts = append(xmlHosts, XmlHost{
					Name: hostName,
					Port: host.Port,
					Count: host.Count,
				})
			}
		}
		xmlRegions = append(xmlRegions, XmlRegion{
			Name: regionName,
			Hosts: xmlHosts,
		})
	}
	return xmlRegions
} 

func parseInputFile(filePath string) (*JsonInput, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading input file [%s]: %v", filePath, err))
	}
	input := new(JsonInput)
	if err := json.Unmarshal(bytes, input); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing input file [%s]: %v", filePath, err))
	}
	return input, nil
}

func marshalBrowsers(browsers XmlBrowsers) ([]byte, error) {
	return xml.Marshal(browsers)
}

func output(quotaName string, browsers XmlBrowsers, outputDirectory string) error {
	filePath := path.Join(outputDirectory, quotaName + ".xml")
	if (*dryRun) {
		return printOutputFile(filePath, browsers)
	} else {
		return saveOutputFile(filePath, browsers)
	}
}

func printOutputFile(filePath string, browsers XmlBrowsers) error {
	bytes, err := xml.Marshal(browsers)
	if (err != nil) {
		return err
	}
	fmt.Println(filePath)
	fmt.Println("---")
	fmt.Println(string(bytes))
	fmt.Println("---")
	return nil
}

func saveOutputFile(filePath string, browsers XmlBrowsers) error {
	bytes, err := marshalBrowsers(browsers)
	if (err != nil) {
		return err
	}
	if err := ioutil.WriteFile(filePath, bytes, fileMode); err != nil {
		return errors.New(fmt.Sprintf("error saving to output file [%s]: %v", filePath, err))
	}
	return nil
}

//Only one [1:10] pattern can be included in host pattern
func parseHostPattern(pattern string) []string {
	re := regexp.MustCompile("(.*)\\[(\\d+):(\\d+)\\](.*)")
	pieces := re.FindStringSubmatch(pattern)
	if len(pieces) == 5 {
		head := pieces[1]
		from, _ := strconv.Atoi(pieces[2])
		to, _ := strconv.Atoi(pieces[3])
		tail := pieces[4]
		if (from <= to) {
			ret := []string{}
			for i := from; i <= to; i++ {
				ret = append(ret, fmt.Sprintf("%s%d%s", head, i, tail))
			}
			return ret
		}
	} 
	return []string{pattern}
}