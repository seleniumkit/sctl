package cmd

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	ggr "github.com/aerokube/ggr/config"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	fileMode = 0644
)

var (
	outputDirectory string
	dryRun          bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate XML quota using JSON input file",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := parseInputFile(inputFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		files := convert(*input)
		names := []string{}
		for name := range files {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			if err := output(name, files[name], outputDirectory); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	},
}

func init() {
	initCommonFlags(generateCmd)
	generateCmd.PersistentFlags().StringVar(&outputDirectory, "outputDirectory", ".", "output directory")
	generateCmd.PersistentFlags().BoolVar(&dryRun, "dryRun", false, "whether to send output to stdout instead of writing files")
}

func convert(input Input) map[string]ggr.Browsers {
	ret := make(map[string]ggr.Browsers)
	hostsMap := input.Hosts
	quotaMap := input.Quota
	aliasesMap := input.Aliases
	for quotaName, quota := range quotaMap {
		ret[quotaName] = createQuota(quotaName, hostsMap, quota)
	}
	for quotaName, aliases := range aliasesMap {
		if _, ok := ret[quotaName]; ok {
			for _, alias := range aliases {
				ret[alias] = ret[quotaName]
			}
		} else {
			fmt.Printf("Missing reference quota %s\n", quotaName)
			os.Exit(1)
		}
	}
	return ret
}

func createQuota(quotaName string, hostsMap Hosts, quota Quota) ggr.Browsers {
	var browsers []ggr.Browser
	for browserName, browser := range quota {
		var xmlVersions []ggr.Version
		for versionName, hostsRef := range browser.Versions {
			regions := hostsMap[hostsRef]
			if regions != nil {
				version, platform := parseVersionPlatform(versionName)
				xmlVersion := ggr.Version{
					Number:  version,
					Regions: jsonRegionsToXmlRegions(regions),
				}
				if platform != "" {
					xmlVersion.Platform = platform
				}
				xmlVersions = append(xmlVersions, xmlVersion)
			} else {
				fmt.Printf("Missing host reference %s for browser %s:%s:%s\n", hostsRef, quotaName, browserName, versionName)
				os.Exit(1)
			}
		}
		xmlBrowser := ggr.Browser{
			Name:            browserName,
			DefaultVersion:  browser.DefaultVersion,
			DefaultPlatform: browser.DefaultPlatform,
			Versions:        xmlVersions,
		}
		browsers = append(browsers, xmlBrowser)
	}
	return ggr.Browsers{
		Browsers: browsers,
	}
}

func parseVersionPlatform(s string) (string, string) {
	const separator = "@"
	pieces := strings.Split(s, separator)
	if len(pieces) > 1 {
		return pieces[0], strings.Join(pieces[1:], separator)
	}
	return s, ""
}

func jsonRegionsToXmlRegions(regions Regions) []ggr.Region {
	var xmlRegions []ggr.Region
	for regionName, region := range regions {
		xmlHosts := ggr.Hosts{}
		for hostPattern, host := range region {
			hostNames := parseHostPattern(hostPattern)
			for _, hostName := range hostNames {
				h := ggr.Host{
					Name:     hostName,
					Port:     host.Port,
					Count:    host.Count,
					Username: host.Username,
					Password: host.Password,
				}
				if host.VNC != "" {
					h.VNC = preProcessVNC(hostName, host.Port, host.VNC)
				}
				xmlHosts = append(xmlHosts, h)
			}
		}
		xmlRegions = append(xmlRegions, ggr.Region{
			Name:  regionName,
			Hosts: xmlHosts,
		})
	}
	return xmlRegions
}

func preProcessVNC(hostName string, port int, vnc string) string {
	const selenoid = "selenoid"
	const hostPattern = "$hostName"
	if vnc == selenoid {
		return fmt.Sprintf("ws://%s:%d/vnc", hostName, port)
	}
	return strings.Replace(vnc, hostPattern, hostName, -1)
}

func parseInputFile(filePath string) (*Input, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading input file [%s]: %v", filePath, err))
	}
	input := new(Input)
	if err := json.Unmarshal(bytes, input); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing input file [%s]: %v", filePath, err))
	}
	return input, nil
}

func marshalBrowsers(browsers ggr.Browsers) ([]byte, error) {
	return xml.MarshalIndent(browsers, "", "    ")
}

func output(quotaName string, browsers ggr.Browsers, outputDirectory string) error {
	filePath := path.Join(outputDirectory, quotaName+".xml")
	if dryRun {
		return printOutputFile(filePath, browsers)
	} else {
		return saveOutputFile(filePath, browsers)
	}
}

func printOutputFile(filePath string, browsers ggr.Browsers) error {
	bytes, err := marshalBrowsers(browsers)
	if err != nil {
		return err
	}
	fmt.Println(filePath)
	fmt.Println("---")
	fmt.Println(string(bytes))
	fmt.Println("---")
	return nil
}

func saveOutputFile(filePath string, browsers ggr.Browsers) error {
	bytes, err := marshalBrowsers(browsers)
	if err != nil {
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
		hostnameFmt := fmt.Sprintf("%%s%%0%dd%%s", len(pieces[2]))
		if from <= to {
			ret := []string{}
			for i := from; i <= to; i++ {
				ret = append(ret, fmt.Sprintf(hostnameFmt, head, i, tail))
			}
			return ret
		}
	}
	return []string{pattern}
}
