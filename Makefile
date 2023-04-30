build-dev:
	docker build --no-cache -t registry.coruja.studio/go/ponzu-cms:$(tag) --build-arg UNAME=$(uname) --build-arg GID=$(IMAGE_GID) --build-arg UID=$(IMAGE_UID) --build-arg GO_VERSION=$(version) .

push-image:
	docker push registry.coruja.studio/go/ponzu-cms:$(tag)

run:
	docker run -it --rm -p 18080:8080 -p 10443:10443 -v $$(pwd)/test:/go/src/git.coruja.studio/go/application --entrypoint bash registry.coruja.studio/go/ponzu-cms:$(tag)

run-latest:
	make run tag=latest uname=$(uname)

build:
	go build github.com/monstrum/ponzu-cms/cmd/ponzu

ponzu-start:
	./ponzu build & ./ponzu run --port 8083
dev:
	make build
	make ponzu-start
