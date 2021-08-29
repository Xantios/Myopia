package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type DynamicHost struct {
	ContainerName string
	Url string
	Port int
	Ip string
}

// MapType can be HOST or SUBDIR
// Make this compatible with RouteType from router
var MapType = "HOST"
var MapDomain = ""

// cli Docker Client Interface
var cli *client.Client

// DockerLog logrus instance
var DockerLog *logrus.Entry

func ContainerMapType(maptype string,domain string) {

	if maptype == "HOST" || maptype == "SUBDIR" {
		MapType = maptype
		MapDomain = domain
		return
	}

	panic("Container map type can be HOST or SUBDIR, invalid value supplied ("+maptype+")")
}

func GetMapType() string {
	return MapType
}

type CreateCallbackFunction func(host DynamicHost)
type RemoveCallbackFunction func(name string)

func Subscribe(socketPath string,create CreateCallbackFunction,remove RemoveCallbackFunction) {

	if socketPath == "" {
		socketPath = "unix:///var/run/docker.sock"
	}

	// Setup logging
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	DockerLog = logrus.WithFields(logrus.Fields{
		"service": "Docker",
	})

	var err error
	cli,err = client.NewClientWithOpts(client.WithHost(socketPath),client.FromEnv)

	if err != nil {
		panic(err)
	}

	filter := filters.NewArgs()
	filter.Add("type","container")
	filter.Add("event","start")
	filter.Add("event","die")

	messageChannel,errorChannel := cli.Events(context.Background(),types.EventsOptions{Filters: filter})

	for {
		select {
			case err := <-errorChannel:
				println("Error: "+err.Error())
			case messageChannel := <-messageChannel:

				// println("Type beat:"+messageChannel.Action)

				if messageChannel.Action == "die" {
					removeContainerFromProxy(messageChannel,remove)
				}

				if messageChannel.Action == "start" {
					addContainerToProxy(messageChannel,create)
				}
		}
	}
}

// addContainerToProxy Check if a label myopia.host exists
/*

	check if value of label syntax matches:
	/somePath or vhost.tld

	First port is assumed if not defined in value
	set port by using:
	/somePath:1337 or vhost.tld:1337
*/
func addContainerToProxy(msg events.Message,callback CreateCallbackFunction) {

	DockerLog.Info("Starting inspect on #"+msg.Actor.ID)
	container,err := cli.ContainerInspect(context.Background(),msg.Actor.ID)

	if err != nil {
		println("Error while inspecting container: "+err.Error())
		return
	}

	DockerLog.Info("Inspecting container: "+strings.TrimPrefix(container.Name,"/"))

	for key,value := range container.Config.Labels {
		// DockerLog.Debug("Label "+key," => "+value)

		if key == "myopia.host" {
			dynamicHost := parseHostString(value,container)
			callback(dynamicHost)
		}
	}

	// No match
	return
}

func removeContainerFromProxy(msg events.Message,callback RemoveCallbackFunction) {

	DockerLog.Info("Container removed #"+msg.Actor.ID)
	name := ""

	for k,v := range msg.Actor.Attributes {
		if k == "name" {
			name = v
		}
	}

	callback(name)
}

func parseHostString(value string,container types.ContainerJSON) DynamicHost {

	var url string
	var port int = 0

	// Extract ports
	var ports []int
	for k,_ := range container.NetworkSettings.NetworkSettingsBase.Ports {
		prefix := strings.SplitN(string(k),"/",2)[0]
		port,_ = strconv.Atoi(prefix)
		ports = append(ports, port)
	}

	// Prefix protocol
	var prefix = ""
	if !strings.HasPrefix(value,"/") {
		prefix = "http://"
	}

	if strings.Contains(value,":") {
		split := strings.SplitN(value,":",2)
		url = split[1]
		tmp, _ := strconv.ParseInt(split[2],10,4)
		port = int(tmp)
	} else {
		url = value
		port = ports[0]
	}

	return DynamicHost{
		ContainerName: strings.TrimPrefix(container.Name,"/"),
		Url: prefix+url,
		Port: port,
		Ip: container.NetworkSettings.DefaultNetworkSettings.IPAddress,
	}
}
