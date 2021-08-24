package main

import "example.com/xantios/tinyproxy/router"

type ConfigStruct struct {
	Config GlobalConf  `yaml:"config"`
	Hosts  []Hosts `yaml:"hosts"`
	Assets []string `yaml:"assets"`
	Domains []string `yaml:"domains"`
}

type GlobalConf struct {
	Debug bool `yaml:"debug"`
}

type Item struct {
	Name        string `yaml:"name"`
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

type Hosts struct {
	Item Item `yaml:"item"`
}

type ConfigItem struct {
	name string
	source string
	destination string
	mapType router.RouteType
	spliceCount int // So we can check partial URLs against each other
}

type ExportConfig struct {
	debug bool
	hosts []ConfigItem
	assets []AssetMap
	domains []string
}

type AssetMap struct {
	url string
	path string
}
