package main

type RouteType string

const (
	MapHost RouteType = "MAP_HOST_TO_HOST"
	MapPath RouteType = "MAP_PATH_TO_HOST"
	MapAsset RouteType = "MAP_PATH_TO_PATH"
)

var defaultRoute = ConfigItem{
	source: "",
	destination: "404_page.html",
	mapType: MapAsset,
	spliceCount: 0,
}

