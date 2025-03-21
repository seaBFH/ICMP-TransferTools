.PHONY: invoke invoke-win invoke-all image
invoke-all: invoke invoke-win
invoke:
	go build -o ./invoke/dist/invoke ./invoke
invoke-win:
	GOOS=windows GOARCH=amd64 go build -o ./invoke/dist/invoke.exe ./invoke

image: invoke-all
	docker build -t dev/it-tool:22.04 .
