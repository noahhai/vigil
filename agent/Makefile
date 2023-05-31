.PHONY: build clean deploy

build:
	env CG_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/vigil-linux-x64 main.go
	env CG_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/vigil-windows-x64.exe main.go
	env CG_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/vigil-darwin-x64 main.go

clean:
	rm -rf bin && mkdir bin

deploy: build
	aws s3 sync ./bin s3://getvigil.io/agent/
	aws cloudfront create-invalidation --distribution-id E2VR2UI5YH62G7 --paths "/agent/*"
