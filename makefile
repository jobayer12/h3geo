.PHONY: build build-angular build-go

build: build-angular build-go

build-angular:
	cd geo-map-app && npm run build

build-go:
	cd geo-api && go build -o geo-api main.go