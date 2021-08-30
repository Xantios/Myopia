default:
	@echo "Please provide a command"

love:
	@echo "not war(craft)"

clean:
	rm -rvf ./bin/ && mkdir bin

container:
	docker build -t myopia .

posix_x86:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=linux -e GOARCH=386 myopia:latest go build -o ./bin/myopia_posix_x86

posix_amd64:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=linux -e GOARCH=amd64 myopia:latest go build -o ./bin/myopia_posix_amd64

pi:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=linux -e GOARCH=arm -e GOARM=5 myopia:latest go build -o ./bin/myopia_posix_arm

mac_x86:
	echo "This option is deprecated!"
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=386 myopia:latest go build -o ./bin/myopia_mac_x86

mac_amd64:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=amd64 myopia:latest go build -o ./bin/myopia_mac_amd64

mac_silicon:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app -e GOOS=darwin -e GOARCH=arm64 myopia:latest go build -o ./bin/myopia_mac_m1

# For current Docker ENV (most likely Linux/AMD64)
build_dev:
	docker run --rm -v $(CURDIR):/usr/src/app -w /usr/src/app myopia:latest go build -o ./bin/myopia_dev

run_dev:
	@echo "\n==> Port gets mapped to 8100:42069\n"
	docker run --name myopia --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(CURDIR):/usr/src/app -w /usr/src/app -p 8100:42069 myopia:latest ./bin/myopia_dev

watch_dev:
	@echo "\n==> Port gets mapped to 8100:42069\n"
	docker run --name myopia --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(CURDIR):/usr/src/app -w /usr/src/app -p 8100:42069 myopia:latest watcher

dev: build_dev run_dev
watch: watch_dev
all: container posix_x86 posix_amd64 pi mac_amd64 mac_silicon
