
export

build_container:
	docker build -t tiger_server --build-arg GOPROXY --build-arg GO111MODULE -f Dockerfile .

run:
	docker run --rm --env GOPROXY --env GO111MODULE tiger_server

