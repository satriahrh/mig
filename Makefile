dist:
	@rm ./mig-*
	GOOS=darwin GOARCH=amd64 go build -o ./mig-darwin64 ./cmd/mig
	GOOS=linux GOARCH=amd64 go build -o ./mig-linux64 ./cmd/mig
	GOOS=linux GOARCH=386 go build -o ./mig-linux386 ./cmd/mig

