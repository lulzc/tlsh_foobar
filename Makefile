SERVER_BINARY_NAME=tlshServer
CLIENT_BINARY_NAME=tlshClient

#GOARCH=amd64 GOOS=darwin go build -o hello-world-darwin main.go
#GOARCH=amd64 GOOS=linux go build -o hello-world-linux main.go
#GOARCH=amd64 GOOS=windows go build -o hello-world-windows main.go

build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -ldflags "-s -w" -o build/${SERVER_BINARY_NAME} server/main.go
	#CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -ldflags "-s -w" -o build/${CLIENT_BINARY_NAME} client/main.go
run:
	go run server/main.go

clean:
	go clean
	rm build/${SERVER_BINARY_NAME}
	#rm ${CLIENT_BINARY_NAME}
# security checks
# pre req:
# go install golang.org/x/vuln/cmd/govulncheck@latest
# go install honnef.co/go/tools/cmd/staticcheck@latest
# https://github.com/golangci/golangci-lint/releases


vulncheck:
	~/go/bin/govulncheck ./...

other:
	~/go/bin/golangci-lint run
