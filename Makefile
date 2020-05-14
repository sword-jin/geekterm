drawin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o geekterm-drawin ./cli/main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o geekterm-linux ./cli/main.go

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o geekterm.exe ./cli/main.go

all: drawin linux windows