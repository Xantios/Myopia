package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)


// MapType can be HOST or SUBDIR
var MapType = "HOST"
var MapDomain = ""

// Client
var cli *client.Client

func ContainerMapType(maptype string,domain string) {

	if maptype == "HOST" || maptype == "SUBDIR" {
		MapType = maptype
		MapDomain = domain
		return
	}

	panic("Container map type can be HOST or SUBDIR, invalid value supplied ("+maptype+")")
}

func Subscribe(socketPath string) {

	if socketPath == "" {
		socketPath = "unix:///var/run/docker.sock"
	}

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

				println("Type beat:"+messageChannel.Action)

				if messageChannel.Action == "die" {
					removeContainerFromProxy(messageChannel)
				}

				if messageChannel.Action == "start" {
					addContainerToProxy(messageChannel)
				}
		}
	}
}

func addContainerToProxy(msg events.Message) {

	println("Starting inspect on #"+msg.Actor.ID)
	container,err := cli.ContainerInspect(context.Background(),msg.Actor.ID)

	fmt.Printf("%#v\n",container.Config.Env)

	if err != nil {
		println("Error while inspecting container: "+err.Error())
		return
	}



	/*for key,value := range msg.Actor.Attributes {
		if key == "name" {
			println("Tering zei Nijntje, een nieuwe container genaamd: "+value)
		}
	}*/

	fmt.Printf("%#v\n",msg)
}

func removeContainerFromProxy(msg events.Message) {
	fmt.Printf("%#v\n",msg)
}
