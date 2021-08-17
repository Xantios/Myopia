## Myopia
Short-sightedness where items in near eyesight are sharp but things in the background are not

## Configuration
It's a yaml file. Here is a example for reference:

```yaml
config:
  debug: true
hosts:
  - item:
      name: "SomeRemoteHost"
      source: "http://vhost.example.cop"
      destination: "http://server.behind.a.firewall.lan:5000"
  - item:
      name: "RemoteHostInAFolder"
      source: "/vhost"
      destination: "https://server.behind.a.firewall.lan:5000"
assets:
  - "/app:/var/www/"
```