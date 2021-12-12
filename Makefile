build:
	GOARCH=amd64 GOOS=linux go build -o handler/main handler/main.go
	zip -j ./handler/main.zip ./handler/main

deploy:build 
	cdktf deploy
