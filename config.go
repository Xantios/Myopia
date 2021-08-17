package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func env(key string,fallback string) string {
	value, ok := os.LookupEnv(key)

	if ok {
		return value
	}

	return fallback
}

func pullConfigFile(configPath string) []byte {
	data,err := ioutil.ReadFile(configPath)

	if err != nil {
		panic("Cant read config file "+configPath)
	}

	return data
}

func PrintConf(configStruct ConfigStruct) {
	s, _ := json.MarshalIndent(configStruct, "", "\t");
	fmt.Print(string(s))
}

func getHosts(hostsConf []Hosts) []ConfigItem {
	var hosts []ConfigItem
	for _,item := range hostsConf {

		// Set MatchType
		var matchType RouteType

		// Check for :// in URL
		matched,_ := regexp.Match("\\S{3,}://",[]byte(item.Item.Source))

		// Check amount of slashes in case paths overlap (eg: /api/myRequest/ should match /api/)
		// spliceCount = (matched) : strings.Count(item.Item.Source,"/")-2 ? strings.Count(item.Item.Source,"/")
		var spliceCount int
		if matched {
			spliceCount = strings.Count(item.Item.Source,"/")-2
		} else {
			spliceCount = strings.Count(item.Item.Source,"/")
		}

		if matched {
			matchType = MapHost
		} else {
			matchType = MapPath
		}

		hosts = append(hosts, ConfigItem{
			name: 		 item.Item.Name,
			source:      item.Item.Source,
			destination: item.Item.Destination,
			mapType:     matchType,
			spliceCount: spliceCount,
		})
	}

	return hosts
}

func getAssets(assetsConf []string) []AssetMap {
	var assets []AssetMap
	for _,item := range assetsConf {
		seperatedItem := strings.SplitN(item,":",2)
		assets = append(assets, AssetMap{
			url: seperatedItem[0],
			path: seperatedItem[1],
		})
	}

	return assets
}

func GetConf(configPath string) ExportConfig {

	// Pull byte[] from disk
	configfile := pullConfigFile(configPath)

	// Marshal into structs so we can access it like native data
	configContent := ConfigStruct{}
	err := yaml.Unmarshal(configfile,&configContent)

	if err != nil {
		log.Fatalf("Yaml Error: %v", err)
	}

	// Convert to more easy to use map
	hosts := getHosts(configContent.Hosts)

	// Pre split the assets tags
	assets := getAssets(configContent.Assets)

	// PrintConf(configContent)
	return ExportConfig{
		debug: configContent.Config.Debug,
		hosts: hosts,
		assets: assets,
	}
}
