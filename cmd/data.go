package cmd

import "encoding/xml"

// Input data

type JsonHost struct {
	Port  int `json:"port"`
	Count int `json:"count"`
}

type JsonRegion map[string]JsonHost

type JsonRegions map[string]JsonRegion

type JsonHosts map[string]JsonRegions

type JsonVersions map[string]string

type JsonBrowser struct {
	DefaultVersion string       `json:"defaultVersion"`
	Versions       JsonVersions `json:"versions"`
}

type JsonQuota map[string]JsonBrowser

type JsonInput struct {
	Hosts   JsonHosts            `json:"hosts"`
	Quota   map[string]JsonQuota `json:"quota"`
	Aliases map[string][]string  `json:"aliases"`
}

// Output data

type XmlBrowsers struct {
	XMLName  xml.Name     `xml:"qa:browsers"`
	XmlNS    string       `xml:"xmlns:qa,attr"`
	Browsers []XmlBrowser `xml:"browser"`
}

type XmlBrowser struct {
	Name           string       `xml:"name,attr"`
	DefaultVersion string       `xml:"defaultVersion,attr"`
	Versions       []XmlVersion `xml:"version"`
}

type XmlVersion struct {
	Number  string      `xml:"number,attr"`
	Regions []XmlRegion `xml:"region"`
}

type XmlHosts []XmlHost

type XmlRegion struct {
	Name  string   `xml:"name,attr"`
	Hosts XmlHosts `xml:"host"`
}

type XmlHost struct {
	Name  string `xml:"name,attr"`
	Port  int    `xml:"port,attr"`
	Count int    `xml:"count,attr"`
}
