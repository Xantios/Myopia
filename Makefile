default:
	@echo "Please provide a command"

love:
	@echo "not war(craft)"

clean:
	rm -rvf ./bin/ && mkdir bin

container:
	docker build -t tiny_proxy .

posix_x86:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=linux -e GOARCH=386 tiny_proxy:latest go build -o ./bin/tinyproxy_posix_x86

posix_amd64:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=linux -e GOARCH=amd64 tiny_proxy:latest go build -o ./bin/tinyproxy_posix_amd64

pi:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=linux -e GOARCH=arm -e GOARM=5 tiny_proxy:latest go build -o ./bin/tinyproxy_posix_arm

mac_x86:
	echo "This option is deprecated!"
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=386 tiny_proxy:latest go build -o ./bin/tinyproxy_mac_x86

mac_amd64:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=amd64 tiny_proxy:latest go build -o ./bin/tinyproxy_mac_amd64

mac_silicon:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=arm64 tiny_proxy:latest go build -o ./bin/tinyproxy_mac_m1

build_dev:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app tiny_proxy:latest go build -o ./bin/tinyproxy_dev

run_dev:
	@echo "\n==> Port gets mapped to 8100:42069\n"
	docker run --name tinyproxy --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -p 8100:42069 tiny_proxy:latest ./bin/tinyproxy_dev

watch_dev:
	@echo "\n==> Port gets mapped to 8100:42069\n"
	docker run --name tinyproxy --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -p 8100:42069 tiny_proxy:latest watcher

dev: build_dev run_dev
watch: watch_dev
all: container posix_x86 posix_amd64 pi mac_amd64 mac_silicon
