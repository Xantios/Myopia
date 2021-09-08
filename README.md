## Myopia
Short-sightedness where items in near eyesight are sharp but things in the background are not

## Configuration
It's a yaml file. Here is an example for reference:

```yaml
config:
  debug: true
hosts:
  - item:
      name: "SomeRemoteHost"
      source: "http://vhost.example.com"
      destination: "http://server.behind.a.firewall.lan:5000/api"
  - item:
      name: "RemoteHostInAFolder"
      source: "/vhost"
      destination: "https://server.behind.a.firewall.lan:5000"
assets:
  - "/app:/var/www/"
domains:
  - "example.com"
```


## Host to host
As you can see, you can add multiple types of mappings.
the first example given listens for a domain and maps it to another.

Example URL: `http://vhost.example.com/example`
Goes to: `http://server.beind.a.firewall.lan/api/example`

## Path to host
The second example given maps a host to a path. 

Example URL: `http://localhost/vhost/example`
Goest to: `http://server.behind.a.firewall.lan/example`

## Serving assets
Local files (also called assets or artifacts) can be server out. make sure paths to overlap hosts

## Domains
By design localhost is always allowed. all items found in the hosts array are added to this allow list to.
you can add more domains to this list by adding them in this list

## Running
make run
