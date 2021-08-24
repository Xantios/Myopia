package router

type RouteType string

const (
	MapHost RouteType = "MAP_HOST_TO_HOST"
	MapPath RouteType = "MAP_PATH_TO_HOST"
	MapAsset RouteType = "MAP_PATH_TO_PATH"
)

type Route struct {
	Name string
	Source string
	Destination string
	MapType RouteType
}

// var defaultRoute = ConfigItem{
var defaultRoute = Route{
	Name: "DefaultRoute",
	Source: "",
	Destination: "404_page.html",
	MapType: MapAsset,
}
