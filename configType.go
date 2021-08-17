package main

type ConfigStruct struct {
	Config GlobalConf  `yaml:"config"`
	Hosts  []Hosts `yaml:"hosts"`
	Assets []string `yaml:"assets"`
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
	mapType RouteType
	spliceCount int // So we can check partial URLs against each other
}

type ExportConfig struct {
	debug bool
	hosts []ConfigItem
	assets []AssetMap
}

type AssetMap struct {
	url string
	path string
}
