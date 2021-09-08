package main

import "example.com/xantios/myopia/router"

type ConfigStruct struct {
	Config GlobalConf  `yaml:"config"`
	Hosts  []Hosts `yaml:"hosts"`
	Assets []string `yaml:"assets"`
	Domains []string `yaml:"domains"`
}

type GlobalConf struct {
	Debug bool `yaml:"debug"`
	Docker bool `yaml:"docker"`
	Secure bool `yaml:"secure"`
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
}

type ExportConfig struct {
	debug bool
	docker bool
	secure bool
	hosts []router.Route
	assets []AssetMap
	domains []string
}

type AssetMap struct {
	Url string
	Path string
}
