package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var quotaName string

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show quota statistics using JSON input file",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := parseInputFile(inputFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		printStat(*input, quotaName)
	},
}

func init() {
	initCommonFlags(statCmd)
	statCmd.PersistentFlags().StringVar(&quotaName, "quotaName", "", "quota name")
}

func printStat(input JsonInput, quotaName string) {
	hostsMap := input.Hosts
	quotaMap := input.Quota
	aliasesMap := input.Aliases

	for qn, aliases := range aliasesMap {
		if _, ok := quotaMap[qn]; ok {
			for _, alias := range aliases {
				quotaMap[alias] = quotaMap[qn]
			}
		} else {
			fmt.Printf("Missing reference quota %s\n", qn)
			os.Exit(1)
		}
	}

	if len(quotaName) > 0 {
		if _, ok := quotaMap[quotaName]; ok {
			processJsonQuota(quotaName, quotaMap[quotaName], hostsMap)
		} else {
			fmt.Printf("Missing quota %s\n", quotaName)
			os.Exit(1)
		}
	} else {
		for qn, jsonQuota := range quotaMap {
			processJsonQuota(qn, jsonQuota, hostsMap)
		}
	}

}

func processJsonQuota(quotaName string, quota JsonQuota, hostsMap JsonHosts) {
	fmt.Println(quotaName)
	fmt.Println("---")
	for browserName, browser := range quota {
		for versionName, hostsRef := range browser.Versions {
			regions := hostsMap[hostsRef]
			if regions != nil {
				regionsStat := ""
				for regionName, region := range regions {
					regionTotal := 0
					for hostPattern, host := range region {
						hostNames := parseHostPattern(hostPattern)
						regionTotal += len(hostNames) * host.Count
					}
					regionsStat += fmt.Sprintf(" %s = %d", regionName, regionTotal)
				}
				if len(regionsStat) > 0 {
					fmt.Printf("%s %s%s\n", browserName, versionName, regionsStat)
				}
			}
		}
	}
	fmt.Println()

}
